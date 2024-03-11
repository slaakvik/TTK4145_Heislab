package peers

import (
	"Heis/network/conn"
	"fmt"
	"net"
	"sort"
	"time"
	"Heis/network/localip"
	"os"
	"strings"
)

type PeerUpdate struct {
	Peers    []string
	New      string
	Lost     []string
	Master   string
	IsMaster bool
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func Transmitter(port int, id string, isMaster *bool, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			if *isMaster && len(id) > 0 && id[0] != 'M' {
				id = "M" + id[1:]
			}
		}
		conn.WriteTo([]byte(id), addr)
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])
		fmt.Println("Jeg fikk noe")
		id := string(buf[:n])
		fmt.Printf("Her er IDen: %s\n", id)
		
		fmt.Println("Før lastseen")
		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()

		}
		fmt.Println("Etter lastseen")

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}
		// Check if the id is the master and set MasterId accordingly
		if id != "" {
			if (id[0] == 'M') && (p.Master == "") {
				p.Master = id
				updated = true //i am unsure if this is going to be problems if it is more than one master on the network
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}

// ---[what i have defined]---
func PrintUpdatedPeers(p PeerUpdate) {
	println()
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:      %q\n", p.Peers)
	fmt.Printf("  New:        %q\n", p.New)
	fmt.Printf("  Lost:       %q\n", p.Lost)
	fmt.Printf("  Master:     %q\n", p.Master)
	//fmt.Printf("  isMaster:   %v\n,", p.IsMaster)
	println()
}

func ExtractIpFromPeers(p PeerUpdate, peersIp []string) []string {
	//peersIp = ""
	for _, peer := range p.Peers {
		data := strings.Split(peer, "-")
		peersIp = append(peersIp, data[1])
	}
	return peersIp
}

func ExtractIpFromPeer(peer string) string {
	data := strings.Split(peer, "-")
	return data[1]
}

func ExtractProcessIdFromPeers(p PeerUpdate, peersProcessId []string) []string {
	//peersProcessId = nil // maybe not the best solution, but I need to reset the slice every time to not just add new peer. Alternetively, we could just remove the lost peers.
	for _, peer := range p.Peers {
		data := strings.Split(peer, "-")
		peersProcessId = append(peersProcessId, data[2])
	}
	return peersProcessId
}

func ExtractProcessIdFromPeer(peer string) string {
	data := strings.Split(peer, "-")
	return data[2]
}

func GetPeerUpdate(peerCh chan PeerUpdate, masterCh chan string) {
	for {
		p := <-peerCh
		fmt.Println("Hei her er jeg inni peerUpdate")
		PrintUpdatedPeers(p)
		//Send master to masterCh
		if p.Master != ""{
			masterId := p.Master
			masterCh <- masterId
		} // else : you are master

	}
}

func CheckId() string {
	id := ""
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("S-%s-%d", localIP, os.Getpid())
	return id
}




//---------------------------


/* func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {
    var buf [1024]byte
    var p PeerUpdate
    lastSeen := make(map[string]time.Time)
    conn := conn.DialBroadcastUDP(port)

    for {
        conn.SetReadDeadline(time.Now().Add(interval))
        n, _, _ := conn.ReadFrom(buf[0:])

        id := string(buf[:n])
        isMasterByte := buf[n]
        // Convert the byte to a boolean value
        if isMasterByte == 1{
			p.IsMaster=true
		} else {
			p.IsMaster=false
		}

        // Adding new connection
        p.New = ""
        if id != "" {
            if _, idExists := lastSeen[id]; !idExists {
                p.New = id
            }
            lastSeen[id] = time.Now()
        }

        // If isMaster is true, set the corresponding Master string
        if p.IsMaster && id != "" {
            p.Master = id
        }

        // Removing dead connection
        p.Lost = make([]string, 0)
        for k, v := range lastSeen {
            if time.Now().Sub(v) > timeout {
                p.Lost = append(p.Lost, k)
                delete(lastSeen, k)
            }
        }

        // Sending update if necessary
        if p.New != "" || len(p.Lost) > 0 || p.Master != "" {
            p.Peers = make([]string, 0, len(lastSeen))
            for k := range lastSeen {
                p.Peers = append(p.Peers, k)
            }
            sort.Strings(p.Peers)
            sort.Strings(p.Lost)
            peerUpdateCh <- p
        }
    }
}
*/
