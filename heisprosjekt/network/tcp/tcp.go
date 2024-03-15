package tcp

import (
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

/**
 * @var variables regarding the connections for the TCP.
 */
const (
	CONN_HOST = "localhost"
	CONN_PORT = 8081
	CONN_TYPE = "tcp"
)

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

/**
 * @func Transmit for both the master and slaves
 */
func Transmit(conn net.Conn, data interface{}) {
	// fmt.Println("[Transmit] Nå har jeg akkurat kommet inni Transmit")
	buffer, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("[error] Failed to encode data with error: %v\n", err)
		return
	}
	// fmt.Println("[Transmit] Klarte å omforme data til Marshal")

	// fmt.Println("[Transmit] Skal til å lage en buffer")
	buffer, err = json.Marshal(TaggedJson{reflect.TypeOf(data).Name(), buffer})
	if err != nil {
		fmt.Printf("[error] Failed to make buffer with error:")
	}

	// fmt.Println("[Transmit] Skal til å skrive")
	_, err = conn.Write(buffer)
	if err != nil {
		fmt.Printf("[error] Failed to write: %v\n", err)
		//conn.Close() her?
		return
	}
	// fmt.Println("[Transmit] Hvis ingen error, fikk jeg skrevet meldingen")
	// fmt.Println("[Transmit] Nå går jeg ut")

}

/**
 * @func for the master and the slave
 */
func Receive(conn net.Conn, data ...interface{}) {
	// fmt.Println("[Receive] Akkurat kommet meg inni")
	defer conn.Close()                       // må denne være her??
	channels := make(map[string]interface{}) // a map with called channels with each data's type, written as a string, as keys

	for _, channel := range data {
		if channel == nil || reflect.TypeOf(channel).Kind() != reflect.Chan {
			panic("Arguments contains one or more non channel type\n")
		}

		channels[reflect.TypeOf(channel).Elem().Name()] = channel
	}
	// fmt.Println("[Receive] Fyllt ett map med ønskede kanaler")

	var tj TaggedJson
	buffer := make([]byte, 1024)
	// fmt.Println("[Receive] På vei inn i for-loopen")
	for {
		// fmt.Println("[Receive] Er i for-loopen, og leser")
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
	fmt.Println("[Receive] Connection kunne ikke brukes mer, går ut av denne funksjonen")
}
