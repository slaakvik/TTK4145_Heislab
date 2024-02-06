package elevator

import (
	"Heis/driver-go/elevio"
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	Floor              int
	Dirn               elevio.MotorDirection
	Requests           [][]int
	Behaviour          ElevatorBehaviour
	DoorOpenDuration_s float32
	DoorOpen           bool
}
