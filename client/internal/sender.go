package internal

import (
	"bufio"
	"net"
	"os"
	"strings"
)

func RunSender(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	var cmd string
	for {
		print("> ")

		if !scanner.Scan() || scanner.Err() != nil {
			return
		}

		cmd = scanner.Text()
		parts := strings.Fields(cmd)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "/quit":
			conn.Write([]byte(cmd + " \x00"))
			conn.Close()
			return
		default:
			_, err := conn.Write([]byte(cmd + " \x00"))
			if err != nil {
				return
			}
		}
	}
}
