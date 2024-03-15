package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/fsm"
	"Heis/master"
	"os"
	"os/exec"
	"time"

	//"Heis/master"
	"Heis/timer"

	// "Heis/network/bcast"

	"Heis/network/establish_connection"
	"Heis/network/netfuncs"
	"Heis/network/peers"
	"Heis/slave"
	"fmt"
	"net"
	// "Heis/network/tcp"
)

//_________________________________________________________________________________________________

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

//_________________________________________________________________________________________________

func main() {

	//_________________________________________________________________________________________________

	// go primary(backup())

	//Heisann35!
	masterPort := "8070"
	//slavePort := "8090"

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)
	//elevio.Init("localhost:15654", _numFloors)
	// elevio.Init("localhost:15655", _numFloors)

	

	// Elevator initialization
	eleviId := netfuncs.InitNet()

	buffer := 10
	// Channels to update the the FSM-functions:
	newOrderCh := make(chan map[string]elevator.Elevator, buffer)
	elevUpdateRealtimeCh := make(chan elevator.Elevator, buffer)

	// Channels for door-timer and light:
	doorTimerCh := make(chan bool, buffer)
	timedOut := make(chan int, buffer)
	lightsCh := make(chan int, buffer)

	// Channels for inputs:
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// Master channels:
	masterConnCh := make(chan net.Conn, buffer)
	connectionsCh := make(chan map[string]net.Conn, buffer)
	sendMasterIdToReceive := make(chan string, buffer)
	masterIdToAlertMasterCh := make(chan string, buffer)
	masterIdToSendAndReceiveToMasterCh := make(chan string, buffer)
	sendMapToSlavesCh := make(chan map[string]elevator.Elevator, buffer)
	getElevFromSlaveRx := make(chan elevator.Elevator, buffer)
	isMasterCh1 := make(chan bool, buffer)
	isMasterCh2:= make(chan bool, buffer)
	isMasterCh3:= make(chan bool, buffer)


	// Slave channels:
	slaveConnCh := make(chan net.Conn, buffer)
	sendMyselfToMasterTx := make(chan elevator.Elevator, buffer)
	receiveMapFromMasterCh := make(chan map[string]elevator.Elevator, buffer)


	// Channels for Heartbeat
	peerUpdateCh1 := make(chan peers.PeerUpdate)
	peerUpdateCh2 := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	// Goroutines for inputs
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Goroutines for Heartbeat
	go peers.Transmitter(15623, eleviId, peerTxEnable)
	go peers.Receiver(15623, peerUpdateCh1)
	go peers.PeerUpdates(eleviId, peerUpdateCh1, peerUpdateCh2, isMasterCh1, isMasterCh2, isMasterCh3, sendMasterIdToReceive, masterIdToAlertMasterCh)

	/* if <-isMasterCh1{
		isMaster = true
	} else {
		isMaster = false
	} */

	// Goroutines for establishing connection
	/* go establish_connection.EstablishConnToSlaves(eleviId, masterPort, masterConnCh, connectionsCh,
		sendMasterIdToReceive) */
	go establish_connection.EstablishConnToSlaves2(eleviId, masterPort, masterConnCh, connectionsCh, sendMasterIdToReceive, isMasterCh2)
	go slave.AlertMaster(masterPort, eleviId, masterIdToAlertMasterCh, masterIdToSendAndReceiveToMasterCh, slaveConnCh)
	
	go fsm.ButtonsAndRequests(masterPort, eleviId, /*isMaster,*/ elevUpdateRealtimeCh,
		drv_buttons, sendMapToSlavesCh, getElevFromSlaveRx, receiveMapFromMasterCh,
		newOrderCh, lightsCh, sendMyselfToMasterTx, isMasterCh3)

	go fsm.FloorObstrStop(masterPort, eleviId, elevUpdateRealtimeCh, drv_floors,
		newOrderCh, doorTimerCh, timedOut, lightsCh, sendMyselfToMasterTx)

	go timer.Timer(doorTimerCh, timedOut)

	go slave.SendAndReceiveToMaster(eleviId, slaveConnCh, masterIdToSendAndReceiveToMasterCh, receiveMapFromMasterCh, sendMyselfToMasterTx)

	go master.SendAndReceiveToSlaves(eleviId, peerUpdateCh2, masterConnCh, connectionsCh, sendMapToSlavesCh, getElevFromSlaveRx)

	select {}
}

//_________________________________________________________________________________________________
// Process pairs


func primary(counter int) {
	var addr string = "localhost:8070"
	fmt.Printf("This is now a primary.")
 	exec.Command("gnome-terminal", "--", "go", "run", "main.go"/*, "-isMaster=true"*/).Run() // Commentet -isMaster=true so that the elevator initializes as a slave.
	for {
		conn, err := net.Dial("udp", addr)
		if err != nil {
			fmt.Println("The following error occured", err)
		}
		conn.Write([]byte(""))
		time.Sleep(20*time.Millisecond)
	}
}

func backup() int {
	var addr string = "localhost:8070"
	println("This is a process pair backup.")
	udpConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	buffer := make([]byte, 1024)
	for {
		err = udpConn.SetReadDeadline(time.Now().Add(4*time.Second))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		_, _, err := udpConn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return 0
		}
	}
}

