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

func FSM(buttons chan elevio.ButtonEvent, floors chan int, stop chan bool, obstr chan bool, elevatorID string) {
	elev := elevator.InitElev()

	if elevio.GetFloor() == -1 {
		elev = onInitBetweenFloors(elev)
		elevator.Elevator_print(elev)
	}
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
				setAllLights(elev)
				quit := make(chan int)
				go doorTimer(quit)
				<-quit
				elev = onDoorTimeout(elev)
			} else {
				//kall cost func med nye requests.
				//switch case som sjekker om heisen er idle eller moving/dooropen.
				// dersom heisen er moving eller dooropen skal request listen oppdateres utifra cost func
				//dersom heisen er idle skal den i tillegg velge retnig og sette motordirection etter å ha hentet requests fra costfunc.
				//finn en måte å hente ut riktig keys og values fra cost func, og flett inn de nye hall callsene inn i request listen til heisen.

				//Tar inn id som parameter, og sjekker
				if btn_type == 2 {
					elev.Requests[btn_floor][btn_type] = true
				}

				NewHRs := runCostFunc(elev, btn_floor, int(btn_type), elevatorID)
				MyNewHRs := NewHRs[elevatorID]
				switch elev.Behaviour {
				case elevator.EB_DoorOpen:
					elev.Requests = elevator.MergeHallAndRequests(elev.Requests, MyNewHRs)

				case elevator.EB_Moving:
					elev.Requests = elevator.MergeHallAndRequests(elev.Requests, MyNewHRs)

				case elevator.EB_Idle:
					elev.Requests = elevator.MergeHallAndRequests(elev.Requests, MyNewHRs)
					pair := requests.Requests_chooseDirection(elev)
					elev.Dirn = pair.Dirn
					elevio.SetMotorDirection(elev.Dirn)
					elev.Behaviour = pair.Behaviour
				}

			}
			setAllLights(elev)

		case a := <-floors: // a er etasjen heisen er i
			fmt.Printf("Floor: %+v\n", a)
			//onFloorArrival(&elev, a)
			elev.Floor = a
			elevio.SetFloorIndicator(elev.Floor)
			switch elev.Behaviour {
			case elevator.EB_Moving:
				if requests.Requests_shouldStop(elev) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.Behaviour = elevator.EB_DoorOpen
					setAllLights(elev)
					quit := make(chan int)
					go doorTimer(quit)
					<-quit
					elev = onDoorTimeout(elev)
				}
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

func runCostFunc(elev elevator.Elevator, btn_floor int, btn_type int, elevId string) map[string][][2]bool {
	tempElev := elev //egt unødvendig å ha tempElev, heisen vår endres kun i denne funksjonen hvis vi endrer den her, så vi likke påvirke heisen andre steder.
	tempElev.Requests[btn_floor][btn_type] = true
	input := cost_fns.InputToCost(tempElev, elevId)
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
