package main

import (
	"fmt"
	"log"
	"net"
)

const PORT = 8888

func main() {
	s := newServer()
	go s.run()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatalf("unable to start tcp server: %s", err.Error())
	}

	defer l.Close()
	log.Printf("server started at :%d", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("unable to accept tcp connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
