package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
)

func Server(mutexLog *MutexLog, wg *sync.WaitGroup) {
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

		go handleConnection(conn, mutexLog, wg)
	}
}

// ----- LOCAL FUNCTIONS ------

func handleConnection(conn net.Conn, mutexLog *MutexLog, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	defer conn.Close()
	defer fmt.Println("Closing connection...")

	pid, err := receiveRequest(conn)
	fmt.Println("Received request from", pid)
	for err == nil {
		err = mutexLog.Lock(pid)
		if err != nil {
			sendToClient(conn, "CLOSE")
			return
		}
		sendToClient(conn, "GRANT")

		receiveRelease(conn) // wait for client to release

		mutexLog.Unlock()

		pid, err = receiveRequest(conn)
	}
}

func receiveRequest(conn net.Conn) (string, error) {

	msgTuple, err := readFromClient(conn)
	if err != nil {
		return "", err
	}

	code := msgTuple[0]
	processId := msgTuple[1]

	if code != "1" {
		return "", fmt.Errorf("error: invalid code received %s", code)
	}

	return processId, nil
}

func receiveRelease(conn net.Conn) error {

	msgTuple, err := readFromClient(conn)
	if err != nil {
		return err
	}

	code := msgTuple[0]
	if code != "3" {
		return fmt.Errorf("error: invalid code received %s", code)
	}

	return nil
}

// CLIENT READ/SEND FUNCTIONS

func readFromClient(conn net.Conn) ([2]string, error) {
	var msgTuple [2]string

	byteSlice := make([]byte, 20)
	_, err := conn.Read(byteSlice)
	if err != nil {
		return msgTuple, err
	}

	msg := string(byteSlice)
	fmt.Println("readFromClient:", msg)
	msgSplit := strings.Split(msg, "|")

	if len(msgSplit) < 2 {
		return msgTuple, fmt.Errorf("error: invalid message received %s", msg)
	}

	msgTuple = [2]string{msgSplit[0], msgSplit[1]}

	return msgTuple, err
}

func sendToClient(conn net.Conn, response string) error {
	if len(response) != 5 {
		return fmt.Errorf("error: text must be 5 characters long")
	}

	var buf bytes.Buffer
	buf.WriteString(response)
	_, err := conn.Write(buf.Bytes())

	return err
}
