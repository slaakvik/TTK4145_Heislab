package master

import (
	"Heis/network/peers"
	"fmt"
	"strconv"
)

/*
-------------------------------------------------------
--! @file
--! @brief This file contain functions and variables regarding the master
-------------------------------------------------------
*/

/**
 * @func ChooseMaster is not assigning or making a master, it only chooses who of the peers should be assigned, which is the peer with the lowest processID
 *
 */
func ChooseMaster(peersProcessId []string) (string, error) {
	var SmallestProcessId, err = strconv.Atoi(peersProcessId[0])
	if err != nil {
		fmt.Printf("[error]: error converting string to int: %s\n", err)
		return "", err
	}
	var SmallestProcessIdIndex = 0
	for i, peer_str := range peersProcessId {
		peer_int, err := strconv.Atoi(peer_str)
		if err != nil {
			fmt.Printf("[error]: error converting string to int: %s\n", err)
			return "", err
		}
		if peer_int < SmallestProcessId {
			SmallestProcessId = peer_int
			SmallestProcessIdIndex = i
		}
	}
	return peersProcessId[SmallestProcessIdIndex], nil
}

/**
 * @func ChechIfMasterIsLost checks if master is lost. Return true if it is lost.
 *
 */
func CheckIfMasterIsLost(masterID string, p peers.PeerUpdate) bool {
	for _, lostPeer := range p.Lost {
		if lostPeer == masterID {
			fmt.Printf("Master is lost!\n")
			return true
		}
	}
	return false
}

/**
 * @func ChechIfMasterIsEmpty checks if the masterID is empty. Returns true if it is.
 *
 */
func ChechIfMasterIsEmpty(masterID string) bool {
	return masterID == ""
}

/**
 * @func PrintMaster printes the master.
 *
 */
func PrintMaster(masterID string) {
	fmt.Printf("Master: %s\n", masterID)
}
