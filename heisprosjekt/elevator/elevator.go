package elevator

import (
	"Heis/driver-go/elevio"
	"fmt"
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2
)

type Elevator struct {
	Floor     int
	Dirn      elevio.MotorDirection
	Requests  [elevio.NumFloors][elevio.NumButtons]bool
	Behaviour ElevatorBehaviour
	//DoorOpenDuration_s float32
	//DoorOpen           bool
}

// func Elevator_print(
// 	elevator Elevator) {
// 	fmt.Printf("Floor: %d, Dirn: %d, Requests: %v, Behaviour: %d, DoorOpenDuration: %f, DoorOpen: %t\n",
// 		elevator.Floor, elevator.Dirn, elevator.Requests, elevator.Behaviour, elevator.DoorOpenDuration_s, elevator.DoorOpen)
// }

func Elevator_print(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n"+
			"  |stop = %-12.12s|\n",
		es.Floor,
		fmt.Sprintf("%v", es.Dirn),
		fmt.Sprintf("%v", es.Behaviour),
		fmt.Sprintf("%v", elevio.GetStop()),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevio.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevio.NumButtons; btn++ {
			if (f == elevio.NumFloors-1 && btn == int(elevio.BT_HallUp)) ||
				(f == 0 && btn == int(elevio.BT_HallDown)) {
				fmt.Print("|     ")
			} else {
				if es.Requests[f][btn] {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func GetCabCalls(elev Elevator) [elevio.NumFloors]bool {
	cabRequests := [elevio.NumFloors]bool{false, false, false, false}
	for f := 0; f < elevio.NumFloors; f++ {
		cabRequests[f] = elev.Requests[f][elevio.NumButtons-1]
	}
	return cabRequests
}

func GetHallCalls(elev Elevator) [elevio.NumFloors][2]bool {
	HallCalls := [elevio.NumFloors][2]bool{{false, false},
		{false, false},
		{false, false},
		{false, false}}

	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < (elevio.NumButtons - 1); b++ {
			HallCalls[f][b] = elev.Requests[f][b]
		}
	}
	return HallCalls
}

func MergeHallAndCabCall(cabs []bool, halls [][2]bool) [elevio.NumFloors][elevio.NumButtons]bool {
	requests := [elevio.NumFloors][elevio.NumButtons]bool{{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false}}

	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < (elevio.NumButtons - 1); b++ {
			requests[f][b] = halls[f][b]
		}
		requests[f][elevio.NumButtons-1] = cabs[f]
	}
	return requests
}
