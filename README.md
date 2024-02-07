# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering


Onsdag kveld: 
Vi bør lage en funksjon som endrer både motordirection og elevatiour direction. 

feilen vår lå i at vi ikke endret direction på motoren da vi endret elevator direction.

/*
funksjonen vil se slik ut, bare ikke skriv den i elevator modulen. skriv den heller i requests. requests har allerede importert elevator, så vi ønsker ikke å importere requests i elevator. 
func Elevator_directionChange(elev *Elevator,pair requests.DirnBehaviourPair){
	elev.Dirn = pair.Dirn
	elev.Behaviour = pair.Behaviour
	elevio.SetMotorDirection(elev.Dirn)
}
*/

Hvis vi trykker en knapp som er i samme etasjen som heisen befinner seg i, så får ikke heisen til å ta nye ordre senere. låser seg fast. ellers funker den greit (tror vi). 




--------------------JONATHAN SITT--------------------------
___________MAIN.GO____________________
package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/fsm"
	"fmt"
)

func main() {

	numFloors := 4
	numButtons := 3
	elevio.Init("localhost:15657", numFloors)

	//Initialiserer en heisstruct
	e elevator.InitElev()

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	/* Lage en tilsvarende funksjon som oppdaterer
		heisstructen?
	go elevator.status(drv_buttons, drv_floors,
						drv_obstr, drv_stop)*/

	if elevio.GetFloor() == -1 {
		fsm.Fsm_onInitBetweenFloors()
	}

	fmt.Printf("Started!\n")

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			fsm.Fsm_onRequestButtonPress(e, a.Floor, a.Button)

		case a := <-drv_floors: // a er etasjen heisen er i
			fmt.Printf("%+v\n", a)
			fsm.Fsm_onFloorArrival(a)

			/*if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}*/
			elevio.SetMotorDirection(d)

			// Fikser disse funksjonene senere
			/*
				case a := <-drv_obstr:
					fmt.Printf("%+v\n", a)
					if a {
						elevio.SetMotorDirection(elevio.MD_Stop)
					} else {
						elevio.SetMotorDirection(d)
					}

				case a := <-drv_stop:
					fmt.Printf("%+v\n", a)
					for f := 0; f < numFloors; f++ {
						for b := elevio.ButtonType(0); b < 3; b++ {
							elevio.SetButtonLamp(b, f, false)
						}
					}
			*/
		}
	}

}




_____________HUSKERIKKE.GO_________________
import (
	"fmt"
	// "Heis/noe noe elevator_io_types ??"
)
```
func fsm_onRequestButtonPress(btn_floor int, btn_type Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.toString())
	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			timer_start(elevator.config.doorOpenDuration_s)
		} else {
			elevator.requests[btn_floor][btn_type] = 1
		}
	
	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = 1
	
	case EB_idle:
		elevator.requests[btn_floor][btn_type] = 1
		DirnBehaviourPair pair = requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
			switch pair.behaviour {
			case EB_DoorOpen:
				outputDevice.doorLight(1)
				timer_start(elevator.config.doorOpenDuration_s)
				elevator = requests_clearAtCurrentFloor(elevator)

			case EB_Moving:
				outputDevice.motorDirection(elevator.dirn)
			
			case EB_Idle:
		}
	}

	setAllLights(elevator)

	fmt.Println("\nNew state:")
```


---------------------------------ANDERS SITT-------------------------------

type DirnBehaviourPair struct {
	Dirn              dirn
	ElevatorBehaviour behaviour
}


```
func fsm_onFloorArrival(newFloor int) int{
	//printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
    //elevator_print(elevator);

	elevator.floor = newFloor

	outputDevice.floorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(	elevator) {
			output.motorDirection(D_Stop)
			outputDevice.doorLight(1);
            elevator = requests_clearAtCurrentFloor(elevator);
            //timer_start(elevator.config.doorOpenDuration_s);
            setAllLights(elevator);
            elevator.behaviour = EB_DoorOpen;
		}
	default:

	}
	//printf("\nNew state:\n");
    //elevator_print(elevator); 

}
```

//Trenger vi å bruke pekere her? (*elevator)?


```
func fsm_onDoorTimeout() { //????
	//printf("\n\n%s()\n", __FUNCTION__);
    //elevator_print(elevator);

	switch(elevator.behaviour){
	case EB_DoorOpen:
		//timer_start(elevator.config.doorOpenDuration_s);
		elevator = requests_clearAtCurrentFloor(elevator)
		setAllLights(elevator)
	case EB_Moving:
	case EB_Idle:
		outputDevice.doorLight(0)
		outputDevice.motorDirection(elevator.dirn)
	default:

	}

	//printf("\nNew state:\n");
    //elevator_print(elevator);


}
```
```
func Fsm_onFloorArrival(newFloor int){
	//printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
    //elevator_print(elevator);
	e.Floor = newFloor

	elevio.GetFloor()

	switch e.Behaviour{
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(e){
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			e = requests.Requests_clearAtCurrentFloor(e)
			time.Sleep(time.Duration(e.DoorOpenDuration_s)*time.Second)
			setAllLights(e) // need to rewrite this function
			e.Behaviour = elevator.EB_DoorOpen
		}
	default:
	}

}
```








 
	elevator_print(elevator)
}
