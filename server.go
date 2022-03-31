package main

import (
	"log"
	"net"
)

type Server struct {
	rooms map[string] *Room
	commands chan Command
}


func newServer() *Server{
	return &Server{
		rooms: make(map[string]*Room),
		commands: make(chan Command),
	}
}

func (s *Server) newClient(conn net.Conn){
	log.Printf("new Client has connected: %s", conn.RemoteAddr().String())

	 c := &Client{
		 conn: conn,
		 nick: "anonymous",
		 commands: s.commands,
	 }
}