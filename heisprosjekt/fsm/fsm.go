package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"fmt"
	"time"
)

// Finite state machine

func FSM(buttons chan elevio.ButtonEvent, floors chan int, stop chan bool, obstr chan bool, elevatorID string, isMaster bool, sendElev chan elevator.Elevator, updatedCostTx chan map[string][][2]bool, masterNewOrder chan map[string][][2]bool, receiveElev chan elevator.Elevator, updatedCostRx chan map[string][][2]bool) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorID

	if elevio.GetFloor() == -1 {
		elev = onInitBetweenFloors(elev)
		elevator.Elevator_print(elev)
	}

	//dersom man er master skal man sette opp en cost func map og legge til seg selv:
	mapOfElevs := make(map[string]elevator.Elevator)
	fmt.Println(mapOfElevs)
	fmt.Println("______________")
	//mapOfElevs[elev.ElevID] = elev

	//cost function test:
	// elev.Floor = elevio.GetFloor()
	// input := cost_fns.InputToCost(elev)
	// cost_fns.HRA_funcs(input)
	// elevator.Elevator_print(elev)
	for {
		//elevator.Elevator_print(elev)
		//fmt.Print(eleviId)
		select {
		case a := <-buttons:
			btn_floor := a.Floor
			btn_type := a.Button
			fmt.Printf("Button: %+v\n", a)
			//onRequestButtonPress(&elev, a.Floor, a.Button)

			if elev.Behaviour != elevator.EB_Moving && requests.Requests_shouldClearImmediately(elev, btn_floor, btn_type) {
				elev.Behaviour = elevator.EB_DoorOpen
				mapOfElevs[elev.ElevID] = elev
				sendElevToMaster(isMaster, sendElev, elev)

				setAllLights(elev)
				quit := make(chan int)
				go doorTimer(quit)
				<-quit
				elev = onDoorTimeout(elev)

			} else {
				if isMaster {
					if btn_type == 2 {
						fmt.Println("cab callll")
						elev.Requests[btn_floor][btn_type] = true
						mapOfElevs[elev.ElevID] = elev
						fmt.Println(mapOfElevs)

						// //Her må vi kjøre cost funcksjonen
						// runcostresult := runCostFunc(mapOfElevs)
						// //Sende resultat fra cost funksjon til alle.
						// updatedCostTx <- runcostresult

						// MyNewHRs := runcostresult[elev.ElevID]
						// elev = onRequest(elev, MyNewHRs)
						// mapOfElevs[elev.ElevID] = elev
					} else {
						temp_elev := elev
						temp_elev.Requests[btn_floor][btn_type] = true
						mapOfElevs[elev.ElevID] = temp_elev
						// //Her må vi kjøre cost funcksjonen
						// runcostresult := runCostFunc(mapOfElevs)
						// //Sende resultat fra cost funksjon til alle.
						// updatedCostTx <- runcostresult
						// fmt.Println("Hall calllll")
						// // masterNewOrder <- runcostresult
						// MyNewHRs := runcostresult[elev.ElevID]
						// elev = onRequest(elev, MyNewHRs)
						// mapOfElevs[elev.ElevID] = elev
					}
					//Her må vi kjøre cost funcksjonen
					runcostresult := runCostFunc(mapOfElevs)
					//Sende resultat fra cost funksjon til alle.
					updatedCostTx <- runcostresult
					//fmt.Println(mapOfElevs)
					fmt.Println("Map of elevators:")
					fmt.Println(mapOfElevs)
					fmt.Println("Output from cost function")
					fmt.Println(runcostresult)
					fmt.Println()

					MyNewHRs := runcostresult[elev.ElevID]
					elev = onRequest(elev, MyNewHRs)
					mapOfElevs[elev.ElevID] = elev

				} else {

					if btn_type == 2 {
						elev.Requests[btn_floor][btn_type] = true
						sendElev <- elev
					} else {
						temp_elev := elev
						temp_elev.Requests[btn_floor][btn_type] = true
						sendElev <- temp_elev
					}
				}
			}

		case a := <-receiveElev:
			//fmt.Printf("ElevatorRx: %+v\n", a)
			if isMaster {
				mapOfElevs[a.ElevID] = a
				//Her må vi kjøre cost funcksjonen
				runcostresult := runCostFunc(mapOfElevs)
				fmt.Println()
				fmt.Println(runcostresult)
				fmt.Println()

				//Sende resultat fra cost funksjon til alle.
				updatedCostTx <- runcostresult
				//masterNewOrder <- runcostresult
				MyNewHRs := runcostresult[elev.ElevID]
				elev = onRequest(elev, MyNewHRs)
				mapOfElevs[elev.ElevID] = elev
			}
		// case a := <-masterNewOrder:
		// 	if isMaster {
		// 		fmt.Printf("Master got a new order: %+v\n", a)
		// 		MyNewHRs := a[elev.ElevID]
		// 		elev = onRequest(elev, MyNewHRs)
		// 		mapOfElevs[elev.ElevID] = elev

		// 	}
		case a := <-updatedCostRx:
			if !isMaster {
				//fmt.Printf("CostRx: %+v\n", a)
				MyNewHRs := a[elev.ElevID]
				elev = onRequest(elev, MyNewHRs)
				mapOfElevs[elev.ElevID] = elev
				//fmt.Println(mapOfElevs)
				sendElevToMaster(isMaster, sendElev, elev)
			}

		case a := <-floors: // a er etasjen heisen er i
			fmt.Printf("Floor: %+v\n", a)
			//onFloorArrival(&elev, a)
			//if elev.Floor != a {

			elev.Floor = a
			elevio.SetFloorIndicator(elev.Floor)
			mapOfElevs[elev.ElevID] = elev
			sendElevToMaster(isMaster, sendElev, elev)
			//	}
			switch elev.Behaviour {
			case elevator.EB_Moving:
				if requests.Requests_shouldStop(elev) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.Behaviour = elevator.EB_DoorOpen
					mapOfElevs[elev.ElevID] = elev
					sendElevToMaster(isMaster, sendElev, elev)
					setAllLights(elev)
					quit := make(chan int)
					go doorTimer(quit)
					<-quit
					elev = onDoorTimeout(elev)
				}
				mapOfElevs[elev.ElevID] = elev
				sendElevToMaster(isMaster, sendElev, elev)
			default:
			}
		case a := <-stop:
			fmt.Printf("Stop button: %+v\n", a)
			// fsm.Fsm_onStopButtonPress(&elev)
			stop_functionality(a, elev)

		case a := <-obstr:
			fmt.Printf("Obstruction %+v\n", a)
			obstruction_functionality(a, elev)
		}
	}
}

