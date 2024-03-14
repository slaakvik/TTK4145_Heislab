package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"

	// "Heis/network/establish_connection"
	// "Heis/network/tcp"
	"Heis/requests"
	"fmt"
)

// Finite state machine

func ButtonsAndRequests(masterPort string, elevatorID string, isMaster bool, elevUpdateRealtimeCh <-chan elevator.Elevator,
	drv_buttons chan elevio.ButtonEvent, /*elevatorTx chan elevator.Elevator,*/
	sendMapToSlavesCh chan<- map[string]elevator.Elevator, getElevFromSlave chan elevator.Elevator,
	receiveMapFromMasterCh <-chan map[string]elevator.Elevator, newOrderCh chan<- map[string]elevator.Elevator,
	lightsCh <-chan int, sendMyselfToMaster chan elevator.Elevator, doorTimerChForBtnFSM chan bool) {
	fmt.Println("[ButtonsAndRequests] akkurat kommet inni")
	elev := elevator.InitElev()
	elev.ElevID = elevatorID

	mapOfElevs := make(map[string]elevator.Elevator)
	mapOfElevs[elev.ElevID] = elev

	fmt.Println("______________")

	for {
		select {
		case a := <-elevUpdateRealtimeCh:
			fmt.Printf("[ButtonsAndRequests] mottok en heis på elevUpdateRealtimeCh: %v\n", a)
			elev = a
			mapOfElevs[elev.ElevID] = elev
			fmt.Printf("[ButtonsAndRequests] mapOfElevs ser slik ut nå: %v\n", mapOfElevs)
		case a := <-drv_buttons:
			btn_floor := a.Floor
			btn_type := a.Button
			fmt.Printf("[ButtonsAndRequests] Button: %+v\n", a)
			if elev.Behaviour != elevator.EB_Moving && requests.ShouldClearImmediately(elev, btn_floor, btn_type) {
				fmt.Println("[ButtonsAndRequests] kom meg forbi if-state i drv_buttons")
				elev.Behaviour = elevator.EB_DoorOpen
				mapOfElevs[elev.ElevID] = elev
				go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
				elevio.SetDoorOpenLamp(true)
				elevio.SetButtonLamp(btn_type, btn_floor, true)
				doorTimerChForBtnFSM <- true

			} else {
				if isMaster {
					elev.Requests[btn_floor][btn_type] = true
					fmt.Println(" [ButtonsAndRequests] master har fått buttonpress")
					mapOfElevs[elev.ElevID] = elev
					mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
					sendMapToSlavesCh <- mapOfElevs
					fmt.Println("[ButtonsAndRequests] Sendte til alle")
					newOrderCh <- mapOfElevs

				} else {

					elev.Requests[btn_floor][btn_type] = true
					// elevatorTx <- elev
					go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
					// tcpConn, _ := establish_connection.TransmitConn(masterPort, elev.ElevID) // vi må ha masterIp
					// tcp.Transmit(tcpConn, elev)
				}
			}
		case a := <-getElevFromSlave:
			fmt.Println("[ButtonsAndRequests] mottok noe på getElevFromSlave")
			if isMaster {
				fmt.Printf("[ButtonsAndRequests] Jeg mottok en heis nå: %v\n", a)
				mapOfElevs[a.ElevID] = a
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)

				sendMapToSlavesCh <- mapOfElevs
				newOrderCh <- mapOfElevs
			}

		case a := <-receiveMapFromMasterCh:
			fmt.Println("[ButtonsAndRequests] mottok noe på receiveMapFromMasterCh")
			if !isMaster {
				fmt.Printf("[ButtonsAndRequests] Jeg mottok et map nå: %v\n", a)
				mapOfElevs = a
				fmt.Println("[ButtonsAndRequests] skal sende mappet mapOfElevs til newOrderCh")
				newOrderCh <- mapOfElevs
			}
		case <-lightsCh:
			fmt.Println("[ButtonsAndRequests] mottok noe på lightsCh")
			elevator.SetAllLights(elev, mapOfElevs)

		}
	}
}

