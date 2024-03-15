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
	"flag"
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
	// elevio.Init("localhost:15654", _numFloors)
	// elevio.Init("localhost:15655", _numFloors)

	// Master false by default
	var (
		isMaster bool
	)

	flag.BoolVar(&isMaster, "isMaster", false, "")
	flag.Parse()

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

	// Slave channels:
	slaveConnCh := make(chan net.Conn, buffer)
	sendMyselfToMasterTx := make(chan elevator.Elevator, buffer)
	receiveMapFromMasterCh := make(chan map[string]elevator.Elevator, buffer)

	// Blocking channels:
	connEstablishedForSlave := make(chan struct{})
	ListenAccepted := make(chan struct{})

	// Channels for Heartbeat
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	// Goroutines for inputs
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Goroutines for Heartbeat
	go peers.Transmitter(15623, eleviId, peerTxEnable)
	go peers.Receiver(15623, peerUpdateCh)
	go peers.PeerUpdates(peerUpdateCh, sendMasterIdToReceive, masterIdToAlertMasterCh)

	fmt.Println("[main] Nå er jeg på vei inn i EstablishConnToSlaves")
	go establish_connection.EstablishConnToSlaves(eleviId, masterPort, masterConnCh, connectionsCh,
		sendMasterIdToReceive, ListenAccepted)
	fmt.Println("[main] Nå har jeg kommet meg forbi EstablishConnToSlaves")
	go slave.AlertMaster(masterPort, eleviId, masterIdToAlertMasterCh, masterIdToSendAndReceiveToMasterCh, slaveConnCh, connEstablishedForSlave)
	fmt.Println("[main] Nå har jeg nådd sperren for !isMaster")

	/* if !isMaster {
		<-connEstablishedForSlave
	} */

	fmt.Println("[main] Started!")
	go fsm.ButtonsAndRequests(masterPort, eleviId, isMaster, elevUpdateRealtimeCh,
		drv_buttons /*elevatorTx, */, sendMapToSlavesCh, getElevFromSlaveRx, receiveMapFromMasterCh,
		newOrderCh, lightsCh, sendMyselfToMasterTx)

	go fsm.FloorObstrStop(masterPort, isMaster, eleviId, elevUpdateRealtimeCh, drv_floors,
		newOrderCh, doorTimerCh, timedOut, lightsCh, sendMyselfToMasterTx)

	//go fsm.FSM(drv_buttons, drv_floors, drv_stop, drv_obstr, eleviId, isMaster, ElevatorTx, CostTx, masterOrders, ElevatorRx, CostRx)
	//go netfuncs.Network_FSM(peerUpdateCh)
	go timer.Timer(doorTimerCh, timedOut)

	// Slave
	go slave.SendAndReceiveToMaster(eleviId, slaveConnCh, masterIdToSendAndReceiveToMasterCh, receiveMapFromMasterCh, sendMyselfToMasterTx)

	// Master
	fmt.Println("[main] Nå har vi nådd sperren for master")
	/* if isMaster {
		fmt.Println("[main] Jeg er master, og har nådd sperren")
		<-ListenAccepted
		fmt.Println("[main] Jeg er master, og har kommet meg")

	} */
	fmt.Println("[main] Nå er vi kommet forbi sperren for master")
	go master.SendAndReceiveToSlaves(masterConnCh, connectionsCh, sendMapToSlavesCh, getElevFromSlaveRx)
	// // go tcp.Receive(masterPort, eleviId, elevatorRx)
	// for {
	// 	select {
	// 		case c := <-sendMyselfToMasterTx:
	// 			fmt.Printf("[main] mottok denne heisen fra en slave: %v\n", c)
	// 	case c := <-newOrderCh:
	// 			fmt.Printf("[main] mottok denne ordren: %v\n", c)
	// 		}
		// }

	select {}
}

//_________________________________________________________________________________________________
// Process pairs


func primary(counter int) {
	var addr string = "localhost:8070"
	fmt.Printf("This is now a primary.")
 	exec.Command("gnome-terminal", "--", "go", "run", "main.go", "-isMaster=true").Run()
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

