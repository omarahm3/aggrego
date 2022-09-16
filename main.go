package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var DEFAULT_PORT = "8888"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	s := newServer()
	go s.run()

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("unable to start tcp server: %s", err.Error())
	}

	defer l.Close()
	log.Printf("server started at :%s", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("unable to accept tcp connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
