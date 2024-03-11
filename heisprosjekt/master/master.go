package master

import (
	"Heis/network/peers"
	"fmt"
	"strconv"
	"Heis/network/tcp"
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


func GetMaster(port string, masterCh chan string) {
	for {
		m := <-masterCh
		tcp.TransmitConn(port, m)
	}
}
