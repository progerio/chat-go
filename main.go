package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()

	listener, err := net.Listen("tcp", ":8088")

	if err != nil {
		log.Fatalf("unable to start server %s", err.Error())
	}
	defer listener.Close()
	log.Printf("started server on :8088")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
