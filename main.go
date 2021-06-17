package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/superc03/colirc/communication"
)

var (
	client_port int
)

func init() {
	if in, err := strconv.Atoi(os.Getenv("COLIRC_CLIENT_PORT")); err != nil {
		client_port = 6667
	} else {
		client_port = in
	}
}

func main() {
	// Logger
	l := log.New(os.Stdout, "[ColIRC] ", log.LstdFlags)

	// Client2Server TCP Connection
	clientServer := communication.NewServer(client_port, l)
	l.Printf("Server Started on Port %d :)", client_port)

	// Wait for Shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// Stop Servers
	<-sigChan
	l.Println("Received Termination Signal. Shutting Down")
	clientServer.Stop()
	l.Println("Server Shutdown Completed. Have a Nice Day :)")
}
