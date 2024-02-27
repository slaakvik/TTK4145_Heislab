package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/fsm"
	"Heis/network/bcast"
	"Heis/network/netfuncs"
	"Heis/network/peers"
	"fmt"
)

func main() {
	//heisann, din gamle ørn4!

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
	go peers.Transmitter(15647, eleviId, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)
	//broadcast
	ElevatorTx := make(chan netfuncs.HelloMsg)
	ElevatorRx := make(chan netfuncs.HelloMsg)
	go bcast.Transmitter(16569, ElevatorTx)
	go bcast.Receiver(16569, ElevatorRx)

	go netfuncs.Bcast_message(ElevatorTx, eleviId)

	fmt.Printf("Started!\n")
	elevator.Elevator_print(elev)

	if elevio.GetFloor() == -1 {
		fsm.Fsm_onInitBetweenFloors(&elev)
	}

	for {
		elevator.Elevator_print(elev)

		select {
		case a := <-drv_buttons:
			//fmt.Printf("%+v\n", a)
			//elevio.SetButtonLamp(a.Button, a.Floor, true)  //Denne er vel litt for tidlig
			fsm.Fsm_onRequestButtonPress(&elev, a.Floor, a.Button)

		case a := <-drv_floors: // a er etasjen heisen er i
			//fmt.Printf("%+v\n", a)
			fsm.Fsm_onFloorArrival(&elev, a)

		case a := <-drv_stop:
			//fmt.Printf("%+v\n", a)
			// fsm.Fsm_onStopButtonPress(&elev)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetStopLamp(true)
			} else {
				elevio.SetMotorDirection(elev.Dirn)
				elevio.SetStopLamp(false)
			}
		case a := <-peerUpdateCh:
			//fmt.Printf("%+v\n", a)
			netfuncs.PrintPeerUpdate(&a)
		case a := <-ElevatorRx:
			fmt.Printf("Received: %+v\n", a)
			// Fikser disse funksjonene senere
			/*
				case a := <-drv_obstr:
					fmt.Printf("%+v\n", a)
					if a {
						elevio.SetMotorDirection(elevio.MD_Stop)
					} else {
						elevio.SetMotorDirection(d)
					}
			*/
		}
	}
}
