package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/network/localip"
	"Heis/network/tcp"
	"flag"
	"fmt"
	//"os"
)

// this was just to make main smaller. Just temporary
/* func CheckId(id string) string {
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		return id
	}
	return id
} */

func main() {
	//kalosj9

	//-----
	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)

	//Initialiserer en heisstruct
	elev := elevator.InitElev()
	//------

	//Channels
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	elevatorRx := make(chan elevator.Elevator) // TCP-receive
	elevatorTx := make(chan elevator.Elevator) // TCP-transmit
	//heartbeat := make(chan )

	//Threads
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//temp:
	localip, err := localip.LocalIP()
	if err != nil {
		fmt.Printf("[error] Failed to get localIP %v\n", err)
		return
	}

	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

	fmt.Printf("Is Master: %v\n", isMaster)

	if isMaster {
		//fmt.Printf("Her1\n")
		go tcp.Receive2pkt0(8080, elevatorRx) // starter en receiver for lytting
		for {
			select {
			case data := <-elevatorRx:
				fmt.Printf("Dette kommer fra slaven\n %v\n", data)
				// data er hele heisen her. Den må brukes videre
			}
		}
	} else {

		tcp.Transmit2pkt0(8080, localip, elevatorTx)
		elevatorTx<-elev
	}


/* 	select {
	case a := <-drv_buttons:
		fmt.Printf("\ndrv_buttons%v\n", a)

	case a := <-drv_floors:
		fmt.Printf("\ndrv_floors%v\n", a)

	case a := <-drv_stop:
		fmt.Printf("\ndrv_stop%v\n", a)

	case a := <-drv_obstr:
		fmt.Printf("\ndrv_obs %v\n", a) */
		
	
	
	//felles for alle casene over: lage kopi av elev og send det til transmitter
		
	//case a:=<-elevatorRx:
		// ta imot ny hallrequst bytt den ut med den gamle
		//
		//

}


/*




 */
