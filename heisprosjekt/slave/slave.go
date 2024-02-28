package slave

import (
	"Heis/network/peers"
	"fmt"


)

//"fmt"
//"Heis/network/peers"
//"net"

/**
 * This file contain functions and variables regarding the slave
 */



/**
 * @func MakeSlaves makes the slaves, which is the remainding elements in Peers (Master is first element)
 */
func MakeSlaves(p peers.PeerUpdate, slavesID *[]string){
	*slavesID = p.Peers[1:]
}


/**
 * @func PrintSlaves 
 */
func PrintSlaves(slavesID []string){
	fmt.Printf("The Slaves:\n")
	for _,slave := range slavesID{
		fmt.Printf("Slave:    %q\n", slave)
	}
	
}

 /* func PrintUpdatedPeers(p PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
} */