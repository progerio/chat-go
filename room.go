package main

import "net"

type Room struct {
	name    string
	members map[net.Addr]*Client
}

func (r *Room) broadcast(sender *Client, msg string) {
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}
