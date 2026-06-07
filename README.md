# Chat In Go (chat-in-go)

A lightweight TCP-based Chat & Pub/Sub application built in Go. This project is a hands-on environment for training and reinforcing skills in **Go concurrency (Goroutines, Channels, sync primitives)** and **Network socket programming (TCP)**.

---

## 🚀 Key Features

* **Concurrency**: Spin up lightweight goroutines for managing concurrent TCP connections, sending/receiving streams asynchronously, and monitoring OS interrupts.
* **TCP Socket Communication**: Pure networking implementation using the standard library's `net` package, utilizing null-byte `\x00` termination framing to handle TCP stream boundaries.
* **Pub/Sub Topic System**: Users can dynamically join topics, list active topics, leave topics, and publish messages to subscribers.
* **Thread-Safe State Management**: Active connections, topics, and subscriber mappings are protected using `sync.Mutex`.
* **Graceful Shutdown**: The server and clients handle OS signals (like `SIGINT` / `Ctrl+C`) to cleanly terminate active TCP connections and allow pending goroutines to finish via `sync.WaitGroup`.

---

## 🛠️ Concurrency & Network Concepts Covered

* **Goroutines (`go` statement)**:
  * Server: Spawned for each incoming TCP connection (`internal.HandleConnection`).
  * Client: Split into two concurrent processes, one for receiving (`internal.RunReceiver`) and one for sending (`internal.RunSender`).
* **WaitGroups (`sync.WaitGroup`)**: Used to coordinate clean shutdowns, ensuring the server/client doesn't exit until all connection handlers have finished their cleanup.
* **Mutexes (`sync.Mutex`)**: Protecting singleton instances, client maps, and topic subscriber lists against race conditions.
* **Channels (`chan`)**: Utilized for handling OS shutdown signals (`os.Signal`) and coordinating background shutdowns.
* **TCP Framing**: Reading and writing strings delimited by `\x00` to correctly separate discrete messages from a continuous TCP stream.

---

## 📂 Project Structure

```text
├── client/
│   ├── internal/
│   │   ├── receiver.go   # Reads messages from TCP socket and prints to stdout
│   │   └── sender.go     # Reads commands/text from stdin and sends via TCP socket
│   └── main.go           # Entrypoint for client, dials the server, spawns receiver/sender
├── server/
│   ├── internal/
│   │   ├── channel.go    # Manages global channels, topics, and active client list
│   │   ├── handler.go    # Handles individual TCP connection flow and parses commands
│   │   └── topic.go      # Represents a chat topic and handles subscriber broadcasts
│   └── main.go           # Entrypoint for server, listens on port 8080
├── go.mod                # Go module definition
└── README.md             # Project documentation
```

---

## 💬 Command Reference

Once connected, clients can use the following commands:

| Command | Description |
| :--- | :--- |
| `/help` | Display list of available commands. |
| `/topics` | List all existing topics and the number of active participants in each. |
| `/join <topic>` | Subscribe to a specific topic and receive broadcasted messages. |
| `/publish <topic> <message>` | Send a message to all subscribers of `<topic>` (implicitly joins you to the topic). |
| `/leave <topic>` | Unsubscribe from `<topic>` to stop receiving messages from it. |
| `/quit` | Terminate the connection and close the client cleanly. |

---

## 🏁 Getting Started

### 1. Start the Server
By default, the server listens on `localhost:8080`.

```bash
go run server/main.go
```

### 2. Connect a Client
Run the client executable by providing the address and port of the server.

```bash
go run client/main.go 127.0.0.1 8080
```

Open multiple terminal windows and run the client command in each to test real-time pub/sub messaging!
