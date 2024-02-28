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
func DetermineMaster(p peers.PeerUpdate, masterID *string) {
	*masterID = p.Peers[0]
}

/**
 * @func ChechIfMasterExists checks whenever master exists or not and returns a bool.
 *
 */
func ChechIfMasterExists(masterID string) bool {
	return masterID != ""
}

/**
 * @func PrintMaster printes the master.
 *
 */
func PrintMaster(masterID string) {
	fmt.Printf("%s is Master\n", masterID)
}


