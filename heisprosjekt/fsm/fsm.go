package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/network/tcp"
	"Heis/requests"
	"fmt"
)

// Finite state machine

func ButtonsAndRequests(masterPort string, elevatorID string, isMaster bool, elevUpdateRealtimeCh <-chan elevator.Elevator, drv_buttons chan elevio.ButtonEvent /*elevatorTx chan elevator.Elevator,*/, mapOfElevsTx chan map[string]elevator.Elevator, elevatorRx chan elevator.Elevator, mapOfElevsRx chan map[string]elevator.Elevator, newOrderCh chan<- map[string]elevator.Elevator, lightsCh <-chan int) {
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
				// elevatorTx <- elev
				tcpConn, _ := tcp.TransmitConn(masterPort, elev.ElevID)
				tcp.Transmiter(tcpConn, elev)

			}

		case a := <-elevatorRx:
			if isMaster {
				fmt.Println("Jeg mottok en heis nå: ", a)
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
		case <-lightsCh:
			elevator.SetAllLights(elev, mapOfElevs)

		}
	}
}

func FloorObstrStop(masterPort string, isMaster bool, elevatorId string, elevUpdateRealtimeCh chan<- elevator.Elevator, drv_floors chan int /*elevatorTx chan elevator.Elevator,*/, newOrderCh <-chan map[string]elevator.Elevator, doorTimerCh chan bool, timedOut chan int, lightsCh chan<- int) {
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
			sendElevToMaster(masterPort, isMaster, elev) //(trenger vi denne?)
			elev = requests.OnRequest(elev, lightsCh)
			elevUpdateRealtimeCh <- elev
			sendElevToMaster(masterPort, isMaster, elev)

		case a := <-drv_floors:
			fmt.Printf("Floor: %+v\n", a)
			elev.Floor = a
			elevUpdateRealtimeCh <- elev
			elevio.SetFloorIndicator(elev.Floor)
			sendElevToMaster(masterPort, isMaster, elev)
			switch elev.Behaviour {
			case elevator.EB_Moving:
				if requests.ShouldStop(elev) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.Behaviour = elevator.EB_DoorOpen
					elevUpdateRealtimeCh <- elev
					sendElevToMaster(masterPort, isMaster, elev)
					lightsCh <- 1
					doorTimerCh <- true
				}

			default:

			}

		case <-timedOut:
			fmt.Println("gikk inn i timedout")
			elev = requests.OnDoorTimeout(elev, doorTimerCh, lightsCh, elevUpdateRealtimeCh)
			fmt.Println("gikk inn i timedout")
			elevUpdateRealtimeCh <- elev

			sendElevToMaster(masterPort, isMaster, elev)
			fmt.Println("Managed to send floor update")
			// case a := <-drv_stop:
			// 	fmt.Printf("Stop button: %+v\n", a)
			// 	stop_functionality(a, elev)

		}
	}

}

func sendElevToMaster(masterPort string, isMaster bool, elev elevator.Elevator) {
	if !isMaster {
		// updateElev <- elev
		tcpConn, _ := tcp.TransmitConn(masterPort, elev.ElevID)
		tcp.Transmiter(tcpConn, elev)
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
