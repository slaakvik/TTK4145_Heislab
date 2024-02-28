package main

import (
	/*
		"Heis/driver-go/elevio"
		"Heis/elevator"
		"Heis/fsm"
	*/
	"Heis/master"
	"Heis/slave"
	"Heis/network/localip"
	"Heis/network/peers"
	"flag"
	"fmt"
	"os"
)

func main() {
	var masterID string = ""
	var slavesID []string

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerEnableTx := make(chan bool)
	peerUpdateRx := make(chan peers.PeerUpdate)

	go peers.Transmitter(15647, id, peerEnableTx)
	go peers.Receiver(15647, peerUpdateRx)

	for {
		select {
		case p := <-peerUpdateRx:
			peers.PrintUpdatedPeers(p)
			if master.ChechIfMasterExists(masterID) {
				slave.MakeSlaves(p, &slavesID)
				slave.PrintSlaves(slavesID)
				master.PrintMaster(masterID)
				// slave functions
			} else {
				master.DetermineMaster(p, &masterID)
				master.PrintMaster(masterID)
			}
		}
	}
}
