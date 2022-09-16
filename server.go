package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	errMissingArgs       = errors.New("there are missing args you dumbo! follow this pattern fool: [/cmd args...]")
	errClientIsNotJoined = errors.New("are you stupid? trying to send message to room when you're not joined yet?")
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client is connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "5wl",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) run() {
	log.Print("listenting for server commands")
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			log.Print("server received /nick command")
			s.nick(cmd)
		case CMD_JOIN:
			log.Print("server received /join command")
			s.join(cmd)
		case CMD_ROOMS:
			log.Print("server received /rooms command")
			s.listRooms(cmd)
		case CMD_MSG:
			log.Print("server received /msg command")
			s.msg(cmd)
		case CMD_QUIT:
			log.Print("server received /quit command")
			s.quit(cmd)
		}
	}
}

func (s *server) nick(cmd command) {
	if len(cmd.args) == 1 {
		cmd.client.err(errMissingArgs)
		return
	}

	cmd.client.nick = cmd.args[1]
	cmd.client.msg(fmt.Sprintf("your nickname is set to: %s", cmd.args[1]))
}

func (s *server) join(cmd command) {
	if len(cmd.args) == 1 {
		cmd.client.err(errMissingArgs)
		return
	}

	roomName := cmd.args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		// Create the room and add it to server rooms list
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	// Add client to list of room members
	r.members[cmd.client.conn.RemoteAddr()] = cmd.client
	// Make sure to remove client from old room
	s.quitCurrentRoom(cmd.client)
	// Add this room to client
	cmd.client.room = r

	// Notify all users
	r.broadcast(cmd.client, fmt.Sprintf("%s has joined this room, please make sure to roast them till they leave", cmd.client.nick))
	cmd.client.msg(fmt.Sprintf("hey %s, welcome to %s, i have notified all members to give you warm welcome :)", cmd.client.nick, r.name))
}

func (s *server) listRooms(cmd command) {
	var rooms []string
	for _, r := range s.rooms {
		rooms = append(rooms, r.name)
	}

	cmd.client.msg(fmt.Sprintf("available rooms:\n%s", strings.Join(rooms, ", ")))
}

func (s *server) msg(cmd command) {
	if cmd.client.room == nil {
		cmd.client.err(errClientIsNotJoined)
		return
	}

	msg := fmt.Sprintf("%s: %s", cmd.client.nick, strings.Join(cmd.args[1:len(cmd.args)], " "))
	cmd.client.room.broadcast(cmd.client, msg)
}

func (s *server) quit(cmd command) {
	log.Printf("client [%s] has disconnected", cmd.client.conn.RemoteAddr().String())
	s.quitCurrentRoom(cmd.client)
	cmd.client.msg("bye bye MF! i never been happier!")
	cmd.client.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	// In case user was in another room before joining this one
	// we need to get him/her the fuck out of it
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		// Broadcast to all users that this user has left
		c.room.broadcast(c, fmt.Sprintf("%s has left this damn room", c.nick))
	}
}

func newServer() *server {
	log.Print("initializing new server instance")
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}