func FloorObstrStop(masterPort string, isMaster bool, elevatorId string, elevUpdateRealtimeCh chan<- elevator.Elevator, drv_floors chan int /*elevatorTx chan elevator.Elevator,*/, newOrderCh <-chan map[string]elevator.Elevator, doorTimerCh chan bool, timedOut chan int, lightsCh chan<- int, sendMyselfToMaster chan elevator.Elevator) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorId
	fmt.Println("[FloorObstrStop] kommet meg inni")
	if elevio.GetFloor() == -1 {
		elev = elevator.OnInitBetweenFloors(elev)
		elevUpdateRealtimeCh <- elev
		elevator.Elevator_print(elev)
	}
	fmt.Println("[FloorObstrStop] kommet til for-loopen")
	for {
		select {
		case a := <-newOrderCh:
			elev = a[elev.ElevID]
			fmt.Printf("[FloorObstrStop] mottok på newOrderCh: %v\n",a)
			fmt.Println("[FloorObstrStop] sender den til updateRealTimeCh")
			elevUpdateRealtimeCh <- elev
			fmt.Println("[FloorObstrStop] kommet meg inni")
			go SendElevToMaster(isMaster, elev, sendMyselfToMaster) //(trenger vi denne?)
			elev = requests.OnRequest(elev, lightsCh)
			elevUpdateRealtimeCh <- elev
			go SendElevToMaster(isMaster, elev, sendMyselfToMaster) // på tråd?

		case a := <-drv_floors:
			fmt.Printf("[FloorObstrStop] Floor: %+v\n", a)
			elev.Floor = a
			elevUpdateRealtimeCh <- elev
			elevio.SetFloorIndicator(elev.Floor)
			go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
			switch elev.Behaviour {
			case elevator.EB_Moving:
				if requests.ShouldStop(elev) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.Behaviour = elevator.EB_DoorOpen
					elevUpdateRealtimeCh <- elev
					go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
					lightsCh <- 1
					doorTimerCh <- true
				}

			default:

			}

		case <-timedOut:
			fmt.Println("[FloorObstrStop] gikk inn i timedout")
			elev = requests.OnDoorTimeout(elev, doorTimerCh, lightsCh, elevUpdateRealtimeCh)
			fmt.Println("gikk inn i timedout")
			elevUpdateRealtimeCh <- elev

			SendElevToMaster(isMaster, elev, sendMyselfToMaster)
			fmt.Println("[FloorObstrStop] Managed to send floor update")
			// case a := <-drv_stop:
			// 	fmt.Printf("Stop button: %+v\n", a)
			// 	stop_functionality(a, elev)

		}
	}
}

func SendElevToMaster(isMaster bool, elev elevator.Elevator, sendMyselfToMaster chan elevator.Elevator) {
	if !isMaster {
		fmt.Println("[SendElevToMaster] kommet meg inni, og jeg er ikke master ")
		fmt.Println("[FloorObstrStop] sender heisen til sendMyselfToMasterCh")
		sendMyselfToMaster <- elev
		// tcpConn, _ := establish_connection.TransmitConn(masterPort, elev.ElevID) //vi må ha masterIp
		// tcp.Transmit(tcpConn, elev)
	}
	fmt.Println("[SendElevToMaster] går ut")
}

func stop_functionality(stop bool, elev elevator.Elevator) {
	fmt.Println("[stop_functionality] kommet meg inni ")
	if stop {
		fmt.Println("[stop_functionality] forbi if stop ")
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetStopLamp(true)
	} else {
		fmt.Println("[stop_functionality] forbi else ")
		elevio.SetMotorDirection(elev.Dirn)
		elevio.SetStopLamp(false)
	}
	fmt.Println("[stop_functionality] går ut ")
}