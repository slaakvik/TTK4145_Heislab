package establish_connection

import (
	"Heis/network/peers"
	"fmt"
	"net"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = 8081
	CONN_TYPE = "tcp"
)

/**
 * @func for the slave
 */
func TransmitConn(port string, id string, masterIp string) (net.Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", masterIp+":"+port) // dette er bare foreløpig. IP-som mates inn må være destinasjonen
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

/**
 * @func for the master
 */
func ReceiveConn(id string, port string, connCh chan net.Conn, connectionsCh chan<- map[string]net.Conn) (net.Conn, error) {
	masterIp := peers.ExtractIpFromPeer(id)
	// Resolve address
	addr, err := net.ResolveTCPAddr("tcp", masterIp+":"+port)
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

		connections[id] = tcpConn
		connCh <- tcpConn
		connectionsCh <- connections

	}
}

/**
 * @func
 */
func AddConnections(id string, connsUpdateCh chan map[string]net.Conn, peerchan chan peers.PeerUpdate, connsCh chan map[string]net.Conn) {
	var conns map[string]net.Conn
	for {
		select {
		case c := <-connsUpdateCh:
			conns = c
			//fmt.Println("Updated connections:", c)
			//fmt.Println("Connections lagt til:", conns)
			connsCh <- conns // Må ha dennne, ikke fjern den Anders! (kanskje?)
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
