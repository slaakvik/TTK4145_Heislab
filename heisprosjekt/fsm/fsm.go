package fsm

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"time"
)

// Finite state machine

// Timer goroutine
func FSM_DoorTimer( /*elev *elevator.Elevator,*/ quit chan int) { // Door opens for three seconds
	time.Sleep(3 * time.Second)
	quit <- 1
}

func Obstruction_functionality(obstruct bool, elev elevator.Elevator) {
	if obstruct {
		elevio.SetMotorDirection(elevio.MD_Stop)
		if elevio.GetFloor() != -1 {
			elevio.SetDoorOpenLamp(true)
		}
	} else {
		elevio.SetMotorDirection(elev.Dirn)
		setAllLights(elev)
	}
}
func Stop_functionality(stop bool, elev elevator.Elevator) {
	if stop {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetStopLamp(true)
	} else {
		elevio.SetMotorDirection(elev.Dirn)
		elevio.SetStopLamp(false)
	}
}

func OnInitBetweenFloors(elev elevator.Elevator) elevator.Elevator {

	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Behaviour = elevator.EB_Moving
	elev.Dirn = elevio.MD_Down
	return elev
}

func setAllLights(elev elevator.Elevator) {
	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < elevio.NumButtons; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elev.Requests[f][b])
		}
	}
	elevio.SetDoorOpenLamp(elev.Behaviour == elevator.EB_DoorOpen)

}

func OnRequestButtonPress(elev *elevator.Elevator, btn_floor int, btn_type elevio.ButtonType) {

	// fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.toString())
	// elevator_print(elevator)
	//

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(*elev, btn_floor, btn_type) {

			setAllLights(*elev)
			quit := make(chan int)
			go FSM_DoorTimer(quit)

			<-quit
			*elev = OnDoorTimeout(*elev)

			// elev.Behaviour = elevator.EB_Idle
			// elevio.SetDoorOpenLamp(0)

		} else {
			elev.Requests[btn_floor][btn_type] = true
		}

	case elevator.EB_Moving:
		elev.Requests[btn_floor][btn_type] = true

	case elevator.EB_Idle:
		if requests.Requests_shouldClearImmediately(*elev, btn_floor, btn_type) {
			elev.Behaviour = elevator.EB_DoorOpen
			setAllLights(*elev)
			quit := make(chan int)
			go FSM_DoorTimer(quit)
			<-quit
			*elev = OnDoorTimeout(*elev)

		} else {
			elev.Requests[btn_floor][btn_type] = true
			pair := requests.Requests_chooseDirection(*elev)
			elev.Dirn = pair.Dirn
			//elevio.SetMotorDirection(elev.Dirn)
			elev.Behaviour = pair.Behaviour
			setAllLights(*elev)
		}
		switch elev.Behaviour {
		case elevator.EB_DoorOpen:
			//elevio.SetDoorOpenLamp(true)
			//elev.DoorOpen=true
			//time.Sleep(time.Duration(elev.DoorOpenDuration_s * float32(time.Second)))
			//time.NewTimer(time.Duration(elev.DoorOpenDuration_s))

			*elev = requests.Requests_clearAtCurrentFloor(*elev) //kan dette ha noe å gjøre med at vi ikke klarer å motta bestillinger fra en etasje vi allerede er i?
			//elev.Behaviour = elevator.EB_Idle

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)

		case elevator.EB_Idle:
		}
	}

	setAllLights(*elev)

	//fmt.Println("\nNew state:")
}

// func OnDoorTimeout(elev *elevator.Elevator) {

// 	//printf("\n\n%s()\n", __FUNCTION__);
// 	//elevator_print(elevator);

// 	requests.Requests_clearAtCurrentFloor(elev)
// 	pair := requests.Requests_chooseDirection(*elev)
// 	elev.Dirn = pair.Dirn
// 	elev.Behaviour = pair.Behaviour
// 	setAllLights(*elev) //@ need to rewrite this function. Hvorfor er denne funksjonen her i det hele tatt?
// 	elevio.SetMotorDirection(elev.Dirn)

// 	//fmt.Println("\nNew state:")
// 	//elevator.Elevator_print(*elev)

// }

func OnDoorTimeout(elev elevator.Elevator) elevator.Elevator {

	//printf("\n\n%s()\n", __FUNCTION__);
	//elevator_print(elevator);

	elev = requests.Requests_clearAtCurrentFloor(elev)
	pair := requests.Requests_chooseDirection(elev)
	elev.Dirn = pair.Dirn
	elev.Behaviour = pair.Behaviour
	setAllLights(elev) //@ need to rewrite this function. Hvorfor er denne funksjonen her i det hele tatt?
	elevio.SetMotorDirection(elev.Dirn)

	//fmt.Println("\nNew state:")
	//elevator.Elevator_print(*elev)
	return elev
}
func OnFloorArrival(elev *elevator.Elevator, newFloor int) { //elevptr
	//printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
	//elevator_print(elevator);
	elev.Floor = newFloor

	elevio.SetFloorIndicator(elev.Floor)

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(*elev) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elev.Behaviour = elevator.EB_DoorOpen
			setAllLights(*elev)
			quit := make(chan int)
			go FSM_DoorTimer(quit)
			<-quit
			*elev = OnDoorTimeout(*elev)
			//elevio.SetDoorOpenLamp(true)
			// elev = requests.Requests_clearAtCurrentFloor(*elev)
			//requests.Requests_clearAtCurrentFloor(elev)

			//time.Sleep(time.Duration(elev.DoorOpenDuration_s) * time.Second)

			//@ Det under skrev vi onsdag kveld.
			// pair := requests.Requests_chooseDirection(*elev)
			// elev.Dirn = pair.Dirn
			// elev.Behaviour = pair.Behaviour
			// setAllLights(*elev) //@ need to rewrite this function. Hvorfor er denne funksjonen her i det hele tatt?
			// elevio.SetMotorDirection(elev.Dirn)
		}
	default:
		// elev.Dirn = elevio.MD_Stop
		// elevio.SetMotorDirection(elev.Dirn)
		//@ Disse linjene må vel være med? Slik at heisen stopper hvis den ikke allerede er moving
	}

}

// func Fsm_onStopButtonPress(elev *elevator.Elevator) {
// 	fmt.Printf("%s()\n", "STOPP DA!")
// }
