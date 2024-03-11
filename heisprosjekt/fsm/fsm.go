package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"Heis/timer"
	"fmt"
	"time"
)

// Finite state machine

func ButtonsAndRequests(elevatorID string, isMaster bool /*elevUpdateBtnAndOrdersCh chan<- elevator.Elevator, */, elevUpdateRealtimeCh <-chan elevator.Elevator, drv_buttons chan elevio.ButtonEvent, ElevatorTx chan elevator.Elevator, mapOfElevsTx chan map[string]elevator.Elevator, ElevatorRx chan elevator.Elevator, mapOfElevsRx chan map[string]elevator.Elevator, newOrderCh chan<- map[string]elevator.Elevator) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorID

	mapOfElevs := make(map[string]elevator.Elevator)
	mapOfElevs[elev.ElevID] = elev

	//newOrderCh <- mapOfElevs

	//dersom man er master skal man sette opp en cost func map og legge til seg selv:
	fmt.Println(mapOfElevs)
	fmt.Println("______________")

	for {
		select {
		case a := <-elevUpdateRealtimeCh:
			elev = a
			mapOfElevs[elev.ElevID] = elev
		case a := <-drv_buttons:
			btn_floor := a.Floor
			btn_type := a.Button
			fmt.Printf("Button: %+v\n", a)

			// if elev.Behaviour != elevator.EB_Moving && requests.Requests_shouldClearImmediately(elev, btn_floor, btn_type) {

			// 	elev.Behaviour = elevator.EB_DoorOpen
			// 	elevUpdateBtnAndOrdersCh <- elev

			// 	mapOfElevs[elev.ElevID] = elev
			// 	sendElevToMaster(isMaster, ElevatorTx, elev)
			// 	elevator.SetAllLights(elev)
			// 	go func() {

			// 		quit := make(chan int)
			// 		go timer.DoorTimer(quit)
			// 		<-quit
			// 	}()
			// 	elev = timer.OnDoorTimeout(elev)
			// 	elevUpdateBtnAndOrdersCh <- elev
			// 	elevio.SetMotorDirection(elev.Dirn)

			// } else {
			if isMaster {
				elev.Requests[btn_floor][btn_type] = true

				// elevUpdateBtnAndOrdersCh <- elev //Trenger vi denne?

				mapOfElevs[elev.ElevID] = elev

				// costFunctionResults := cost_fns.RunCostFunc(mapOfElevs)

				// CostTx <- costFunctionResults
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
				mapOfElevsTx <- mapOfElevs
				newOrderCh <- mapOfElevs

				// MyNewHRs := costFunctionResults[elev.ElevID]
				// elev = requests.OnRequest(elev, MyNewHRs)
				// elevUpdateBtnAndOrdersCh <- elev
				// mapOfElevs[elev.ElevID] = elev

			} else {

				elev.Requests[btn_floor][btn_type] = true
				ElevatorTx <- elev

			}
			// }
		case a := <-ElevatorRx:
			if isMaster {
				mapOfElevs[a.ElevID] = a
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
				mapOfElevsTx <- mapOfElevs
				newOrderCh <- mapOfElevs
				// MyNewHRs := costFunctionResults[elev.ElevID]
				// elev = requests.OnRequest(elev, MyNewHRs)
				// elevUpdateBtnAndOrdersCh <- elev
				// mapOfElevs[elev.ElevID] = elev
			}
		case a := <-mapOfElevsRx:
			if !isMaster {
				mapOfElevs = a
				newOrderCh <- mapOfElevs
				// MyNewHRs := a[elev.ElevID]
				// elev = requests.OnRequest(elev, MyNewHRs)
				// elevUpdateBtnAndOrdersCh <- elev
				// mapOfElevs[elev.ElevID] = elev
				// sendElevToMaster(isMaster, ElevatorTx, elev)
			}

		}
	}
}

func FloorObstrStop(isMaster bool, elevatorId string /* elevUpdateBtnAndOrdersCh <-chan elevator.Elevator,*/, elevUpdateRealtimeCh chan<- elevator.Elevator, drv_floors chan int, drv_stop chan bool, drv_obstr chan bool, ElevatorTx chan<- elevator.Elevator, newOrderCh <-chan map[string]elevator.Elevator) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorId

	if elevio.GetFloor() == -1 {
		elev = elevator.OnInitBetweenFloors(elev)
		elevator.Elevator_print(elev)
	}

	for {
		select {
		// case a := <-elevUpdateBtnAndOrdersCh:
		// 	elev = a
		// 	fmt.Println("Elev floor at btn press: ", elev.Floor)
		case a := <-newOrderCh:
			elev = a[elev.ElevID]
			// elevUpdateRealtimeCh <- elev
			// sendElevToMaster(isMaster, ElevatorTx, elev) //(trenger vi denne?)
			elev = requests.OnRequest(elev)
			elevUpdateRealtimeCh <- elev

		case a := <-drv_floors:
			fmt.Printf("Floor: %+v\n", a)
			elev.Floor = a
			elevUpdateRealtimeCh <- elev
			elevio.SetFloorIndicator(elev.Floor)
			sendElevToMaster(isMaster, ElevatorTx, elev)
			switch elev.Behaviour {
			case elevator.EB_Moving:
				go func() {
					if requests.Requests_shouldStop(elev) {
						elevio.SetMotorDirection(elevio.MD_Stop)
						elev.Behaviour = elevator.EB_DoorOpen
						elevUpdateRealtimeCh <- elev
						sendElevToMaster(isMaster, ElevatorTx, elev)
						elevator.SetAllLights(elev)
						// quit := make(chan int)
						// go timer.DoorTimer(quit)
						// <-quit
						time.Sleep(3 * time.Second)
						elev = timer.OnDoorTimeout(elev)
						elevUpdateRealtimeCh <- elev
						elevio.SetMotorDirection(elev.Dirn)
						sendElevToMaster(isMaster, ElevatorTx, elev)
						fmt.Println("Managed to send floor update")
					}
				}()
			default:

			}
		case a := <-drv_stop:
			fmt.Printf("Stop button: %+v\n", a)
			stop_functionality(a, elev)

		case a := <-drv_obstr:
			fmt.Printf("Obstruction %+v\n", a)
			elev = obstruction_functionality(a, elev)
			elevUpdateRealtimeCh <- elev
			sendElevToMaster(isMaster, ElevatorTx, elev)

		}
	}

}

