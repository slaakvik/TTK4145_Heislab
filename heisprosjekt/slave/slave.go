package slave

import (
	"fmt"
	"Heis/network/peers"
	"net"
	
)

/**
 * This file contain functions and variables regarding the slave
 */


/**
 * @var
 * slice consisting of the slaves on the network.
 */
var SlavesID []string


//slavene er SlavesID = .p.Peers[1,:]
