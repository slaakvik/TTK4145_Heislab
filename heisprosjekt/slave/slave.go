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
func NotifyMaster(port string, id string, sendMasterIdToNotifyMasterCh chan string, sendMasterIdToGetNotifyFromMaster chan string,
	sendElevToMaster chan elevator.Elevator, slaveConnCh chan<- net.Conn, connEstablished chan struct{}) {
	fmt.Println("[NotifyMaster] akkurat kommet inn")
	var elev elevator.Elevator
	var slaveConn net.Conn = nil // ??
	var err error                // ??
	fmt.Println("[NotifyMaster] går inn i for-loopen")
	for {
		select {
		case c := <-sendMasterIdToNotifyMasterCh:
			if id != c {
				fmt.Println("[NotifyMaster] Nå har jeg mottatt master på sendMasterIdToNotifyMaster, jeg er ikke master")
				//time.Sleep(1 * time.Millisecond)
				masterIp := peers.ExtractIpFromPeer(c)
				fmt.Println("[NotifyMaster] Master IP:", masterIp)
				fmt.Printf("[NotifyMaster] slaveConn før: %v\n", slaveConn)
				slaveConn, err = establish_connection.TransmitConn(port, id, masterIp)
				fmt.Printf("[NotifyMaster] slaveConn etter: %v\n", slaveConn)

				fmt.Println("[NotifyMaster] lager en eksempelheis som skal sendes")
				elevator := elevator.InitElev()
				fmt.Println("[NotifyMaster] Nå er jeg i ferd med å sende heisen til Transmit")
				tcp.Transmit(slaveConn, elevator)
				fmt.Println("[NotifyMaster] Nå skal heisen være sendt")

				//if slaveConn != nil {
				fmt.Println("[NotifyMaster] sjekkes om connEstablished kan closes")
				/* if _, ok := <-connEstablished; ok {
					close(connEstablished)
					fmt.Println("[NotifyMaster] closeEstablished kunne closes")
				} */
				fmt.Println("[NotifyMaster] kommet meg forbi closeEstablished ")

				//close(connEstablished)
				fmt.Println("[NotifyMaster] skal sende master-id på sendMasterIdToGetNotifyFromMaster")
				sendMasterIdToGetNotifyFromMaster <- c
				fmt.Println("[NotifyMaster] skal sende slaveConn på slaveConnCh")
				slaveConnCh <- slaveConn
				//}
				//fmt.Printf("Her1\n")
				if err != nil {
					fmt.Printf("[error] Failed to Dial: %v\n", err)
					return
				}
				fmt.Printf("[NotifyMaster] Dette er slaveConn som ble sendt: %v\n", slaveConn)
				// we need to send this SlaveConn to a chennel so a goroutine can start a receiver for the slave

			}
		case c := <-sendElevToMaster:
			//if slaveConn != nil {
			fmt.Println("[NotifyMaster] mottok en heis på sendElevToMaster")
			elev = c
			fmt.Println("[NotifyMaster] skal sende heisen gjennom Transmit")
			tcp.Transmit(slaveConn, elev)
			fmt.Println("[NotifyMaster] Nå skal heisen være sendt til master")

			//}

		}
	}
}

/**
 * @func slave tries to connect to master
 */
func GetNotifyFromMaster(id string, slaveConnCh <-chan net.Conn, sendMasterIdToGetNotifyFromMaster chan string,
	receiveMapFromMasterCh chan map[string]elevator.Elevator) {
	fmt.Println("[GetNotifyFromMaster] nå er jeg inni")
	masterId := ""
	var slaveConn net.Conn
	for {
		fmt.Println("[GetNotifyFromMaster] inni for-loopen")
		if masterId != id {
			fmt.Println("[GetNotifyFromMaster] dette kommer fordi jeg ikke er master")
			select {
			case c := <-sendMasterIdToGetNotifyFromMaster:
				fmt.Println("[GetNotifyFromMaster] mottok master på sendMasterIdToGetNotifyFromMaster")
				masterId = c
			case c := <-slaveConnCh:
				fmt.Printf("[GetNotifyFromMaster] mottok en slaveConn på slaveConnCh: %v\n", c)
				slaveConn = c
				fmt.Println("[GetNotifyFromMaster] starter en ReceiveHandler for slaven")
				go tcp.ReceiveHandler(slaveConn, receiveMapFromMasterCh /* kanaler som slave lytter på*/)
				fmt.Println("[GetNotifyFromMaster] nå skal ReceiveHandler være starta")
				//default: ??
			}
		}
	}
}
