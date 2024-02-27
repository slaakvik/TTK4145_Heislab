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
 * @struct
 *
 */

var MasterID string

/**
 * @func DetermineMaster is not assigning or making a master, it only determines who of the peers should be assigned
 *
 */
func DetermineMaster(p peers.PeerUpdate) int {
	var lowestPeer int = 0
	for i, v := range p.Peers {
		if v < p.Peers[lowestPeer] {
			lowestPeer = i
		}
	}
	return lowestPeer
}

/**
 * @func
 *
 */
func MakeMaster(p peers.PeerUpdate) {
	var lowestPeer int
	var newMasterID string

	lowestPeer = DetermineMaster(p)
	newMasterID = p.Peers[lowestPeer]
	fmt.Printf("%s is Master\n", newMasterID)
}



