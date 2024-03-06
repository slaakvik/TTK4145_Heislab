package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/master"
	"Heis/network/localip"
	"Heis/network/peers"
	"Heis/network/tcp"
	"Heis/slave"
	"flag"
	"fmt"
	"os"
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

	// Master and slaves
	var MasterId string
	//var MasterIp string ????? Why does it say not declared and used when I have used it??
	var MasterProcessId string

	var SlavesId []string
	var SlavesIp []string

	// Regarding the peers
	var PeersIp []string
	var PeersProcessId []string

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
			peers.PrintUpdatedPeers(data)

			PeersProcessId = peers.ExtractProcessIdFromPeers(data, PeersIp)
			masterIdIndex, err := master.ChooseMasterIndex(PeersProcessId)
			MasterId = master.ChooseMaster(masterIdIndex, data)
			if err != nil {
				fmt.Printf("[error] Error with choosing master %v", err)
				return
			}

			fmt.Printf("Noe galt 1\n")
			SlavesId = slave.ChooseSlaves(data, SlavesId, MasterId)
			fmt.Printf("Her er master: %v\n", MasterId)

			// Extracting the IP-adresses from the master and the slaves
			SlavesIp = slave.ExtractIpFromSlaves(SlavesId, SlavesIp)
			MasterIp := peers.ExtractIpFromPeer(MasterId)
			
			MasterProcessId = peers.ExtractProcessIdFromPeer(MasterId)
			//var katt = os.Getpid()
			if master.CheckIfYouAreMaster(MasterProcessId) {
				fmt.Printf("HÆLLÆ 1\n")
				for _, slaveIp := range SlavesIp {
					go tcp.Transmit(tcp.CONN_PORT, slaveIp, elev)
					go tcp.Receive(tcp.CONN_PORT, slaveIp, elevatorRx)
				}
			} else {
				fmt.Printf("HÆLLÆ 2\n")
				go tcp.Transmit(tcp.CONN_PORT, MasterIp, elev)
				go tcp.Receive(tcp.CONN_PORT, MasterIp, elevatorRx)
			} 
		case data := <-elevatorRx:
			fmt.Printf("Heisen: %v\n", data)
		}
	}
}




/** -----[Plan]-----
 * 
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
