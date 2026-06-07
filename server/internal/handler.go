package internal

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"sync"
)

func HandleConnection(conn net.Conn, wg *sync.WaitGroup, channel *Channel) {
	defer wg.Done()
	defer conn.Close()

	id := generateId()
	println("User", id, "joined the server")
	defer println("User", id, "left the server")

	channel.AddClient(conn)

	helpMessage := buildHelpMessage()
	reader := bufio.NewReader(conn)

	for {
		message, err := readMessage(reader)
		if err != nil {
			return
		}

		parts := strings.Fields(message)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "/topics":
			topics := channel.ListTopics()

			var topicList strings.Builder
			if len(topics) == 0 {
				topicList.WriteString("No topics available.\n")
				if _, err := conn.Write([]byte(topicList.String() + " \x00")); err != nil {
					return
				}
				continue
			}

			topicList.WriteString("Available Topics (Key - Participants:\n")
			for _, topic := range topics {
				topicList.WriteString("\t")
				topicList.WriteString(topic.Key)
				topicList.WriteString(" - ")
				topicList.WriteString(strconv.Itoa(topic.Participants))
				topicList.WriteString("\n")
			}

			if _, err := conn.Write([]byte(topicList.String() + " \x00")); err != nil {
				return
			}
		case "/join":
			if len(parts) == 1 {
				if _, err := conn.Write([]byte("You didn't specify a topic to join.\n\x00")); err != nil {
					return
				}
				continue
			}
			topicKey := parts[1]
			channel.AddParticipant(topicKey, conn, id)
		case "/leave":
			if len(parts) == 1 {
				if _, err := conn.Write([]byte("You didn't specify a topic to leave.\n\x00")); err != nil {
					return
				}
				continue
			}
			topicKey := parts[1]
			channel.RemoveParticipant(topicKey, conn, id)
		case "/quit":
			return
		case "/help":
			if _, err := conn.Write([]byte(helpMessage + " \x00")); err != nil {
				return
			}
		case "/publish":
			switch len(parts) {
			case 1:
				if _, err := conn.Write([]byte("You didn't specify a topic to publish to.\n\x00")); err != nil {
					return
				}
				continue
			case 2:
				if _, err := conn.Write([]byte("You didn't specify a message to publish.\n\x00")); err != nil {
					return
				}
				continue
			}

			topicKey := parts[1]

			var broadcastMessage strings.Builder
			broadcastMessage.WriteString("New post added to ")
			broadcastMessage.WriteString(topicKey)
			broadcastMessage.WriteString(" by user ")
			broadcastMessage.WriteString(id)
			broadcastMessage.WriteString(":\n")
			broadcastMessage.WriteString(strings.Join(parts[2:], " "))

			channel.AddParticipant(topicKey, conn, id)
			channel.BroadcastMessage(topicKey, broadcastMessage.String()+" \x00", conn)
		default:
			continue
		}
	}
}

func readMessage(reader *bufio.Reader) (string, error) {
	message, err := reader.ReadString('\x00')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(message), nil
}

func buildHelpMessage() string {
	var helpMessage strings.Builder
	helpMessage.WriteString("Available commands:\n")
	helpMessage.WriteString("/topics - List all available topics\n")
	helpMessage.WriteString("/join <topic> - Join a topic\n")
	helpMessage.WriteString("/leave <topic> - Leave a topic\n")
	helpMessage.WriteString("/publish <topic> <message> - Publish a message to a topic\n")
	helpMessage.WriteString("/quit - Quit the server\n")
	helpMessage.WriteString("/help - Show this message\n")

	return helpMessage.String()
}

func generateId() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
