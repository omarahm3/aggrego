package main

import "net"

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, s string) {
	for addr, m := range r.members {
		if addr == sender.conn.RemoteAddr() {
			continue
		}

		m.msg(s)
	}
}
