package master

import (
	"Heis/network/peers"
	"fmt"
)

/*
-------------------------------------------------------
--! @file
--! @brief This file contain functions and variables regarding the master
-------------------------------------------------------
*/

/**
 * @func DetermineMaster is not assigning or making a master, it only determines who of the peers should be assigned
 *
 */
func DetermineMaster(masterID *string, p peers.PeerUpdate) {
	*masterID = p.Peers[0]
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
