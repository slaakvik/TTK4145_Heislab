package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"Heis/timer"
	"fmt"
)

// Finite state machine

func ButtonsAndRequests(elevatorID string, isMaster bool, elevUpdateRealtimeCh <-chan elevator.Elevator, drv_buttons chan elevio.ButtonEvent, ElevatorTx chan elevator.Elevator, mapOfElevsTx chan map[string]elevator.Elevator, ElevatorRx chan elevator.Elevator, mapOfElevsRx chan map[string]elevator.Elevator, newOrderCh chan<- map[string]elevator.Elevator) {
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
			if isMaster {
				elev.Requests[btn_floor][btn_type] = true

				mapOfElevs[elev.ElevID] = elev
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
				mapOfElevsTx <- mapOfElevs
				newOrderCh <- mapOfElevs

			} else {

				elev.Requests[btn_floor][btn_type] = true
				ElevatorTx <- elev

			}

		case a := <-ElevatorRx:
			if isMaster {
				mapOfElevs[a.ElevID] = a
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
				mapOfElevsTx <- mapOfElevs
				newOrderCh <- mapOfElevs

			}
		case a := <-mapOfElevsRx:
			if !isMaster {
				mapOfElevs = a
				newOrderCh <- mapOfElevs

			}

		}
	}
}

func FloorObstrStop(isMaster bool, elevatorId string, elevUpdateRealtimeCh chan<- elevator.Elevator, drv_floors chan int, drv_stop chan bool, drv_obstr chan bool, ElevatorTx chan<- elevator.Elevator, newOrderCh <-chan map[string]elevator.Elevator, doorTimerCh chan bool, timedOut chan int) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorId

	if elevio.GetFloor() == -1 {
		elev = elevator.OnInitBetweenFloors(elev)
		elevUpdateRealtimeCh <- elev
		elevator.Elevator_print(elev)
	}

	for {
		select {
		case a := <-newOrderCh:
			elev = a[elev.ElevID]

			elevUpdateRealtimeCh <- elev
			sendElevToMaster(isMaster, ElevatorTx, elev) //(trenger vi denne?)
			elev = requests.OnRequest(elev)
			elevUpdateRealtimeCh <- elev
			sendElevToMaster(isMaster, ElevatorTx, elev)

		case a := <-drv_floors:
			fmt.Printf("Floor: %+v\n", a)
			elev.Floor = a
			elevUpdateRealtimeCh <- elev
			elevio.SetFloorIndicator(elev.Floor)
			sendElevToMaster(isMaster, ElevatorTx, elev)
			switch elev.Behaviour {
			case elevator.EB_Moving:
				if requests.Requests_shouldStop(elev) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.Behaviour = elevator.EB_DoorOpen
					elevUpdateRealtimeCh <- elev
					sendElevToMaster(isMaster, ElevatorTx, elev)
					elevator.SetAllLights(elev)
					doorTimerCh <- true
				}

			default:

			}

		case <-timedOut:
			fmt.Println("gikk inn i timedout")
			elev = timer.OnDoorTimeout(elev, doorTimerCh)
			fmt.Println("gikk inn i timedout")
			elevUpdateRealtimeCh <- elev

			sendElevToMaster(isMaster, ElevatorTx, elev)
			fmt.Println("Managed to send floor update")
			// case a := <-drv_stop:
			// 	fmt.Printf("Stop button: %+v\n", a)
			// 	stop_functionality(a, elev)

		}
	}

}

func sendElevToMaster(isMaster bool, updateElev chan<- elevator.Elevator, elev elevator.Elevator) {
	if !isMaster {
		updateElev <- elev
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
