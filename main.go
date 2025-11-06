package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Anurag-S1ngh/redis/handler"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		go handler.HandleConnection(conn)
	}
}
