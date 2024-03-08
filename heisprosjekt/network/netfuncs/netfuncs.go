package netfuncs

import (
	//"Heis/elevator"
	"Heis/network/localip"
	"Heis/network/peers"
	"flag"
	"fmt"
	"os"
	//"time"
)

func Network_FSM(peerChan chan peers.PeerUpdate) {
	for {
		select {
		case a := <-peerChan:
			//fmt.Printf("%+v\n", a)
			printPeerUpdate(a)
			// case a := <-ElevatorRx:
			// 	fmt.Printf("Received: %+v\n", a)
		}
	}
}

// type HelloMsg struct {
// 	Message string
// 	elev    elevator.Elevator
// 	Iter    int
// }

func InitNet() string {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	return id
}

func printPeerUpdate(p peers.PeerUpdate) {
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}

// func Bcast_message(ch chan HelloMsg, elev elevator.Elevator, id string) {
// 	message := HelloMsg{"Hello from " + id, elev, 0}
// 	for {
// 		message.Iter++
// 		ch <- message
// 		time.Sleep(1 * time.Second)
// 	}
// }