// Timer goroutine
func doorTimer( /*elev *elevator.Elevator,*/ quit chan int) { // Door opens for three seconds
	time.Sleep(3 * time.Second)
	quit <- 1
}

func sendElevToMaster(isMaster bool, updateElev chan elevator.Elevator, elev elevator.Elevator) {
	if !isMaster {
		updateElev <- elev
	}
}

func obstruction_functionality(obstruct bool, elev elevator.Elevator) {
	if obstruct {
		elevio.SetMotorDirection(elevio.MD_Stop)
		if elevio.GetFloor() != -1 {
			elevio.SetDoorOpenLamp(true)
		}
	} else {
		elevio.SetMotorDirection(elev.Dirn)
		setAllLights(elev)
	}
}
func stop_functionality(stop bool, elev elevator.Elevator) {
	if stop {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetStopLamp(true)
	} else {
		elevio.SetMotorDirection(elev.Dirn)
		elevio.SetStopLamp(false)
	}
}

func onInitBetweenFloors(elev elevator.Elevator) elevator.Elevator {

	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Behaviour = elevator.EB_Moving
	elev.Dirn = elevio.MD_Down
	return elev
}

func setAllLights(elev elevator.Elevator) {
	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < elevio.NumButtons; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elev.Requests[f][b])
		}
	}
	elevio.SetDoorOpenLamp(elev.Behaviour == elevator.EB_DoorOpen)

}

