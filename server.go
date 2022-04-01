package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	rooms    map[string]*Room
	commands chan Command
}

func newServer() *Server {
	return &Server{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
	}
}

func (s *Server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *Server) newClient(conn net.Conn) {
	log.Printf("new Client has connected: %s", conn.RemoteAddr().String())

	c := &Client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *Server) nick(c *Client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("All right, I will all you %s", c.nick))
}
func (s *Server) join(c *Client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]

	if !ok {
		r = &Room{
			name:    roomName,
			members: make(map[net.Addr]*Client),
		}

		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", r.name))
}

func (s *Server) listRooms(c *Client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms are: %s", strings.Join(rooms, ", ")))
}

func (s *Server) msg(c *Client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}
	c.room.broadcast(c, c.nick+":"+strings.Join(args[1:], " "))
}

func (s *Server) quit(c *Client, args []string) {
	log.Printf("client has disconnected : %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.msg("sad to see you go :( ")
	c.conn.Close()
}

func (s *Server) quitCurrentRoom(c *Client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left thre room", c.nick))
	}
}
