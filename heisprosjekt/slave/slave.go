package slave

import (
	"Heis/elevator"
	"Heis/network/establish_connection"
	"Heis/network/tcp"
	"fmt"
	"net"
)

/*
---------------------------------------------------------------------------
--! @file
--! @brief This file contain functions and variables regarding the slaves
---------------------------------------------------------------------------
*/

/**
 * @func slave tries to connect to master
 */
func NotifyMaster(port string, id string, sendMasterIdToNotifyMaster chan string, sendMasterIdToGetNotifyFromMaster chan string, slaveConnCh chan<- net.Conn, connEstablished chan struct{}) {
	// var elev elevator.Elevator
	var slaveConn net.Conn = nil // ??
	var err error                // ??
	// for {
	// 	select {
	// case m := <-sendMasterIdToNotifyMaster:
	fmt.Println("Nå har jeg mottatt master, jeg er ikke master")
	//time.Sleep(1 * time.Millisecond)
	// masterIp := peers.ExtractIpFromPeer(m)
	masterIp := "localhost"
	fmt.Println("Master IP:", masterIp)
	fmt.Printf("slaveConn før: %v\n", slaveConn)
	slaveConn, err = establish_connection.TransmitConn(port, id, masterIp)
	fmt.Printf("slaveConn etter: %v\n", slaveConn)
	//if slaveConn != nil {
	// fmt.Println("closer nå jeg (jeg er slave)")
	//close(connEstablished)
	// fmt.Println("nå har jeg closet")
	// fmt.Println("Slave prøver å sende masterId")
	// sendMasterIdToGetNotifyFromMaster <- "p-localhost-p"
	slaveConnCh <- slaveConn
	//}
	fmt.Printf("Her1\n")
	if err != nil {
		fmt.Printf("[error] Failed to Dial: %v\n", err)
		return
	}
	fmt.Printf("Dette er connen: %v\n", slaveConn)
	if _, ok := <-connEstablished; ok {
		close(connEstablished)
	}
	// we need to send this SlaveConn to a chennel so a goroutine can start a receiver for the slave

	// case m := <-sendElevToMaster:
	// 	//if slaveConn != nil {
	// 	fmt.Println("skal sende heis til master")
	// 	elev = m
	// 	tcp.Transmit(slaveConn, elev)
	// 	fmt.Println("Nå har jeg sendt heis til master")

	//}

}

// 	}
// }

/**
 * @func slave tries to connect to master
 */
func GetNotifyFromMaster(id string, slaveConnCh <-chan net.Conn, sendMasterIdToGetNotifyFromMaster chan string,
	receiveMapFromMasterCh chan map[string]elevator.Elevator, sendElevToMaster chan elevator.Elevator) {
	// masterId := ""
	var slaveConn net.Conn
	for {
		// if masterId != id {
		// 	select {
		// 	case c := <-sendMasterIdToGetNotifyFromMaster:
		// 		masterId = c
		select {
		case c := <-slaveConnCh:
			slaveConn = c
			go tcp.ReceiveHandler(slaveConn, receiveMapFromMasterCh /* kanaler som slave lytter på*/)
			// default:
		case c := <-sendElevToMaster:
			tcp.Transmit(slaveConn, c)

		}
	}
}

// }
