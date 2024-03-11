package main

import (
	//"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/timer"

	// "Heis/elevator"
	"Heis/fsm"

	"Heis/network/bcast"
	"Heis/network/netfuncs"
	"Heis/network/peers"
	"fmt"
)

//_________________________________________________________________________________________________

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

//_________________________________________________________________________________________________

func main() {

	//_________________________________________________________________________________________________

	//heisann, din gamle ørn36!

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)

	//Initialiserer en heisstruct
	//elev := elevator.InitElev()
	// elevUpdateBtnAndOrdersCh := make(chan elevator.Elevator)
	newOrderCh := make(chan map[string]elevator.Elevator)
	elevUpdateRealtimeCh := make(chan elevator.Elevator)
	doorTimerCh := make(chan bool)
	timedOut := make(chan int)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//networkchannels:
	//Peers
	eleviId := netfuncs.InitNet()
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15623, eleviId, peerTxEnable)
	go peers.Receiver(15623, peerUpdateCh)
	//broadcast
	ElevatorTx := make(chan elevator.Elevator)
	ElevatorRx := make(chan elevator.Elevator)
	go bcast.Transmitter(16523, ElevatorTx)
	go bcast.Receiver(16523, ElevatorRx)

	//go netfuncs.Bcast_message(ElevatorTx, elev, eleviId)
	//send cost func result to network
	mapOfElevsTx := make(chan map[string]elevator.Elevator)
	mapOfElevsRx := make(chan map[string]elevator.Elevator)
	go bcast.Transmitter(16524, mapOfElevsTx)
	go bcast.Receiver(16524, mapOfElevsRx)
	//Master slave
	isMaster := true
	//masterOrders := make(chan map[string][][2]bool)

	fmt.Printf("Started!\n")
	//cost function test:
	// input := cost_fns.InputToCost(elev)
	// cost_fns.HRA_funcs(input)
	//elevator.Elevator_print(elev)

	// if elevio.GetFloor() == -1 {
	// 	elev = fsm.OnInitBetweenFloors(elev)
	// 	elevator.Elevator_print(elev)
	// }
	go fsm.ButtonsAndRequests(eleviId, isMaster, elevUpdateRealtimeCh, drv_buttons, ElevatorTx, mapOfElevsTx, ElevatorRx, mapOfElevsRx, newOrderCh)
	go fsm.FloorObstrStop(isMaster, eleviId, elevUpdateRealtimeCh, drv_floors, drv_stop, drv_obstr, ElevatorTx, newOrderCh, doorTimerCh, timedOut)

	//go fsm.FSM(drv_buttons, drv_floors, drv_stop, drv_obstr, eleviId, isMaster, ElevatorTx, CostTx, masterOrders, ElevatorRx, CostRx)
	go netfuncs.Network_FSM(peerUpdateCh)
	go timer.Timer(doorTimerCh, timedOut)
	select {}
}

// for {
// 	//elevator.Elevator_print(elev)
// 	//fmt.Print(eleviId)
// 	select {
// 	// case a := <-drv_buttons:
// 	// 	fmt.Printf("Button: %+v\n", a)
// 	// 	fsm.OnRequestButtonPress(&elev, a.Floor, a.Button)

// 	// case a := <-drv_floors: // a er etasjen heisen er i
// 	// 	fmt.Printf("Floor: %+v\n", a)
// 	// 	fsm.OnFloorArrival(&elev, a)

// 	// case a := <-drv_stop:
// 	// 	fmt.Printf("Stop button: %+v\n", a)
// 	// 	// fsm.Fsm_onStopButtonPress(&elev)
// 	// 	fsm.Stop_functionality(a, elev)

// 	// case a := <-drv_obstr:
// 	// 	fmt.Printf("Obstruction %+v\n", a)
// 	// 	fsm.Obstruction_functionality(a, elev)

// 	case a := <-peerUpdateCh:
// 		//fmt.Printf("%+v\n", a)
// 		netfuncs.PrintPeerUpdate(a)
// 	case a := <-ElevatorRx:
// 		fmt.Printf("Received: %+v\n", a)

// 	}
// }
