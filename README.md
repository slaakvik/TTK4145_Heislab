# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering


--------------------JONATHAN SITT--------------------------


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



---------------------------------ANDERS SITT-------------------------------

type DirnBehaviourPair struct {
	Dirn              dirn
	ElevatorBehaviour behaviour
}



func fsm_onFloorArrival(newFloor int) int{
	//printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
    //elevator_print(elevator);

	elevator.floor = newFloor

	outputDevice.floorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(	elevator) {
			output.motorDirection(D_Stop)
			outputDevice.doorLight(1);
            elevator = requests_clearAtCurrentFloor(elevator);
            //timer_start(elevator.config.doorOpenDuration_s);
            setAllLights(elevator);
            elevator.behaviour = EB_DoorOpen;
		}
	default:

	}
	//printf("\nNew state:\n");
    //elevator_print(elevator); 

}

//Trenger vi å bruke pekere her? (*elevator)?



func fsm_onDoorTimeout() { //????
	//printf("\n\n%s()\n", __FUNCTION__);
    //elevator_print(elevator);

	switch(elevator.behaviour){
	case EB_DoorOpen:
		//timer_start(elevator.config.doorOpenDuration_s);
		elevator = requests_clearAtCurrentFloor(elevator)
		setAllLights(elevator)
	case EB_Moving:
	case EB_Idle:
		outputDevice.doorLight(0)
		outputDevice.motorDirection(elevator.dirn)
	default:

	}

	//printf("\nNew state:\n");
    //elevator_print(elevator);


}




 
	elevator_print(elevator)
}