/*func FSM(buttons chan elevio.ButtonEvent, floors chan int, stop chan bool, obstr chan bool, elevatorID string, isMaster bool, sendElev chan elevator.Elevator, updatedCostTx chan map[string][][2]bool, masterNewOrder chan map[string][][2]bool, receiveElev chan elevator.Elevator, updatedCostRx chan map[string][][2]bool) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorID

	if elevio.GetFloor() == -1 {
		elev = elevator.OnInitBetweenFloors(elev)
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
		// case a := <-elevChange:
		// 	elev = a
		case a := <-buttons:
			btn_floor := a.Floor
			btn_type := a.Button
			fmt.Printf("Button: %+v\n", a)
			//onRequestButtonPress(&elev, a.Floor, a.Button)

			if elev.Behaviour != elevator.EB_Moving && requests.Requests_shouldClearImmediately(elev, btn_floor, btn_type) {
				elev.Behaviour = elevator.EB_DoorOpen
				//elevChange <- elev
				mapOfElevs[elev.ElevID] = elev
				sendElevToMaster(isMaster, sendElev, elev)

				elevator.SetAllLights(elev)
				quit := make(chan int)
				go timer.DoorTimer(quit)
				<-quit
				elev = timer.OnDoorTimeout(elev)
				elevio.SetMotorDirection(elev.Dirn)

			} else {
				if isMaster {
					if btn_type == 2 {
						fmt.Println("cab callll")
						elev.Requests[btn_floor][btn_type] = true
						mapOfElevs[elev.ElevID] = elev
						fmt.Println(mapOfElevs)
					} else {
						temp_elev := elev
						temp_elev.Requests[btn_floor][btn_type] = true
						mapOfElevs[elev.ElevID] = temp_elev

					}
					//Her må vi kjøre cost funcksjonen
					runcostresult := cost_fns.RunCostFunc(mapOfElevs)
					//Sende resultat fra cost funksjon til alle.
					updatedCostTx <- runcostresult
					//fmt.Println(mapOfElevs)
					fmt.Println("Map of elevators:")
					fmt.Println(mapOfElevs)
					fmt.Println("Output from cost function")
					fmt.Println(runcostresult)
					fmt.Println()

					MyNewHRs := runcostresult[elev.ElevID]
					elev = requests.OnRequest(elev, MyNewHRs)
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
				runcostresult := cost_fns.RunCostFunc(mapOfElevs)
				fmt.Println()
				fmt.Println(runcostresult)
				fmt.Println()

				//Sende resultat fra cost funksjon til alle.
				updatedCostTx <- runcostresult
				//masterNewOrder <- runcostresult
				MyNewHRs := runcostresult[elev.ElevID]
				elev = requests.OnRequest(elev, MyNewHRs)
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
				elev = requests.OnRequest(elev, MyNewHRs)
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
					elevator.SetAllLights(elev)
					quit := make(chan int)
					go timer.DoorTimer(quit)
					<-quit
					elev = timer.OnDoorTimeout(elev)
					elevio.SetMotorDirection(elev.Dirn)
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
*/
// Timer goroutine
// Kanskje legge denn inn i en egen go routine slik at vi fortsatt kan ta ordre?.

func sendElevToMaster(isMaster bool, updateElev chan<- elevator.Elevator, elev elevator.Elevator) {
	if !isMaster {
		updateElev <- elev
	}
}

// etterhvert: fiks sånn at master blir klar over at heisen er obstructed (Må ha en "failure state" i elevator structet)
func obstruction_functionality(obstruct bool, elev elevator.Elevator) elevator.Elevator {
	if obstruct {
		elevio.SetMotorDirection(elevio.MD_Stop)
		if elevio.GetFloor() != -1 {
			elevio.SetDoorOpenLamp(true)
		}
		elev.Failure = true

	} else {
		elevio.SetMotorDirection(elev.Dirn)
		elevator.SetAllLights(elev)
		elev.Failure = false
	}
	return elev
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
