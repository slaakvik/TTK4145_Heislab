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
	floor int
	dirn elevio.MotorDirection
	requests[][] bool
	behaviour ElevatorBehaviour
	doorOpenDuration_s	float32
	doorOpen bool

}