func runCostFunc(elevMap map[string]elevator.Elevator) map[string][][2]bool {
	input := cost_fns.InputToCost(elevMap)
	return cost_fns.GetCostOutput(input)
}

func onDoorTimeout(elev elevator.Elevator) elevator.Elevator {

	//printf("\n\n%s()\n", __FUNCTION__);
	//elevator_print(elevator);

	elev = requests.Requests_clearAtCurrentFloor(elev)
	pair := requests.Requests_chooseDirection(elev)
	elev.Dirn = pair.Dirn
	elev.Behaviour = pair.Behaviour
	setAllLights(elev)
	if elev.Behaviour == elevator.EB_DoorOpen {
		quit := make(chan int)
		go doorTimer(quit)
		<-quit
		elev = onDoorTimeout(elev)
	}
	elevio.SetMotorDirection(elev.Dirn) //Det er her det skjærer seg (?), den klarer ikke å ha to hall calls i samme etasje, og ingen andre requests.

	//fmt.Println("\nNew state:")
	//elevator.Elevator_print(*elev)
	return elev
}

func onRequest(elev elevator.Elevator, HallRequests [][2]bool) elevator.Elevator {
	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		elev.Requests = elevator.MergeHallAndRequests(elev.Requests, HallRequests)
		fmt.Println("HeiHei")
	case elevator.EB_Moving:
		elev.Requests = elevator.MergeHallAndRequests(elev.Requests, HallRequests)

	case elevator.EB_Idle:
		elev.Requests = elevator.MergeHallAndRequests(elev.Requests, HallRequests)
		pair := requests.Requests_chooseDirection(elev)
		elev.Dirn = pair.Dirn
		elevio.SetMotorDirection(elev.Dirn)
		elev.Behaviour = pair.Behaviour
	}

	setAllLights(elev)
	return elev
}

// func onFloorArrival(elev *elevator.Elevator, newFloor int) { //elevptr
// 	//printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
// 	//elevator_print(elevator);
// 	elev.Floor = newFloor

// 	elevio.SetFloorIndicator(elev.Floor)

// 	switch elev.Behaviour {
// 	case elevator.EB_Moving:
// 		if requests.Requests_shouldStop(*elev) {
// 			elevio.SetMotorDirection(elevio.MD_Stop)
// 			elev.Behaviour = elevator.EB_DoorOpen
// 			setAllLights(*elev)
// 			quit := make(chan int)
// 			go doorTimer(quit)
// 			<-quit
// 			*elev = onDoorTimeout(*elev)
// 			//elevio.SetDoorOpenLamp(true)
// 			// elev = requests.Requests_clearAtCurrentFloor(*elev)
// 			//requests.Requests_clearAtCurrentFloor(elev)

// 			//time.Sleep(time.Duration(elev.DoorOpenDuration_s) * time.Second)

// 			//@ Det under skrev vi onsdag kveld.
// 			// pair := requests.Requests_chooseDirection(*elev)
// 			// elev.Dirn = pair.Dirn
// 			// elev.Behaviour = pair.Behaviour
// 			// setAllLights(*elev) //@ need to rewrite this function. Hvorfor er denne funksjonen her i det hele tatt?
// 			// elevio.SetMotorDirection(elev.Dirn)
// 		}
// 	default:
// 		// elev.Dirn = elevio.MD_Stop
// 		// elevio.SetMotorDirection(elev.Dirn)
// 		//@ Disse linjene må vel være med? Slik at heisen stopper hvis den ikke allerede er moving
// 	}

// }

// func Fsm_onStopButtonPress(elev *elevator.Elevator) {
// 	fmt.Printf("%s()\n", "STOPP DA!")
// }
