package main

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vprates-22/chat-in-go/server/internal"
)

func main() {
	channel := internal.NewChannel()
	shutdown := make(chan struct{})

	var wg sync.WaitGroup

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		println("Shutting down server...")

		channel.Shutdown()

		close(shutdown)

		listener.Close()
	}()

	println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-shutdown:
				println("Server is shutting down, stop accepting new connections.")
				wg.Wait()
				println("Server shutdown complete.")
				return
			default:
				println("Error accepting connection:", err.Error())
				continue
			}
		}

		wg.Add(1)
		go internal.HandleConnection(conn, &wg, channel)
	}
}
