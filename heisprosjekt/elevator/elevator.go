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
	requests [elevio._numFloors][3] int
	behaviour ElevatorBehaviour
	



}


