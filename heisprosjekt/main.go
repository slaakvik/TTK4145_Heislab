package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/fsm"

	//"Heis/master"
	"Heis/timer"

	// "Heis/network/bcast"
	"Heis/network/establish_connection"
	"Heis/network/netfuncs"
	"Heis/network/peers"
	"Heis/network/tcp"
	"Heis/slave"
	"flag"
	"fmt"
	"net"
	// "Heis/network/tcp"
)

//_________________________________________________________________________________________________

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

//_________________________________________________________________________________________________

func main() {

	//_________________________________________________________________________________________________

	//Heisann35!
	masterPort := "8070"
	slavePort := "8090"

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	//elevio.Init("localhost:15657", _numFloors)
	elevio.Init("localhost:15668", _numFloors)

	// Master false by default
	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

	// Elevator initialization
	eleviId := netfuncs.InitNet()

	// Channels to update the the FSM-functions:
	newOrderCh := make(chan map[string]elevator.Elevator)
	elevUpdateRealtimeCh := make(chan elevator.Elevator)

	// Channels for door-timer and light:
	doorTimerCh := make(chan bool)
	timedOut := make(chan int)
	lightsCh := make(chan int)

	// Channels for inputs:
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// Master channels:
	masterConnCh := make(chan net.Conn)
	connectionsCh := make(chan map[string]net.Conn)
	sendMasterIdToReceive := make(chan string)
	sendMasterIdToNotifyMaster := make(chan string)
	sendMasterIdToGetNotifyFromMaster := make(chan string)
	sendMapToSlavesCh := make(chan map[string]elevator.Elevator)
	getElevFromSlaveRx := make(chan elevator.Elevator)

	// Slave channels:
	slaveConnCh := make(chan net.Conn)
	sendMyselfToMasterTx := make(chan elevator.Elevator)
	receiveMapFromMasterCh := make(chan map[string]elevator.Elevator)

	// Blocking channels:
	connEstablishedForSlave := make(chan struct{})
	ListenAccepted := make(chan struct{})

	// Channels for Heartbeat
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	// Goroutines for inputs
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Goroutines for Heartbeat
	go peers.Transmitter(15623, eleviId, peerTxEnable)
	go peers.Receiver(15623, peerUpdateCh)
	go peers.GetPeerUpdate(peerUpdateCh, sendMasterIdToReceive, sendMasterIdToNotifyMaster)

	fmt.Println("[main] Nå er jeg på vei inn i receiveConn")
	go establish_connection.ReceiveConn(eleviId, masterPort, masterConnCh, connectionsCh,
		sendMasterIdToReceive, ListenAccepted)
	fmt.Println("[main] Nå har jeg kommet meg forbi ReceiveConn")
	go slave.NotifyMaster(masterPort, eleviId, sendMasterIdToNotifyMaster, sendMasterIdToGetNotifyFromMaster,
		sendMyselfToMasterTx, slaveConnCh, connEstablishedForSlave)
	fmt.Println("[main] Nå har jeg nådd sperren for !isMaster")

	/* if !isMaster {
		<-connEstablishedForSlave
	} */

	fmt.Println("[main] Started!")
	go fsm.ButtonsAndRequests(masterPort, eleviId, isMaster, elevUpdateRealtimeCh,
		drv_buttons /*elevatorTx, */, sendMapToSlavesCh, getElevFromSlaveRx, receiveMapFromMasterCh,
		newOrderCh, lightsCh, sendMyselfToMasterTx)

	go fsm.FloorObstrStop(masterPort, isMaster, eleviId, elevUpdateRealtimeCh, drv_floors, /*elevatorTx,*/
		newOrderCh, doorTimerCh, timedOut, lightsCh, sendMyselfToMasterTx)

	//go fsm.FSM(drv_buttons, drv_floors, drv_stop, drv_obstr, eleviId, isMaster, ElevatorTx, CostTx, masterOrders, ElevatorRx, CostRx)
	//go netfuncs.Network_FSM(peerUpdateCh)
	go timer.Timer(doorTimerCh, timedOut)

	// Slave
	go slave.GetNotifyFromMaster(eleviId, slaveConnCh, sendMasterIdToGetNotifyFromMaster, receiveMapFromMasterCh)

	// Master
	fmt.Println("[main] Nå har vi nådd sperren for master")
	/* if isMaster {
		fmt.Println("[main] Jeg er master, og har nådd sperren")
		<-ListenAccepted
		fmt.Println("[main] Jeg er master, og har kommet meg")

	} */
	fmt.Println("[main] Nå er vi kommet forbi sperren for master")
	go tcp.SendAndReceive(masterConnCh, connectionsCh, sendMapToSlavesCh, getElevFromSlaveRx)
	// go tcp.Receive(masterPort, eleviId, elevatorRx)

	select {}
}
