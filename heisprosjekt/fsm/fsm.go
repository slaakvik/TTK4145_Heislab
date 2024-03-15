package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"reflect"

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
	lightsCh <-chan int, sendMyselfToMaster chan elevator.Elevator) {
	
	elev := elevator.InitElev()
	elev.ElevID = elevatorID
	mapOfElevs := make(map[string]elevator.Elevator)
	mapOfElevs[elev.ElevID] = elev
	

	fmt.Println("______________")

	for {
		select {
		case a := <-elevUpdateRealtimeCh:
				elev = a
				if isMaster {
					mapOfElevs[elev.ElevID] = elev
				}else{
					sendMyselfToMaster <- elev
				}
		case a := <-drv_buttons:
			btn_floor := a.Floor
			btn_type := a.Button
			fmt.Printf("Button: %+v\n", a)
			elev.Requests[btn_floor][btn_type] = true

			if isMaster {
				fmt.Println("master har fått buttonpress")
				mapOfElevs[elev.ElevID] = elev
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
				sendMapToSlavesCh <- mapOfElevs		//Slavene begynnner å utføre ordre når den mottar denne. 
				fmt.Println("Sendte til alle")
				newOrderCh <- mapOfElevs
				// tempMap := copyMap(mapOfElevs)
					
				// mapOfElevs[elev.ElevID] = elev
				// mapOfElevs = cost_fns.RunCostFunc(mapOfElevs)
				// if !reflect.DeepEqual(tempMap, mapOfElevs) {
				// 	fmt.Println("received a new map")
				// 	sendMapToSlavesCh <- mapOfElevs
				// 	newOrderCh <- mapOfElevs
				// } 	

			} else {

				sendMyselfToMaster <- elev
			}

		case a := <-getElevFromSlave:
			if isMaster {

					tempMap := copyMap(mapOfElevs)
					
					mapOfElevs[a.ElevID] = a
					mapOfElevs = cost_fns.RunCostFunc(mapOfElevs)
				
					if !reflect.DeepEqual(tempMap, mapOfElevs) {
						sendMapToSlavesCh <- mapOfElevs
						newOrderCh <- mapOfElevs
					} 	
			}

		case a := <-receiveMapFromMasterCh:
			if !isMaster {
				mapOfElevs = copyMap(a)
				elev = mapOfElevs[elev.ElevID]
				newOrderCh <- mapOfElevs
			}
		case <-lightsCh:
			elevator.SetAllLights(elev, mapOfElevs)
		}
	}
}

func FloorObstrStop(masterPort string, isMaster bool, elevatorId string, elevUpdateRealtimeCh chan<- elevator.Elevator, 
	drv_floors chan int /*elevatorTx chan elevator.Elevator,*/, newOrderCh <-chan map[string]elevator.Elevator, doorTimerCh chan bool, 
	timedOut chan int, lightsCh chan<- int, sendMyselfToMaster chan elevator.Elevator) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorId
	// tempMap:= make(map[string]elevator.Elevator)	

	if elevio.GetFloor() == -1 {
		elev = elevator.OnInitBetweenFloors(elev)
		elevUpdateRealtimeCh <- elev
		elevator.Elevator_print(elev)
	}

	for {
		select {
		case a := <-newOrderCh:
			elev = a[elev.ElevID]
			if !isMaster {
				fmt.Println("slave received new orders.")
				fmt.Println("slaves elevator: ", elev)
				// elev.Requests[0][2]= true
			}
			if (elev.Behaviour != elevator.EB_Moving) && requests.ShouldClearImmediately(elev) {
				fmt.Println("inni should clear")
				
				elev.Behaviour = elevator.EB_DoorOpen
				// mapOfElevs[elev.ElevID] = elev
				// lightsCh <- 1
				doorTimerCh <- true
				elevUpdateRealtimeCh <- elev
				
				} else if elev.Behaviour == elevator.EB_Idle {
					fmt.Println("on request while idle")
					pair := requests.ChooseDirection(elev)
					elev.Dirn = pair.Dirn
					elevio.SetMotorDirection(elev.Dirn)
					elev.Behaviour = pair.Behaviour
					elevUpdateRealtimeCh <- elev
					
				}
				lightsCh <- 1

		case a := <-drv_floors:
			fmt.Printf("Floor: %+v\n", a)
			elev.Floor = a
			elevUpdateRealtimeCh <- elev
			elevio.SetFloorIndicator(elev.Floor)
			switch elev.Behaviour {
			case elevator.EB_Moving:
				if requests.ShouldStop(elev) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.Behaviour = elevator.EB_DoorOpen
					elevUpdateRealtimeCh <- elev
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

			fmt.Println("Managed to send floor update")
			// case a := <-drv_stop:
			// 	fmt.Printf("Stop button: %+v\n", a)
			// 	stop_functionality(a, elev)

		}
	}
}

func copyMap(source map[string]elevator.Elevator) map[string]elevator.Elevator {
	newMap := make(map[string]elevator.Elevator, len(source))
	for k, v := range source {
		newMap[k] = v
	}
	return newMap
}
// func sendElevToMaster(isMaster bool, elev elevator.Elevator, sendMyselfToMaster chan elevator.Elevator) {
// 	if !isMaster {
		
// 		sendMyselfToMaster <- elev
// 		// tcpConn, _ := establish_connection.TransmitConn(masterPort, elev.ElevID) //vi må ha masterIp
// 		// tcp.Transmit(tcpConn, elev)
// 	}
// }


// func stop_functionality(stop bool, elev elevator.Elevator) {
// 	if stop {
// 		elevio.SetMotorDirection(elevio.MD_Stop)
// 		elevio.SetStopLamp(true)
// 	} else {
// 		elevio.SetMotorDirection(elev.Dirn)
// 		elevio.SetStopLamp(false)
// 	}
// }

