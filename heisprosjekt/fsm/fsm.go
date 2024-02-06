package fsm

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"fmt"
)

// Finite state machine

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
}

func Fsm_onRequestButtonPress(e elevator.Elevator, btn_floor int, btn_type elevio.ButtonType) {

	// fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.toString())
	// elevator_print(elevator)
	//

	switch e.Behaviour {
	case elevator.EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(e, btn_floor, btn_type) {
			timer_start(e.DoorOpenDuration_s)
		} else {
			e.Requests[btn_floor][btn_type] = 1
		}

	case elevator.EB_Moving:
		e.Requests[btn_floor][btn_type] = 1

	case elevator.EB_Idle:
		e.Requests[btn_floor][btn_type] = 1
		pair := requests.Requests_chooseDirection(e)
		e.Dirn = pair.Dirn
		e.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			outputDevice.doorLight(1)
			timer_start(e.config.doorOpenDuration_s)
			e = requests.Requests_clearAtCurrentFloor(e)

		case elevator.EB_Moving:
			outputDevice.motorDirection(e.dirn)

		case elevator.EB_Idle:
		}
	}

	setAllLights(e)

	fmt.Println("\nNew state:")
}

// func fsm_onRequestButtonPress(btn_floor int, btn_type Button) {
// 	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.toString())
// 	elevator_print(elevator)

// 	switch elevator.behaviour {
// 	case EB_DoorOpen:
// 		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
// 			timer_start(elevator.config.doorOpenDuration_s)
// 		} else {
// 			elevator.requests[btn_floor][btn_type] = 1
// 		}

// 	case EB_Moving:
// 		elevator.requests[btn_floor][btn_type] = 1

// 	case EB_idle:
// 		elevator.requests[btn_floor][btn_type] = 1
// 		DirnBehaviourPair pair = requests_chooseDirection(elevator)
// 		elevator.dirn = pair.dirn
// 		elevator.behaviour = pair.behaviour
// 			switch pair.behaviour {
// 			case EB_DoorOpen:
// 				outputDevice.doorLight(1)
// 				timer_start(elevator.config.doorOpenDuration_s)
// 				elevator = requests_clearAtCurrentFloor(elevator)

// 			case EB_Moving:
// 				outputDevice.motorDirection(elevator.dirn)

// 			case EB_Idle:
// 		}
// 	}

// 	setAllLights(elevator)

// 	fmt.Println("\nNew state:")
// }

func Fsm_onDoorTimeout() {

	//printf("\n\n%s()\n", __FUNCTION__);
	//elevator_print(elevator);

	switch elevator.Elevator.Behaviour {
	case elevator.EB_DoorOpenEB_DoorOpen:
		//timer_start(elevator.config.doorOpenDuration_s);
		elevator = requests.Requests_clearAtCurrentFloor(elevator.Elevator)
		setAllLights(elevator)
	case elevator.EB_Moving:
	case elevator.EB_Idle:
		outputDevice.doorLight(0)
		outputDevice.motorDirection(elevator.dirn)
	default:

	}

	//printf("\nNew state:\n");
	//elevator_print(elevator);

}
