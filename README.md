# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering

import (
	"fmt"
	// "Heis/noe noe elevator_io_types ??"
)

func fsm_onRequestButtonPress(btn_floor int, btn_type Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.toString())
	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			timer_start(elevator.config.doorOpenDuration_s)
		} else {
			elevator.requests[btn_floor][btn_type] = 1
		}
	
	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = 1
	
	case EB_idle:
		elevator.requests[btn_floor][btn_type] = 1
		DirnBehaviourPair pair = requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
			switch pair.behaviour {
			case EB_DoorOpen:
				outputDevice.doorLight(1)
				timer_start(elevator.config.doorOpenDuration_s)
				elevator = requests_clearAtCurrentFloor(elevator)

			case EB_Moving:
				outputDevice.motorDirection(elevator.dirn)
			
			case EB_Idle:
		}
	}

	setAllLights(elevator)

	fmt.Println("\nNew state:")
	elevator_print(elevator)
}
