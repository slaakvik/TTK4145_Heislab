package master

import (
	"Heis/elevator"
	"Heis/network/establish_connection"
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

/**
 * @func
 *
 */
func NotifyMaster(port string, id string, SendMasterCh chan string, sendElevToMaster chan elevator.Elevator) {
	var elev elevator.Elevator
	var tcpConn net.Conn
	for {
		select {
		case m := <-SendMasterCh:
			if id != m {
				fmt.Println("--------------[jajaja]---------------")
				//time.Sleep(1 * time.Millisecond)
				masterIp := peers.ExtractIpFromPeer(m)
				fmt.Println("Master IP:", masterIp)
				tcpConn, err := establish_connection.TransmitConn(port, id, masterIp)
				//fmt.Printf("Her1\n")
				if err != nil {
					fmt.Printf("[error] Failed to Dial: %v\n", err)
					return
				}
				fmt.Printf("Dette er connen: %v\n", tcpConn)
			}
		case m := <-sendElevToMaster:
			if tcpConn != nil {
				fmt.Println("--------------[jajaja2]---------------")
				elev = m
				fmt.Printf("SJå på denne da: %v\n", tcpConn)
				tcp.Transmit(tcpConn, elev)
				fmt.Println("--------------[jajaja3]---------------")

			}

		}
	}
}
