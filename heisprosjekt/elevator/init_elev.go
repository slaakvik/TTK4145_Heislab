package elevator

import (
	"Heis/driver-go/elevio"
)

func InitElev() Elevator { //elevator_unitialized?

	return Elevator{
		Floor: -1,
		Dirn:  elevio.MD_Stop,
		Requests: [elevio.NumFloors][elevio.NumButtons]bool{{false, false, false},
			{false, false, false},
			{false, false, false},
			{false, false, false}}, //vi må vel gjøre slik at man setter antall etasjer (og knapper) og itererer gjennom en for-løkke for å lage Requests
		Behaviour: EB_Idle,
		ElevID:    "",
		Failure:   false,
		//DoorOpenDuration_s: 3.0,
		//DoorOpen:           false,
	}
}

func OnInitBetweenFloors(elev Elevator) Elevator{

	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Behaviour = EB_Moving
	elev.Dirn = elevio.MD_Down
	return elev
}
