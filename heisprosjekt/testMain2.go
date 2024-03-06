package main

import (
	"Heis/network/tcp"
	"flag"
	"fmt"
)


func main() {
    var (
        isServer bool
    )

	tullekanalen := make(chan string)

    flag.BoolVar(&isServer, "IsServer", false, "")
    flag.Parse()

    fmt.Printf("Is server: %v", isServer)

	if isServer {
		go tcp.Receive(8080, "localhost", tullekanalen )
		for {
			select{
			case data := <-tullekanalen:
				fmt.Printf("Dette kommer fra tullekanalen %v\n", data)
			}
		}
	} else {
		tcp.Transmit(8080, "localhost", "tøys")
	}

}