package main

import (
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

	id := CheckId() //set the id to default if id is not provided at launch

	connectionsUpdateCh := make(chan map[string]net.Conn)

	connectionsCh := make(chan map[string]net.Conn)

	slaveTxCh := make(chan string)

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
		go tcp.ReceiveConn(id, masterPort, connectionsUpdateCh)
		go tcp.AddConnections(id, connectionsUpdateCh, peerUpdateRx, connectionsCh)
		go tcp.SendAndReceive(connectionsCh)

	} else {
		//go tcp.ReceiveConn(slavePort, connectionsUpdateCh)
		go tcp.SlavePeerUpdateRx(peerUpdateRx, slaveTxCh)
		go tcp.NotifyMaster(masterPort, id, slaveTxCh)
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
