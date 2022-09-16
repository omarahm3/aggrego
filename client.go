package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	CLIENT_CMD_NICK  = "/nick"
	CLIENT_CMD_JOIN  = "/join"
	CLIENT_CMD_ROOMS = "/rooms"
	CLIENT_CMD_MSG   = "/msg"
	CLIENT_CMD_QUIT  = "/quit"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Printf("error reading client message: %s", err.Error())
			return
		}

		msg = strings.Trim(msg, "\n\r")
		args := strings.Split(msg, " ")

		if len(args) == 0 {
			log.Print("message is empty")
			continue
		}

		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case CLIENT_CMD_NICK:
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case CLIENT_CMD_MSG:
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case CLIENT_CMD_JOIN:
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case CLIENT_CMD_ROOMS:
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}
		case CLIENT_CMD_QUIT:
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte(fmt.Sprintf("ERR: %s\n", err.Error())))
}

func (c *client) msg(s string) {
	c.conn.Write([]byte(fmt.Sprintf("> %s\n", s)))
}
