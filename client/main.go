package main

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vprates-22/chat-in-go/client/internal"
)

func main() {
	if len(os.Args) != 3 {
		println("usage: go run client/main.go <address> <port>")
		return
	}

	address := os.Args[1]
	port := os.Args[2]

	serverAddress := address + ":" + port

	println("Client is connecting to server at:", serverAddress)

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		println("Error connecting to server:", err.Error())
		return
	}
	// defer conn.Close()

	println("Client is connected to:", serverAddress)

	var wg sync.WaitGroup

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigchan)

	go func() {
		<-sigchan
		println("Shutting down client...")
		conn.Close()
	}()

	wg.Add(1)
	go internal.RunReceiver(&wg, conn)
	go internal.RunSender(conn)

	wg.Wait()
}
