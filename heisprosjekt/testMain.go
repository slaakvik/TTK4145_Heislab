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


/**
 * @var constant variables regarding the TCP 
 * these ports can be made dynamicly by iteration through the p.Peers that are active on the netework.
 * Then you know exactly how many ports to hold on to (for the master).
 * M is short for "master" and S is short for "slave".
 */
const (
	M_CONN_PORT_R1 = 50000
	M_CONN_PORT_R2 = 50001
	M_CONN_PORT_T1 = 50002 
	M_CONN_PORT_T2 = 50003 

	S_CONN_PORT_R1 = 55000
	S_CONN_PORT_T1 = 55001
	S_CONN_PORT_R2 = 55002
	S_CONN_PORT_T2 = 55003
)




func main() {

	//
	var ProcessId = os.Getpid()
	fmt.Printf("ProsessId er: %v\n", ProcessId)

	// Master and slaves
	var MasterId string
	var MasterIp string 
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
			SlavesId = slave.ChooseSlaves(data, SlavesId, MasterId)
			fmt.Printf("Master: %v\n", MasterId)

			// Extracting the IP-adresses from the master and the slaves
			SlavesIp = slave.ExtractIpFromSlaves(SlavesId, SlavesIp)
			MasterIp = peers.ExtractIpFromPeer(MasterId)
			
			MasterProcessId = peers.ExtractProcessIdFromPeer(MasterId)
			//var katt = os.Getpid()
			if master.CheckIfYouAreMaster(MasterProcessId, ProcessId) {
				for _, slaveIp := range SlavesIp {
					go tcp.Transmit(M_CONN_PORT_T1, slaveIp, elev)
					go tcp.Receive(M_CONN_PORT_R1, slaveIp, elevatorRx)
				}
			} else {
				fmt.Printf("Slave: %v\n", id)
				go tcp.Receive(S_CONN_PORT_R1, MasterIp, elevatorRx)
				go tcp.Transmit(S_CONN_PORT_T1, MasterIp, elev)
			} 
		case data := <-elevatorRx:
			fmt.Printf("Heisen: %v\n", data)
		}
	}
}




/** -----[Plan]-----
 * Må ha en variant som er for å kjøre TCP på en PC, og en annen for flere. Dette er bare for å gjøre testing lettere og mer fleksibelt. 
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
