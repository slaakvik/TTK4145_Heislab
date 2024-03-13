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

	//Heisann34!

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	// elevio.Init("localhost:15657", _numFloors)
	elevio.Init("localhost:15668", _numFloors)

	masterPort := "8070"
	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

	//Initialiserer en heisstruct
	//elev := elevator.InitElev()
	// elevUpdateBtnAndOrdersCh := make(chan elevator.Elevator)
	//Buffer:
	// newOrderCh := make(chan map[string]elevator.Elevator, 10)
	// elevUpdateRealtimeCh := make(chan elevator.Elevator, 10)
	newOrderCh := make(chan map[string]elevator.Elevator)
	elevUpdateRealtimeCh := make(chan elevator.Elevator)

	doorTimerCh := make(chan bool)
	timedOut := make(chan int)
	lightsCh := make(chan int)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Master channels:
	sendMasterIdToReceive := make(chan string)
	sendMasterIdToNotifyMaster := make(chan string)
	sendMasterIdToGetNotifyFromMaster := make(chan string)

	connectionsCh := make(chan map[string]net.Conn)
	masterConnCh := make(chan net.Conn)

	// Slave channels:
	slaveConnCh := make(chan net.Conn)

	//networkchannels:
	// Heartbeat
	eleviId := netfuncs.InitNet()
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15623, eleviId, peerTxEnable)
	go peers.Receiver(15623, peerUpdateCh)
	go peers.GetPeerUpdate(peerUpdateCh, sendMasterIdToReceive, sendMasterIdToNotifyMaster)

	//broadcast
	// elevatorTx := make(chan elevator.Elevator)
	getElevFromSlaveRx := make(chan elevator.Elevator)
	// go bcast.Transmitter(16523, elevatorTx)
	// go bcast.Receiver(16523, elevatorRx)

	sendMyselfToMasterTx := make(chan elevator.Elevator)
	//go netfuncs.Bcast_message(ElevatorTx, elev, eleviId)
	//send cost func result to network
	sendMapToSlavesCh := make(chan map[string]elevator.Elevator)
	mapOfElevsRx := make(chan map[string]elevator.Elevator)
	// go bcast.Transmitter(16524, mapOfElevsTx)
	// go bcast.Receiver(16524, mapOfElevsRx)

	//Master slave
	// isMaster := true
	connEstablishedForSlave := make(chan struct{})
	ListenAccepted := make(chan struct{})

	go establish_connection.ReceiveConn(eleviId, masterPort, masterConnCh, connectionsCh, sendMasterIdToReceive, ListenAccepted)

	go slave.NotifyMaster(masterPort, eleviId, sendMasterIdToNotifyMaster, sendMasterIdToGetNotifyFromMaster, sendMyselfToMasterTx, slaveConnCh, connEstablishedForSlave)
	fmt.Println("Nå har jeg nådd sperren")

	if !isMaster {
		<-connEstablishedForSlave
	}

	fmt.Printf("Started!\n")
	go fsm.ButtonsAndRequests(masterPort, eleviId, isMaster, elevUpdateRealtimeCh, drv_buttons /*elevatorTx, */, sendMapToSlavesCh, getElevFromSlaveRx, mapOfElevsRx, newOrderCh, lightsCh, sendMyselfToMasterTx)
	go fsm.FloorObstrStop(masterPort, isMaster, eleviId, elevUpdateRealtimeCh, drv_floors /*elevatorTx,*/, newOrderCh, doorTimerCh, timedOut, lightsCh, sendMyselfToMasterTx)

	//go fsm.FSM(drv_buttons, drv_floors, drv_stop, drv_obstr, eleviId, isMaster, ElevatorTx, CostTx, masterOrders, ElevatorRx, CostRx)
	//go netfuncs.Network_FSM(peerUpdateCh)
	go timer.Timer(doorTimerCh, timedOut)

	// Slave
	go slave.GetNotifyFromMaster(eleviId, slaveConnCh, sendMasterIdToGetNotifyFromMaster, mapOfElevsRx)

	// Master
	fmt.Println("Nå har vi nådd sperren for master")

	if isMaster {
		<-ListenAccepted
	}
	fmt.Println("Nå er vi kommet forbi sperren for master")
	go tcp.SendAndReceive(masterConnCh, connectionsCh, mapOfElevsCh, getElevFromSlaveRx)
	// go tcp.Receive(masterPort, eleviId, elevatorRx)

	select {}
}
