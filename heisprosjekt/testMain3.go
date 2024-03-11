package main

import (
	"Heis/master"
	"Heis/network/peers"
	"Heis/network/tcp"
	"flag"
	"fmt"
)

const (
	MasterPort = "8039"
)

func PeerLytter(peerCh <-chan tcp.Peer, guttaneCh <-chan map[string]tcp.Peer, p chan peers.PeerUpdate) { //burde vel være read only kanal? siden den bare skal motta på
	//connections := make(map[string]net.Conn)
	for {
		select {
		case a := <-peerCh:
			//conne
			fmt.Printf("Ludvig ser slik ut han:\n")
			fmt.Printf("  id:    %v\n", a.Id)
			fmt.Printf("  Floors:      %v\n", a.Floors)
			fmt.Printf("  Color:     %v\n", a.Color)
			//case a:= <-connectionsCh:

		case a := <-guttaneCh:
			//fmt.Println(a)
			for k, v := range a {
				println()
				fmt.Printf("Nøkkel: %v\n", k)
				fmt.Printf("Her er structet: %v\n", v)
				println()
			}
		case a := <-p:
			peers.PrintUpdatedPeers(a)
			println()

		}
	}

}

func main() {
	//kalosj! 500

	id := peers.CheckId() //set the id to default if id is not provided at launch
	// lager et map med ulike Peers

	peerLudvig := tcp.Peer{
		Id:     id,
		Floors: 10,
		Color:  "blue",
	}

	peerJonathan := tcp.Peer{
		Id:     "Jonathan",
		Floors: 12,
		Color:  "green",
	}

	peerOlav := tcp.Peer{
		Id:     "Olav",
		Floors: 12,
		Color:  "red",
	}

	//peerCh := make(chan tcp.Peer)
	guttaneCh := make(chan map[string]tcp.Peer)

	guttane := make(map[string]tcp.Peer)
	guttane[peerLudvig.Id] = peerLudvig
	guttane[peerJonathan.Id] = peerJonathan
	guttane[peerOlav.Id] = peerOlav

	//--------

	//connectionsUpdateCh := make(chan map[string]net.Conn)
	//connectionsCh := make(chan map[string]net.Conn)

	//masterUpdateCh := make(chan string)
	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

	peerEnableTx := make(chan bool)
	peerUpdateRx := make(chan peers.PeerUpdate)
	masterCh := make(chan string)

	//fmt.Printf("Dette er isMaster: %v\n", isMaster)
	go peers.Transmitter(15679, id, isMaster, peerEnableTx)
	go peers.Receiver(15679, peerUpdateRx)
	go peers.GetPeerUpdate(peerUpdateRx, masterCh)

	fmt.Printf("Dette er min ID: %v\n", id)

	
	if isMaster {
		go tcp.Receive(MasterPort, id, guttaneCh)
		/* go tcp.ReceiveConn(id, masterPort, connectionsUpdateCh)
		go tcp.AddConnections(id, connectionsUpdateCh, peerUpdateRx, connectionsCh) */
		// CONNECTIONS map
		//go tcp.Receive(MasterPort, id, guttaneCh)
		//go PeerLytter(peerCh, guttaneCh, peerUpdateRx)
		//go tcp.SendAndReceive(connectionsCh, isMaster)

	} else {
		go master.GetMaster(MasterPort, masterCh)
		//go PeerLytter(peerCh, guttaneCh, peerUpdateRx)
		//go tcp.ReceiveConn(slavePort, connectionsUpdateCh)
		//go tcp.GetIdOfNewMaster(peerUpdateRx, masterUpdateCh)
		//go tcp.NotifyMaster(masterPort, id, masterUpdateCh)
		//tcpConn, _ := tcp.TransmitConn2(masterPort, id)
		//tcp.Transmiter(tcpConn, guttane)
		//fmt.Printf("TcpConn: %v\n", tcpConn)
		// if err != nil {
		// 	fmt.Printf("[error] Failed to Dial: %v\n", err)
		// 	return
		// }

	}
	select {}
}
