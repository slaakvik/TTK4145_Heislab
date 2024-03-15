package fsm

import (
	"Heis/cost_fns"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/requests"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Mutex to synchronize access to shared map
var mapMutex sync.Mutex

// Finite state machine

func ButtonsAndRequests(masterPort string, elevatorID string, isMaster bool, elevUpdateRealtimeCh <-chan elevator.Elevator,
	drv_buttons chan elevio.ButtonEvent, sendMapToSlavesCh chan<- map[string]elevator.Elevator, getElevFromSlave chan elevator.Elevator,
	receiveMapFromMasterCh <-chan map[string]elevator.Elevator, newOrderCh chan<- map[string]elevator.Elevator,
	lightsCh <-chan int, sendMyselfToMaster chan elevator.Elevator) {
	// fmt.Println("[ButtonsAndRequests] akkurat kommet inni")
	elev := elevator.InitElev()
	elev.ElevID = elevatorID

	// Initialize mapOfElevs
	mapOfElevs := make(map[string]elevator.Elevator)
	mapOfElevs[elev.ElevID] = elev

	fmt.Println("______________")

	for {
		select {
		case a := <-elevUpdateRealtimeCh:
			// fmt.Printf("[ButtonsAndRequests] mottok en heis på elevUpdateRealtimeCh: %v\n", a)
			mapMutex.Lock()
			elev = a
			mapOfElevs[elev.ElevID] = elev
			go writeToLocalBackup(elevator.GetCabCalls(elev))
			mapMutex.Unlock()
			// fmt.Printf("[ButtonsAndRequests] mapOfElevs ser slik ut nå: %v\n", mapOfElevs)
		case a := <-drv_buttons:
			btn_floor := a.Floor
			btn_type := a.Button
			fmt.Printf("[ButtonsAndRequests] Button: %+v\n", a)

				if isMaster {
					mapMutex.Lock()
					elev.Requests[btn_floor][btn_type] = true
					mapOfElevs[elev.ElevID] = elev
					mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
					mapMutex.Unlock()
					sendMapToSlavesCh <- mapOfElevs
					// fmt.Println("[ButtonsAndRequests] Sendte til alle")
					newOrderCh <- mapOfElevs

				} else {
					mapMutex.Lock()
					elev.Requests[btn_floor][btn_type] = true
					mapOfElevs[elev.ElevID] = elev
					mapMutex.Unlock()
					go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
				}
			
		case a := <-getElevFromSlave:
			// fmt.Println("[ButtonsAndRequests] mottok noe på getElevFromSlave")
			if isMaster {
				// fmt.Printf("[ButtonsAndRequests] Jeg mottok en heis nå: %v\n", a)
				mapMutex.Lock()
				mapOfElevs[a.ElevID] = a
				mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
				mapMutex.Unlock()
				sendMapToSlavesCh <- mapOfElevs
				newOrderCh <- mapOfElevs
			}

		case a := <-receiveMapFromMasterCh:
			// fmt.Println("[ButtonsAndRequests] mottok noe på receiveMapFromMasterCh")
			if !isMaster {
				// fmt.Printf("[ButtonsAndRequests] Jeg mottok et map nå: %v\n", a)
				mapMutex.Lock()
				mapOfElevs = a
				mapMutex.Unlock()
				// fmt.Println("[ButtonsAndRequests] skal sende mappet mapOfElevs til newOrderCh")
				newOrderCh <- mapOfElevs
			}
		case <-lightsCh:
			fmt.Println("[ButtonsAndRequests] mottok noe på lightsCh")
			mapMutex.Lock()
			elevator.SetAllLights(elev, mapOfElevs)
			mapMutex.Unlock()
		}
	}
}

func FloorObstrStop(masterPort string, isMaster bool, elevatorId string, elevUpdateRealtimeCh chan<- elevator.Elevator, drv_floors chan int, 
	newOrderCh <-chan map[string]elevator.Elevator, doorTimerCh chan bool, timedOut chan int, lightsCh chan<- int, sendMyselfToMaster chan elevator.Elevator) {
	elev := elevator.InitElev()
	elev.ElevID = elevatorId
	// fmt.Println("[FloorObstrStop] kommet meg inni")
	if elevio.GetFloor() == -1 {
		elev = elevator.OnInitBetweenFloors(elev)
		elevUpdateRealtimeCh <- elev
		elevator.Elevator_print(elev)
	}
	// fmt.Println("[FloorObstrStop] kommet til for-loopen")
	for {
		select {
		case a := <-newOrderCh:
			fmt.Println("[FloorObstrStop] mottok noe på newOrderCh")
			elev = a[elev.ElevID]
			if !isMaster {
				// fmt.Println("slave received new orders.")
				// fmt.Println("slaves elevator: ", elev)
				// elev.Requests[0][2]= true
			}
			if elev.Behaviour != elevator.EB_Moving && requests.ShouldClearImmediately(elev) {
				fmt.Println("inni should clear")
				
				elev.Behaviour = elevator.EB_DoorOpen
				// mapOfElevs[elev.ElevID] = elev
				doorTimerCh <- true
			} else if elev.Behaviour == elevator.EB_Idle {
				fmt.Println("on request while idle")
				pair := requests.ChooseDirection(elev)
				elev.Dirn = pair.Dirn
				elevio.SetMotorDirection(elev.Dirn)
				elev.Behaviour = pair.Behaviour
				
			}
			lightsCh <- 1
			elevUpdateRealtimeCh <- elev
			SendElevToMaster(isMaster, elev, sendMyselfToMaster) // på tråd?

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
			//  fmt.Printf("Stop button: %+v\n", a)
			//  stop_functionality(a, elev)

		}
	}
}

func SendElevToMaster(isMaster bool, elev elevator.Elevator, sendMyselfToMaster chan elevator.Elevator) {
	if !isMaster {
		fmt.Println("[SendElevToMaster] kommet meg inni, og jeg er ikke master ")
		fmt.Println("[SendElevToMaster] sender heisen til sendMyselfToMasterCh")
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


// Can be run from anywhere. Needs either cab calls, elev or similar as input
func writeToLocalBackup(cabCalls []bool) error {
	file, err := os.Create("localBackup.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	var boolStr strings.Builder
	for _, value := range cabCalls {
		if value {
			boolStr.WriteString("true")
		} else {
			boolStr.WriteString("false")
		}
	}

	_, err = file.WriteString(boolStr.String())
	if err != nil {
		return err
	}
	return nil
}