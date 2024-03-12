package main

import (
	"Heis/elevator"
	"Heis/master"
	"Heis/network/establish_connection"
	"Heis/network/localip"
	"Heis/network/peers"
	"Heis/network/tcp"
	"flag"
	"fmt"
	"net"
	"os"
)

var masterPort = "8091"

//var slavePort = "8074"

// this was just to make main smaller. Just temporary
func CheckId() string {
	id := ""
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	return id
}

func main() {
	//kalosj! 500

	id := CheckId() //set the id to default if id is not provided at launch

	connectionsUpdateCh := make(chan map[string]net.Conn)
	connectionsCh := make(chan map[string]net.Conn)
	connectionCh := make(chan net.Conn)
	transmitCh := make(chan map[string]elevator.Elevator)

	masterUpdateCh := make(chan string)

	peerEnableTx := make(chan bool)
	peerUpdateRx := make(chan peers.PeerUpdate)

	go peers.Transmitter(15678, id, peerEnableTx)
	go peers.Receiver(15678, peerUpdateRx)

	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

	fmt.Printf("Dette er min ID: %v\n", id)

	if isMaster {
		go establish_connection.ReceiveConn(id, masterPort, connectionCh, connectionsUpdateCh)
		go establish_connection.AddConnections(id, connectionsUpdateCh, peerUpdateRx, connectionsCh)
		go tcp.SendAndReceive(connectionCh, connectionsCh, transmitCh)

	} else {
		//go tcp.ReceiveConn(slavePort, connectionsUpdateCh)
		go master.GetIdOfNewMaster(peerUpdateRx, masterUpdateCh)
		go master.NotifyMaster(masterPort, id, masterUpdateCh)
		// tcpConn, err := tcp.TransmitConn(masterPort, id)
		// fmt.Printf("Her1\n")
		// if err != nil {
		// 	fmt.Printf("[error] Failed to Dial: %v\n", err)
		// 	return
		// }
		// fmt.Printf("Dette er connen: %v\n", tcpConn)

		//tcpConn.Write([]byte("ACK - Dette skriver jeg til Master\n"))

		// for {
		// 	fmt.Printf("Nå er jeg i en uendelig loop\n")
		// 	time.Sleep(3 * time.Second)
		// }
	}
	select {}
}
