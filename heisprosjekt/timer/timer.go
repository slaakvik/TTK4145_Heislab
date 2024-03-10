package timer

import (
	"Heis/elevator"
	"Heis/requests"
	"time"
)

func DoorTimer( /*elev *elevator.Elevator,*/ quit chan int) { // Door opens for three seconds
	time.Sleep(3 * time.Second)
	quit <- 1
}

func OnDoorTimeout(elev elevator.Elevator) elevator.Elevator {

	//printf("\n\n%s()\n", __FUNCTION__);
	//elevator_print(elevator);

	elev = requests.Requests_clearAtCurrentFloor(elev)
	pair := requests.Requests_chooseDirection(elev)
	elev.Dirn = pair.Dirn
	elev.Behaviour = pair.Behaviour
	elevator.SetAllLights(elev)
	if elev.Behaviour == elevator.EB_DoorOpen {
		// quit := make(chan int)
		// go DoorTimer(quit)
		// <-quit
		time.Sleep(3 * time.Second)
		elev = OnDoorTimeout(elev)
	}

	//fmt.Println("\nNew state:")
	//elevator.Elevator_print(*elev)
	return elev
}
