package slave

import (
	"Heis/elevator"
	"Heis/network/establish_connection"
	"Heis/network/peers"
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
func AlertMaster(port string, id string, masterIdToAlertMasterCh chan string, masterIdToSendAndReceiveToMasterCh chan string, slaveConnCh chan<- net.Conn, connEstablished chan struct{}) {
	fmt.Println("[AlertMaster] akkurat kommet inn")
	var slaveConn net.Conn = nil // ??
	var err error                // ??
	fmt.Println("[AlertMaster] går inn i for-loopen")
	//for {
	select {
	case c := <-masterIdToAlertMasterCh:
		if id != c {
			fmt.Println("[AlertMaster] Nå har jeg mottatt master på sendMasterIdToNotifyMaster, jeg er ikke master")
			//time.Sleep(1 * time.Millisecond)
			masterIp := peers.ExtractIpFromPeer(c)
			fmt.Println("[AlertMaster] Master IP:", masterIp)
			fmt.Printf("[AlertMaster] slaveConn før: %v\n", slaveConn)
			slaveConn, err = establish_connection.EstablishConnToMaster(port, id, masterIp)
			fmt.Printf("[AlertMaster] slaveConn etter: %v\n", slaveConn)

			fmt.Println("[AlertMaster] lager en eksempelheis som skal sendes")
			elevator := elevator.InitElev()
			fmt.Println("[AlertMaster] Nå er jeg i ferd med å sende heisen til Transmit")
			tcp.Transmit(slaveConn, elevator)
			fmt.Println("[AlertMaster] Nå skal heisen være sendt")

			//if slaveConn != nil {
			fmt.Println("[AlertMaster] sjekkes om connEstablished kan closes")
			/* if _, ok := <-connEstablished; ok {
				close(connEstablished)
				fmt.Println("[NotifyMaster] closeEstablished kunne closes")
			} */
			fmt.Println("[AlertMaster] kommet meg forbi closeEstablished ")

			//close(connEstablished)
			fmt.Println("[AlertMaster] skal sende master-id på sendMasterIdToGetNotifyFromMaster")
			masterIdToSendAndReceiveToMasterCh <- c
			fmt.Println("[AlertMaster] skal sende slaveConn på slaveConnCh")
			slaveConnCh <- slaveConn
			//}
			//fmt.Printf("Her1\n")
			if err != nil {
				fmt.Printf("[error] Failed to Dial: %v\n", err)
				return
			}
			fmt.Printf("[AlertMaster] Dette er slaveConn som ble sendt: %v\n", slaveConn)
			// we need to send this SlaveConn to a chennel so a goroutine can start a receiver for the slave
		}
	}
	//}
}

/**
 * @func slave tries to connect to master
 */
func SendAndReceiveToMaster(id string, slaveConnCh <-chan net.Conn, masterIdToSendAndReceiveToMasterCh chan string,
	receiveMapFromMasterCh chan map[string]elevator.Elevator, sendElevToMaster chan elevator.Elevator) {
	// fmt.Println("[SendAndReceiveToMaster] nå er jeg inni")
	var elev elevator.Elevator
	masterId := ""
	var slaveConn net.Conn
	for {
		// fmt.Println("[SendAndReceiveToMaster] inni for-loopen")
		if masterId != id {
			// fmt.Println("[SendAndReceiveToMaster] dette kommer fordi jeg ikke er master")
			select {
			case c := <-masterIdToSendAndReceiveToMasterCh:
				// fmt.Println("[SendAndReceiveToMaster] mottok master på sendMasterIdToGetNotifyFromMaster")
				masterId = c
			case c := <-slaveConnCh:
				// fmt.Printf("[SendAndReceiveToMaster] mottok en slaveConn på slaveConnCh: %v\n", c)
				slaveConn = c
				// fmt.Println("[SendAndReceiveToMaster] starter en ReceiveHandler for slaven")
				go tcp.Receive(slaveConn, receiveMapFromMasterCh /* kanaler som slave lytter på*/)
				// fmt.Println("[SendAndReceiveToMaster] nå skal ReceiveHandler være starta")
				//default: ??
			case c := <-sendElevToMaster:
				//if slaveConn != nil {
				// fmt.Println("[SendAndReceiveToMaster] mottok en heis på sendElevToMaster")
				elev = c
				// fmt.Println("[SendAndReceiveToMaster] skal sende heisen gjennom Transmit")
				tcp.Transmit(slaveConn, elev)
				// fmt.Println("[SendAndReceiveToMaster] Nå skal heisen være sendt til master")
			}
		}
	}
}
