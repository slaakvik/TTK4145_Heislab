package master

import (
	"Heis/elevator"
	"Heis/network/peers"
	"Heis/network/tcp"
	"fmt"
	"net"
	"strconv"
)

/*
-------------------------------------------------------
--! @file
--! @brief This file contain functions and variables regarding the master
-------------------------------------------------------
*/

/**
 * @func ChooseMasterIndex is not assigning or making a master, it only chooses who of the peers should be assigned, which is the peer with the lowest processID
 *
 */
func ChooseMasterIndex(peersProcessId []string) (int, error) {
	var SmallestProcessId, err = strconv.Atoi(peersProcessId[0])
	if err != nil {
		fmt.Printf("[error]: error converting string to int: %s\n", err)
		return 0, err // returning 0 is a problem here
	}
	var SmallestProcessIdIndex = 0
	for i, peer_str := range peersProcessId {
		peer_int, err := strconv.Atoi(peer_str)
		if err != nil {
			fmt.Printf("[error]: error converting string to int: %s\n", err)
			return 0, err // returning 0 is a problem here
		}
		if peer_int < SmallestProcessId {
			SmallestProcessId = peer_int
			SmallestProcessIdIndex = i
		}
	}
	return SmallestProcessIdIndex, err
}

func ChooseMaster(masterIdIndex int, p peers.PeerUpdate) string {
	masterId := p.Peers[masterIdIndex]
	return masterId
}

/**
 * @func ChechIfMasterIsLost checks if master is lost. Return true if it is lost.
 *
 */
/* func CheckIfMasterIsLost(masterID string, p peers.PeerUpdate) bool {
	for _, lostPeer := range p.Lost {
		if lostPeer == masterID {
			fmt.Printf("Master is lost!\n")
			return true
		}
	}
	return false
} */

/**
 * @func ChechIfMasterIsEmpty checks if the masterID is empty. Returns true if it is.
 *
 */
/* func ChechIfMasterIsEmpty(masterID string) bool {
	return masterID == ""
}
*/
/**
 * @func PrintMaster printes the master.
 *
 */
func PrintMaster(masterId string) {
	fmt.Printf("Master: %s\n", masterId)
}

// Might have to change to check Ip adress aswell to check if peer is master
func CheckIfYouAreMaster(masterProcessId string, processId int) bool {
	masterId_int, err := strconv.Atoi(masterProcessId)
	if err != nil {
		fmt.Printf("[error] Error converting masterId to int %v\n", err)
		//return
	}
	if processId == masterId_int {
		return true
	} else {
		return false
	}
}

/**
 * @func
 *
 */
func GetIdOfNewMaster(peerUpdateRx chan peers.PeerUpdate, sendToMasterCh chan string) {
	master := ""
	for {
		select {
		case p := <-peerUpdateRx:
			peers.PrintUpdatedPeers(p)
			fmt.Println()
			if p.Master != master {
				master = p.Master
				sendToMasterCh <- p.Master
			}
		}
	}
}

// Unsure if the two functions above should be in the master-module or not. Maybe just merge master and slave modules together.

func SendAndReceiveToSlaves(masterConnCh <-chan net.Conn, connectionsCh <-chan map[string]net.Conn,
	sendMapToSlavesCh <-chan map[string]elevator.Elevator, getElevFromSlave chan elevator.Elevator) {
	var connections map[string]net.Conn

	// fmt.Println("[SendAndReceiveToSlaves] kommet meg inni, og skal gå inn i for-loopen")
	for {
		select {
		case c := <-connectionsCh:
			// fmt.Println("[SendAndReceiveToSlaves] mottok en mappet av connections på connectionsCh")
			connections = c
			//fmt.Println("Send and receive sin conn liste", c)
			// fmt.Printf("[SendAndReceiveToSlaves] slik ser mappet ut: %v\n", connections)
			fmt.Println()
		case c := <-masterConnCh:
			// fmt.Println("[SendAndReceiveToSlaves] mottok en masterConn på masterConnCh")
			// fmt.Println("[SendAndReceiveToSlaves] går inn i Receive")
			go tcp.Receive(c, getElevFromSlave)
			// fmt.Println("[SendAndReceiveToSlaves] gått forbi Receive")
		case c := <-sendMapToSlavesCh:
			// fmt.Println("[SendAndReceiveToSlaves] mottok et map av elevs på sendMapToSlavesCh")
			// fmt.Println("[SendAndReceiveToSlaves] går igang med å sende til alle slavene")
			for _, v := range connections {
				// fmt.Println("[SendAndReceiveToSlaves] i for-loopen for å iterere gjennom connections-mappet")
				// fmt.Println("[SendAndReceiveToSlaves] for så å sende conn og map i Transmit")
				tcp.Transmit(v, c) // kanskje ikke goroutine?
				// fmt.Println("[SendAndReceiveToSlaves] iterasjon ferdig")
			}
		}
	}
}
