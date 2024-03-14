package peers

import (
	"Heis/network/conn"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

type PeerUpdate struct {
	Peers  []string
	New    string
	Lost   []string
	Master string
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
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

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Adding new Master
		//p.Master = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.Master = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}
		//---------------------------

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			p.Master = determineMaster(p.Peers, p.Master)
			peerUpdateCh <- p
		}
	}
}

func PeerUpdates(peerCh chan PeerUpdate, SendMasterIdToReceive chan string, sendMasterIdToNotify chan string) {
	for {
		p := <-peerCh
		fmt.Println("Hei her er jeg inni peerUpdate")
		PrintUpdatedPeers(p)
		//Send master to masterCh
		if p.Master != "" {
			masterId := p.Master
			SendMasterIdToReceive <- masterId
			sendMasterIdToNotify <- masterId
		}
	}
}

func determineMaster(peers []string, masterId string) string {
	// Your logic to determine the master goes here
	// For simplicity, you can use the first peer as the master
	// if masterId != "" {
	// 	fmt.Println("Master is already set")
	// 	return masterId
	// }else if len(peers) > 0 {
	for len(peers) > 0 {
		return peers[0]
	}
	// }
	return "" // Return an empty string if there are no peers
}

// ---[what i have defined]---
func PrintUpdatedPeers(p PeerUpdate) {
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
	fmt.Printf("  Master:     %q\n", p.Master)
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

//---------------------------
