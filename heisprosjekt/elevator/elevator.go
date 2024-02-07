package elevator

import (
	"Heis/driver-go/elevio"
	"fmt"
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen                   = 1
	EB_Moving                     = 2
)

type Elevator struct {
	Floor              int
	Dirn               elevio.MotorDirection
	Requests           [elevio.NumFloors][elevio.NumButtons]bool
	Behaviour          ElevatorBehaviour
	DoorOpenDuration_s float32
	DoorOpen           bool
}

func Elevator_print(
	elevator Elevator) {
	fmt.Printf("Floor: %d, Dirn: %d, Requests: %v, Behaviour: %d, DoorOpenDuration: %f, DoorOpen: %t\n",
		elevator.Floor, elevator.Dirn, elevator.Requests, elevator.Behaviour, elevator.DoorOpenDuration_s, elevator.DoorOpen)
}

