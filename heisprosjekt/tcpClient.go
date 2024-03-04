package main

import (
	//"Heis/master"
	//"Heis/network/localip"
	//"Heis/network/peers"
	//"Heis/slave"
	//"flag"
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/network/tcp"
	"time"
	//"os"
)

func main() {
	//gamle kalosj1!

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)

	//Initialiserer en heisstruct
	elev := elevator.InitElev()

	//var elevatorCh = make(chan elevator.Elevator)

	for {
		tcp.Transmit(tcp.CONN_PORT, tcp.CONN_HOST, elev)
		time.Sleep(2 * time.Second)

	}

	//p := tcp.Person{Name: "Anders", Age: 22}
	//h := tcp.House{Name: "Villa", Price: 1000}

	//tcp.Transmit(tcp.CONN_PORT, tcp.CONN_HOST, p)
	//tcp.Transmit(tcp.CONN_PORT, tcp.CONN_HOST, h)

}
