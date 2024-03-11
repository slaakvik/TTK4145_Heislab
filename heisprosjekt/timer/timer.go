package timer

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"fmt"
	"time"
)

func DoorTimer( /*elev *elevator.Elevator,*/ quit chan int) { // Door opens for three seconds
	time.Sleep(3 * time.Second)
	quit <- 1
}

func OnDoorTimeout(elev elevator.Elevator, doorTimerCh chan bool) elevator.Elevator {

	//printf("\n\n%s()\n", __FUNCTION__);
	//elevator_print(elevator);
	fmt.Println("Her")
	elev = requests.Requests_clearAtCurrentFloor(elev)
	fmt.Println("Her")
	pair := requests.Requests_chooseDirection(elev)
	fmt.Println("Her")
	elev.Dirn = pair.Dirn
	elev.Behaviour = pair.Behaviour
	elevator.SetAllLights(elev)
	fmt.Println("Her")

	if elev.Behaviour == elevator.EB_DoorOpen {
		// // quit := make(chan int)
		// // go DoorTimer(quit)
		// // <-quit
		// time.Sleep(3 * time.Second)
		// elev = OnDoorTimeout(elev)
		fmt.Println("Her")

		doorTimerCh <- true

	} else {
		elevio.SetMotorDirection(elev.Dirn)
	}

	//fmt.Println("\nNew state:")
	//elevator.Elevator_print(*elev)
	return elev
}

func Timer(doorTimerChannel chan bool, timedOut chan int) { // Kjører denne som goroutine

	resetDoorTimer := false
	timer := time.NewTimer(3 * time.Second)

	for {
		if elevio.GetObstruction() {
			// if elevio.GetObstruction() {
			timer.Reset(3 * time.Second)
		}

		if resetDoorTimer {
			timer.Reset(3 * time.Second)
			resetDoorTimer = false
		}

		select {
		case a := <-doorTimerChannel:
			resetDoorTimer = a
		case <-timer.C:
			timedOut <- 1
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}
}
