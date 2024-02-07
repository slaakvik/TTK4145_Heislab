package fsm

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"fmt"
	"time"
)

// Finite state machine

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
}

func setAllLights(elev elevator.Elevator) {
	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < elevio.NumButtons; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elev.Requests[f][b])
		}
	}
	elevio.SetDoorOpenLamp(elev.Behaviour == elevator.EB_DoorOpen)
}

func Fsm_onRequestButtonPress(elev elevator.Elevator, btn_floor int, btn_type elevio.ButtonType) {

	// fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.toString())
	// elevator_print(elevator)
	//

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(elev, btn_floor, btn_type) {
			time.Sleep(time.Duration(elev.DoorOpenDuration_s * float32(time.Second)))
		} else {
			elev.Requests[btn_floor][btn_type] = true
		}

	case elevator.EB_Moving:
		elev.Requests[btn_floor][btn_type] = true

	case elevator.EB_Idle:
		elev.Requests[btn_floor][btn_type] = true
		pair := requests.Requests_chooseDirection(elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			time.Sleep(time.Duration(elev.DoorOpenDuration_s * float32(time.Second)))
			elev = requests.Requests_clearAtCurrentFloor(elev)

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)

		case elevator.EB_Idle:
		}
	}

	setAllLights(elev)

	fmt.Println("\nNew state:")
}

func Fsm_onFloorArrival(elev elevator.Elevator, newFloor int) {
	//printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
	//elevator_print(elevator);
	elev.Floor = newFloor

	elevio.GetFloor()

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(elev) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elev = requests.Requests_clearAtCurrentFloor(elev)
			time.Sleep(time.Duration(elev.DoorOpenDuration_s) * time.Second)
			setAllLights(elev) // need to rewrite this function
			elev.Behaviour = elevator.EB_DoorOpen
		}
	default:
	}

}

func Fsm_onDoorTimeout(elev elevator.Elevator) {

	//printf("\n\n%s()\n", __FUNCTION__);
	//elevator_print(elevator);

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		//timer_start(elevator.config.doorOpenDuration_s);
		elev = requests.Requests_clearAtCurrentFloor(elev)
		setAllLights(elev)
	case elevator.EB_Moving:
	case elevator.EB_Idle:
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elev.Dirn)
	default:

	}

	//printf("\nNew state:\n");
	//elevator_print(elevator);

}
