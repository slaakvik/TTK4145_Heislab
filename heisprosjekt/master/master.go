package master

import (
	"Heis/elevator"
	"Heis/network/peers"
	"Heis/network/tcp"
	"fmt"
	"net"
	"sync"
)

/*
--------------------------------------------------------------------------
--! @file
--! @brief This file contain functions and variables regarding the master
--------------------------------------------------------------------------
*/


var mutex sync.Mutex




/**
 * @func ChooseMasterIndex is not assigning or making a master, it only chooses who of the peers should be assigned, which is the peer with the lowest processID
 *
 */
/* func ChooseMasterIndex(peersProcessId []string) (int, error) {
	var SmallestProcessId, err = strconv.Atoi(peersProcessId[0])
	if err != nil {
		fmt.Printf("[error]: error converting string to int: %s\n", err)
		return 0, err
	}
	var SmallestProcessIdIndex = 0
	for i, peer_str := range peersProcessId {
		peer_int, err := strconv.Atoi(peer_str)
		if err != nil {
			fmt.Printf("[error]: error converting string to int: %s\n", err)
			return 0, err
		}
		if peer_int < SmallestProcessId {
			SmallestProcessId = peer_int
			SmallestProcessIdIndex = i
		}
	}
	return SmallestProcessIdIndex, err
} */

/* func ChooseMaster(masterIdIndex int, p peers.PeerUpdate) string {
	masterId := p.Peers[masterIdIndex]
	return masterId
} */

/* func CheckIfYouAreMaster(masterProcessId string, processId int) bool {
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
*/

/* func GetIdOfNewMaster(peerUpdateRx chan peers.PeerUpdate, sendToMasterCh chan string) {
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
*/

/* func SendAndReceiveToSlaves(masterConnCh <-chan net.Conn, connectionsCh <-chan map[string]net.Conn,
	sendMapToSlavesCh <-chan map[string]elevator.Elevator, getElevFromSlave chan elevator.Elevator) {
	var connections map[string]net.Conn

	for {
		select {
		case c := <-connectionsCh:
			connections = c
		case c := <-masterConnCh:
			go tcp.Receive(c, getElevFromSlave)
		case c := <-sendMapToSlavesCh:
			for _, v := range connections {
				tcp.Transmit(v, c)
			}
		}
	}
} */


func SendAndReceiveToSlaves(id string, peerCh chan peers.PeerUpdate, masterConnCh <-chan net.Conn, connectionsCh <-chan map[string]net.Conn,
	sendMapToSlavesCh <-chan map[string]elevator.Elevator, getElevFromSlave chan elevator.Elevator) {
	var connections map[string]net.Conn
	for {
		select {
		case c := <-peerCh:
			peers.PrintUpdatedPeers(c)
			if len(c.Lost) != 0 {
				mutex.Lock()
				for i := 0; i < len(c.Lost); i++ {
					for k := range connections {
						if k == c.Lost[i] {
							delete(connections, k)
						}
					}
				}
				mutex.Unlock()
			}
		case c := <-connectionsCh:
			mutex.Lock()
			connections = c
			fmt.Println()
			mutex.Unlock()
		case c := <-masterConnCh:
			go tcp.Receive(c, getElevFromSlave)
		case c := <-sendMapToSlavesCh:
			mutex.Lock()
			for _, v := range connections {
				tcp.Transmit(v, c) 
			}
			mutex.Unlock()
		}
	}
}
