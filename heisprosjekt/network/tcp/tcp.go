package tcp

import (
	"Heis/elevator"
	"Heis/network/localip"
	"Heis/network/peers"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"
)

/**
 * This file contain functionality regarding TCP
 */

/**
 * @struct TaggedJson is the struct we wish to convert the desired data into before it is transmitted.
 * Reason: Enables easier transmission of different data types
 */
type TaggedJson struct {
	Type string
	JSON []byte
}

type Person struct {
	Name string
	Age  int
}

type House struct {
	Name  string
	Price int
}

type Elevator struct {
	Id    int
	Queue []int
}

/**
 * @var variables regarding the connections for the TCP.
 */
const (
	CONN_HOST = "localhost"
	CONN_PORT = 8081
	CONN_TYPE = "tcp"
)

/**
 * @func Transmit makes a connection, encodes the desired data to json, and writes it to the destination on address with port.
 */
func Transmit(port int, address string, data interface{}) {
	conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		fmt.Printf("[error] Failed to Dial: %v\n", err)
		return
	}
	defer conn.Close()

	buffer, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("[error] Failed to encode data with error: %v\n", err)
		return
	}

	buffer, err = json.Marshal(TaggedJson{reflect.TypeOf(data).Name(), buffer})
	if err != nil {
		fmt.Printf("[error] Failed to make buffer with error:")
	}

	_, err = conn.Write(buffer)
	if err != nil {
		fmt.Printf("[error] Failed to write: %v\n", err)
		return
	}

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("[error] Failed to read: %v\n", err)
		return
	}
	fmt.Printf("Read: %v\n", string(buffer[0:n]))
}

/**
 * @func Receive is listening for incoming data, and decodes the transmitted data from json to the data originally fed into the transmitter.
 */
func Receive(port int, address string, data ...interface{}) {
	listener, err := net.Listen(CONN_TYPE, fmt.Sprintf("%s:%d", address, port))

	channels := make(map[string]interface{}) // a map with called channels with each data's type, written as a string, as keys

	for _, channel := range data {
		if channel == nil || reflect.TypeOf(channel).Kind() != reflect.Chan {
			panic("Arguments contains one or more non channel type\n")
		}

		channels[reflect.TypeOf(channel).Elem().Name()] = channel
	}

	if err != nil {
		fmt.Printf("[error] Failed to create listener with error: %v\n", err)
		return
	}

	var tj TaggedJson
	buffer := make([]byte, 1024)
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("[error] Failed to create connection with error: %v\n", err)
			continue
		}

		length, err := conn.Read(buffer)

		if err != nil {
			fmt.Printf("[error] Failed to read from connection with error: %v\n", err)
			continue
		}

		err = json.Unmarshal(buffer[0:length], &tj)

		if err != nil {
			fmt.Printf("[error] Failed to marshal JSON with error: %v", err)
			continue
		}

		channel, ok := channels[tj.Type]

		if !ok {
			fmt.Printf("[warning] Recieved type we are not listening to: %v\n", tj.Type)
			continue
		}

		value := reflect.New(reflect.TypeOf(channel).Elem())

		err = json.Unmarshal(tj.JSON, value.Interface()) // lagrer dataen vi har fått fra transmitter. Vi lagrer dette i value

		if err != nil {
			fmt.Printf("[error] Failed to unmarshal data with error code: %v\n", err)
			continue
		}

		// Actually send data to the respective channel
		reflect.Select([]reflect.SelectCase{{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(channel),
			Send: reflect.Indirect(value),
		}})

		conn.Write([]byte("ACK"))

		conn.Close()
	}
}

func CanConnectToMaster(address string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), timeout*time.Second)
	if err != nil {
		fmt.Printf("[error] Error establishing tcp-connection %v\n", err)
		return false
	}
	conn.Close()
	return true
}

// ----------------- [Litt endring] ---------------
func Transmit2pkt0(port int, address string, dataCh <-chan elevator.Elevator) {
	conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		fmt.Printf("[error] Failed to Dial: %v\n", err)
		return
	}
	defer conn.Close()

	for data := range dataCh {
		buffer, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("[error] Failed to encode data with error: %v\n", err)
			continue
		}

		buffer, err = json.Marshal(TaggedJson{reflect.TypeOf(data).Name(), buffer})
		if err != nil {
			fmt.Printf("[error] Failed to make buffer with error: %v\n", err)
			continue
		}

		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Printf("[error] Failed to write: %v\n", err)
			continue
		}

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("[error] Failed to read: %v\n", err)
			continue
		}
		fmt.Printf("Read: %v\n", string(buffer[0:n]))
	}
}

