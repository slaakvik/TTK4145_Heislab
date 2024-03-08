package main

import (
	"Heis/network/tcp"
	"flag"
	"fmt"
	"net"
	"time"
	"Heis/network/peers"
	"Heis/network/localip"
	"os"
)

var masterPort = "8080"
var slavePort = "8081"


// this was just to make main smaller. Just temporary
func CheckId(id string) string {
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		return id
	}
	return id
} 

func main() {

	var id string
	id = CheckId(id) //set the id to default if id is not provided at launch

	var connections map[string]net.Conn

	connectionsCh := make(chan map[string]net.Conn)
	peerEnableTx := make(chan bool)
	peerUpdateRx := make(chan peers.PeerUpdate)

	go peers.Transmitter(15647, id, peerEnableTx)
	go peers.Receiver(15647, peerUpdateRx)

	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

	fmt.Printf("Dette er min ID: %v\n",id)

	if isMaster {
		go tcp.ReceiveConn(masterPort, connections, connectionsCh)
		for {
			select {
			case c := <-connectionsCh:
				fmt.Println("Updated connections:", c)
			}
		}
	} else {
		tcpConn, err := tcp.TransmitConn(masterPort, id)
		fmt.Printf("Her1\n")
		if err != nil {
			fmt.Printf("[error] Failed to Dial: %v\n", err)
			return
		}
		tcpConn.Write([]byte("ACK - Dette skriver jeg til Master\n"))
		
	}
	for{
		fmt.Printf("Nå er jeg i en uendelig loop\n")
		time.Sleep(3*time.Second)

	}
}

