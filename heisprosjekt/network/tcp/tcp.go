package tcp

import (
	"Heis/elevator"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
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

// type Person struct {
// 	Name string
// 	Age  int
// }

// type House struct {
// 	Name  string
// 	Price int
// }

// type Elevator struct {
// 	Id    int
// 	Queue []int
// }

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
/*func Transmit(port int, address string, data interface{}) {
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
*/
/**
 * @func Receive is listening for incoming data, and decodes the transmitted data from json to the data originally fed into the transmitter.
 */
/*
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
*/

/*
func CanConnectToMaster(address string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), timeout*time.Second)
	if err != nil {
		fmt.Printf("[error] Error establishing tcp-connection %v\n", err)
		return false
	}
	conn.Close()
	return true
}
*/

//--------------------------------------------------------------------------------------

/**
 * @func Transmit for both the master and slaves
 */
func Transmit(conn net.Conn, data interface{}) {

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
		//conn.Close() her?
		return
	}

}

/**
 * @func for the master and the slave
 */
func ReceiveHandler(conn net.Conn, data ...interface{}) {
	defer conn.Close()                       // må denne være her??
	channels := make(map[string]interface{}) // a map with called channels with each data's type, written as a string, as keys

	for _, channel := range data {
		if channel == nil || reflect.TypeOf(channel).Kind() != reflect.Chan {
			panic("Arguments contains one or more non channel type\n")
		}

		channels[reflect.TypeOf(channel).Elem().Name()] = channel
	}

	var tj TaggedJson
	buffer := make([]byte, 1024)
	for {

		length, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("[error] Failed to read: %v\n", err)
			break // this  makes the code break out of the for-loop if it fails to read the buffer (conn gone)
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
	}
	fmt.Println("Nå er denne tråden lukket")

}

func SendAndReceive(connCh chan net.Conn, connsCh chan map[string]net.Conn, transmitCh chan map[string]elevator.Elevator) {
	var conns map[string]net.Conn
	for {
		select {
		case c := <-connsCh:
			conns = c
			//fmt.Println("Send and receive sin conn liste", c)
			fmt.Println("Send and receive sin conn liste", conns)
			fmt.Println()
		case c := <-connCh:
			go ReceiveHandler(c /* her er kanalene som master skal lytte på*/)

		case c := <-transmitCh:
			for _, v := range conns {
				go Transmit(v, c) // kanskje ikke goroutine?
			}
		}
	}
}
