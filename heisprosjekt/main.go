package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/fsm"
	"fmt"
)

func main() {
	//heisann, din gamle ørn4!

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)

	//Initialiserer en heisstruct
	elev := elevator.InitElev()

	//Channels
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	//heartbeat := make(chan )



	//Threads
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

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
			// Fikser disse funksjonene senere
			/*
				case a := <-drv_obstr:
					fmt.Printf("%+v\n", a)
					if a {
						elevio.SetMotorDirection(elevio.MD_Stop)
					} else {
						elevio.SetMotorDirection(d)
					}

				case a := <-drv_stop:
			*/
		}
	}
}

/*
	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
*/
