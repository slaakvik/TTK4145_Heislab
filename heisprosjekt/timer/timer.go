package timer

import (
	"Heis/driver-go/elevio"
	"time"
)

// func DoorTimer( /*elev *elevator.Elevator,*/ quit chan int) { // Door opens for three seconds
// 	time.Sleep(3 * time.Second)
// 	quit <- 1
// }

func Timer(doorTimerCh chan bool,  timedOut chan int) { // Kjører denne som goroutine

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
		case a:= <- doorTimerCh:
			resetDoorTimer = a
		case <-timer.C:
			timedOut <- 1
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}
}
