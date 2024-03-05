package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/network/localip"
	"Heis/network/peers"
	"Heis/network/tcp"
	"flag"
	"fmt"
	"os"
	"time"
)

// this was just to make main smaller. Just temporary
func CheckId(id string) string {
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		return id
	}
	return id
}

func main() {
	// Regarding the elevator itself
	_numFloors := elevio.NumFloors
	elevio.Init("localhost:15657", _numFloors)
	elev := elevator.InitElev()

	// Regarding the heartbeat
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	id = CheckId(id) //set the id to default if id is not provided at launch

	peerEnableTx := make(chan bool)
	peerUpdateRx := make(chan peers.PeerUpdate)

	go peers.Transmitter(15647, id, peerEnableTx)
	go peers.Receiver(15647, peerUpdateRx)

	// Regarding the receiver
	var elevatorRx = make(chan elevator.Elevator)

	// Regarding the transmitter
	//tcp.Transmit(tcp.CONN_PORT, tcp.CONN_HOST, elev)

	for {
		select {
		case data := <-peerUpdateRx:
			peers.PeersIp = peers.ExtractIpFromPeers(data)
			peers.PrintUpdatedPeers(data)
			for _, peerIp := range peers.PeersIp {
				go tcp.Receive(tcp.CONN_PORT, peerIp, elevatorRx)
				go tcp.Transmit(tcp.CONN_PORT, peerIp, elev)
			}
			time.Sleep(2 * time.Second)

		case data := <-elevatorRx:
			fmt.Printf("Heisen: %v\n", data)
		}
	}
}

/** -----[Plan]-----
 * I will just say that element one in Peers is the master. This is just to have a master.
 *
 *
 */

/**
 * [UDP to localize]

  + |||||| +

 * [Establish TCP connection]

  addr2, port2 {--------> } addr1, port1
			 +	 			+
  addr2, port2 {<--------} addr1, port1

// 1: the node to the left
// 2: the node to the right

* [Remove TCP connection if node goes offline]

*/
