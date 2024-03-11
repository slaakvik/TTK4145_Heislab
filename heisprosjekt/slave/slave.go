package slave

import (
	"Heis/network/peers"
	"fmt"
	"strings"
)

/**
 * This file contain functions and variables regarding the slave
 */

/**
 * @func MakeSlaves makes the slaves, which is the remainding elements in Peers (Master is first element)
 */
func MakeSlaves(slavesID *[]string, p peers.PeerUpdate) {
	*slavesID = p.Peers[1:]
}

/**
 * @func PrintSlaves
 */
func PrintSlaves(slavesID []string) {
	fmt.Printf("Slaves:   %q\n", slavesID)
}



func ChooseSlaves(p peers.PeerUpdate, slavesId []string, masterId string) []string {
	for _, peer := range p.Peers{
		if peer != masterId {
			slavesId = append(slavesId, peer)
		}
	}
	return slavesId
}



func ExtractIpFromSlaves(slavesId []string, slavesIp []string) []string {
	//slavesIp = nil
	for _, slave := range slavesId {
		data := strings.Split(slave, "-")
		slavesIp = append(slavesIp, data[1])
		fmt.Printf("AAAAAAAA: %v\n", data[1])
	}
	return slavesIp
}