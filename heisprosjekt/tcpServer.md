package main

import (
	//"Heis/master"
	//"Heis/network/localip"
	//"Heis/network/peers"
	//"Heis/slave"
	//"flag"
	"Heis/elevator"
	//"Heis/driver-go/elevio"
	"Heis/network/tcp"
	"fmt"
	//"os"
)

func main() {

	var elevatorCh = make(chan elevator.Elevator)
	var personCh = make(chan tcp.Person)
	var houseCh = make(chan tcp.House)

	go tcp.Receive(tcp.CONN_PORT, tcp.CONN_HOST, elevatorCh)

	for {
		select {
		case data := <-personCh:
			fmt.Printf("[Received]: %v\n", data)
		case data := <-houseCh:
			fmt.Printf("[Received]: %v\n", data)
		case data := <-elevatorCh:
			fmt.Printf("Heisen: %v\n", data)
			//elevator.Elevator_print(data)
		}
	}
}
