package main

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
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

	input := cost_fns.InputToCost()

	cost_fns.HRA_funcs(input)

	//_________________________________________________________________________________________________

	//heisann, din gamle ørn12!

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)

	//Initialiserer en heisstruct
	elev := elevator.InitElev()

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
	ElevatorTx := make(chan netfuncs.HelloMsg)
	ElevatorRx := make(chan netfuncs.HelloMsg)
	go bcast.Transmitter(16523, ElevatorTx)
	go bcast.Receiver(16523, ElevatorRx)

	go netfuncs.Bcast_message(ElevatorTx, elev, eleviId)

	fmt.Printf("Started!\n")
	//elevator.Elevator_print(elev)

	if elevio.GetFloor() == -1 {
		elev = fsm.OnInitBetweenFloors(elev)
	}

	for {
		//elevator.Elevator_print(elev)
		//fmt.Print(eleviId)
		select {
		case a := <-drv_buttons:
			fmt.Printf("Button: %+v\n", a)
			fsm.OnRequestButtonPress(&elev, a.Floor, a.Button)

		case a := <-drv_floors: // a er etasjen heisen er i
			fmt.Printf("Floor: %+v\n", a)
			fsm.OnFloorArrival(&elev, a)

		case a := <-drv_stop:
			fmt.Printf("Stop button: %+v\n", a)
			// fsm.Fsm_onStopButtonPress(&elev)
			fsm.Stop_functionality(a, elev)

		case a := <-drv_obstr:
			fmt.Printf("Obstruction %+v\n", a)
			fsm.Obstruction_functionality(a, elev)

		case a := <-peerUpdateCh:
			//fmt.Printf("%+v\n", a)
			netfuncs.PrintPeerUpdate(a)
		case a := <-ElevatorRx:
			fmt.Printf("Received: %+v\n", a)

		}
	}
}