func Receive2pkt0(port int, data ...interface{}) {
	localIp, err := localip.LocalIP()
	if err != nil {
		fmt.Printf("[error] Failed to find local IP: %v\n", err)
		return
	}

	listener, err := net.Listen(CONN_TYPE, fmt.Sprintf("%s:%d", localIp, port))

	channels := make(map[string]interface{}) // a map with called channels with each data's type, written as a string, as keys

	for _, channel := range data {
		if channel == nil || reflect.TypeOf(channel).Kind() != reflect.Chan {
			panic("Arguments contains one or more non channel type\n")
		}

		channels[reflect.TypeOf(channel).Elem().Name()] = channel
	}

	if err != nil {
		fmt.Printf("[error] Failed to create listener with error: %v\n", err)
		return
	}

	var tj TaggedJson
	buffer := make([]byte, 1024)
	for {
		conn, err := listener.Accept()

		//her burde det vel være en handler som tar inn conn?

		if err != nil {
			fmt.Printf("[error] Failed to create connection with error: %v\n", err)
			continue
		}

		length, err := conn.Read(buffer)

		if err != nil {
			fmt.Printf("[error] Failed to read from connection with error: %v\n", err)
			continue
		}

		err = json.Unmarshal(buffer[0:length], &tj)

		if err != nil {
			fmt.Printf("[error] Failed to marshal JSON with error: %v", err)
			continue
		}

		channel, ok := channels[tj.Type]

		if !ok {
			fmt.Printf("[warning] Recieved type we are not listening to: %v\n", tj.Type)
			continue
		}

		value := reflect.New(reflect.TypeOf(channel).Elem())

		err = json.Unmarshal(tj.JSON, value.Interface()) // lagrer dataen vi har fått fra transmitter. Vi lagrer dette i value

		if err != nil {
			fmt.Printf("[error] Failed to unmarshal data with error code: %v\n", err)
			continue
		}

		// Actually send data to the respective channel
		reflect.Select([]reflect.SelectCase{{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(channel),
			Send: reflect.Indirect(value),
		}})

		conn.Write([]byte("ACK"))

		conn.Close()
	}
}

// To varianter av Receive og Transmit som bare returnerer en connection

// var connectionsMutex sync.Mutex

func ReceiveConn(port string, connectionsCh chan<- map[string]net.Conn) (net.Conn, error) {
	// Resolve address
	addr, err := net.ResolveTCPAddr("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return nil, err
	}
	// Create listener
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return nil, err
	}
	defer listener.Close()

	fmt.Println("Server listening on", addr.String())

	connections := make(map[string]net.Conn)

	// Accept incoming connections
	buffer := make([]byte, 1024)
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Accepted connection on port: " + port)

		// [Her leser jeg slavens ID]
		k, err := tcpConn.Read(buffer)
		if err != nil {
			fmt.Printf("[error] Failed to read: %v\n", err)
			return nil, err
		}
		id := string(buffer[0:k])
		// --------------------------------
		//fmt.Printf("Read: %v\n", id)
		//fmt.Printf("Sjå på her da da: %v\n", connections)

		// Update connections map
		//HandleConn(tcpConn, id, &connections, connectionsCh)
		connections[id] = tcpConn
		connectionsCh <- connections

	}
}

// func HandleConn(conn net.Conn, id string, connections *map[string]net.Conn, connectionsCh chan<- map[string]net.Conn) {
// 	//remoteAddr := conn.RemoteAddr().String() // foreløpig er bare keyen Ip-adressen til remote connection. Burde være peer
// 	if *connections == nil {
// 		*connections = make(map[string]net.Conn)
// 	}
// 	connectionsMutex.Lock()
// 	defer connectionsMutex.Unlock()
// 	fmt.Printf("Sjå på her da: %v\n", *connections)
// 	(*connections)[id] = conn
// 	fmt.Printf("Sjå på her: %v\n", *connections)
// 	connectionsCh <- *connections
// 	//return connections
// }

func TransmitConn(port string, id string) (net.Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:"+port) // dette er bare foreløpig. IP-som mates inn må være destinasjonen
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return nil, err
	}
	tcpConn, err := net.Dial(CONN_TYPE, addr.String())
	if err != nil {
		fmt.Printf("[error] Failed to Dial: %v\n", err)
		return nil, err
	}
	tcpConn.Write([]byte(id))
	return tcpConn, err
}

func AddConnections(id string, connsUpdateCh chan map[string]net.Conn, peerchan chan peers.PeerUpdate, connsCh chan map[string]net.Conn) {
	var conns map[string]net.Conn
	for {
		select {
		case c := <-connsUpdateCh:
			conns = c
			//fmt.Println("Updated connections:", c)
			//fmt.Println("Connections lagt til:", conns)
			connsCh <- conns
		case c := <-peerchan:
			if c.Master == "" {
				c.Master = id
			}
			peers.PrintUpdatedPeers(c)
			//fmt.Println("Peer update:", c)
			for i := 0; i < len(c.Lost); i++ {
				for k := range conns {
					if k == c.Lost[i] {
						delete(conns, k)
					}
				}
			}
			//fmt.Printf("Connections update fra peer: %v\n", conns)
			connsCh <- conns

		}
	}
}

func SendAndReceive(connsCh chan map[string]net.Conn) {
	var conns map[string]net.Conn
	for {
		select {
		case c := <-connsCh:
			conns = c
			//fmt.Println("Send and receive sin conn liste", c)
			fmt.Println("Send and receive sin conn liste", conns)
			fmt.Println()
		}
	}
}
