package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"golang.org/x/sync/semaphore"
)

func Server(monitor *Monitor) {
	// listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Server error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started, listening on port 8080")

	for {
		// wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, monitor)
	}
}

// ----- LOCAL FUNCTIONS ------

func handleConnection(conn net.Conn, monitor *Monitor) {
	monitor.wg.Add(1)
	defer monitor.wg.Done()

	defer conn.Close()
	defer fmt.Println("Closing connection...")

	handlerGrant := semaphore.NewWeighted(1)
	handlerGrant.Acquire(monitor.ctx, 1)

	for {
		msgReceived, err := readFromClient(conn)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		msgSplit := strings.Split(msgReceived, "|")
		if len(msgSplit) < 2 {
			fmt.Println("Error: invalid message received", msgReceived)
			continue
		}

		code := msgSplit[0]
		processId := msgSplit[1]

		switch code {
		case "1": // REQUEST
			fmt.Println("Process", processId, "requested critical section")
			monitor.sendRequest(processId, handlerGrant)

			err = handlerGrant.Acquire(monitor.ctx, 1)
			if err != nil { // server closed
				fmt.Println("Server closed")
				return
			}

			err = sendToClient(conn, "GRANT") // TODO: Change message
			if err != nil {
				fmt.Println("Error sending to client:", err)
				continue
			}

			fmt.Println("Process", processId, "entered critical section")

		// 2 -> GRANT

		case "3": // RELEASE
			fmt.Println("Process", processId, "released critical section")
			monitor.releaseSem.Release(1)

		case "9": // exit
			fmt.Println("Process", processId, "exited")
			return

		default:
			fmt.Println("Error: invalid code received", code)

		}
	}

}

func readFromClient(conn net.Conn) (string, error) {
	byteSlice := make([]byte, 20)
	_, err := conn.Read(byteSlice)
	if err != nil {
		return "", err
	}

	return string(byteSlice), err
}

func sendToClient(conn net.Conn, response string) error {
	if len(response) != 5 {
		return fmt.Errorf("error: text must be 20 characters long")
	}

	var buf bytes.Buffer
	buf.WriteString(response)
	_, err := conn.Write(buf.Bytes())

	return err
}
