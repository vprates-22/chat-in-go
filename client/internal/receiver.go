package internal

import (
	"bufio"
	"io"
	"net"
	"strings"
	"sync"
)

func RunReceiver(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()

	reader := bufio.NewReader(conn)

	for {
		message, err := readMessage(reader)
		if err != nil {
			return
		}
		print(message + "\n> ")
	}
}

func readMessage(reader *bufio.Reader) (string, error) {
	message, err := reader.ReadString('\x00')
	if err == io.EOF || err != nil {
		return "", err
	}

	return strings.TrimSpace(message), nil
}
