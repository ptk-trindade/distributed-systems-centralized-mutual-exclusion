package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const REQUEST = 1
const RELEASE = 3
const CLOSE = 9

type mutexServer struct {
	conn net.Conn
	pid  int
}

func createMutexServer() (*mutexServer, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		err := fmt.Errorf("error connecting to server: %s", err)
		return nil, err
	}

	pid := os.Getpid()
	// pid := fmt.Sprintf("%06d", pid_int)

	mutexServer := &mutexServer{
		conn: conn,
		pid:  pid,
	}

	return mutexServer, nil
}

func (mutex *mutexServer) Lock() error {
	// REQUEST
	msg := fmt.Sprintf("%d|%6d|000", REQUEST, mutex.pid)
	err := sendToServer(mutex.conn, msg)
	if err != nil {
		return fmt.Errorf("error sending message to server: %s", err)
	}

	// RECEIVE GRANT
	msg, err = readFromServer(mutex.conn)
	if err != nil {
		return fmt.Errorf("error reading response from server: %s", err)
	}

	if msg != "GRANT" {
		return fmt.Errorf("expected GRANT, got: %s", msg)
	}

	return nil
}

func (mutex *mutexServer) Unlock() error {
	msg := fmt.Sprintf("%d|%6d|000", RELEASE, mutex.pid)
	err := sendToServer(mutex.conn, msg)

	return err
}

func (mutex *mutexServer) Close() error {
	msg := fmt.Sprintf("%d|%6d|000", CLOSE, mutex.pid)
	err := sendToServer(mutex.conn, msg)

	return err
}

// SERVER READ/SEND FUNCTIONS

func sendToServer(conn net.Conn, message string) error {
	if len(message) != 12 {
		return fmt.Errorf("text must be 12 characters long")
	}

	var buf bytes.Buffer
	buf.WriteString(message)
	_, err := conn.Write(buf.Bytes())

	return err
}

func readFromServer(conn net.Conn) (string, error) {
	// read the response from the server
	byteSlice := make([]byte, 5)
	_, err := conn.Read(byteSlice)
	if err != nil {
		return "", err
	}

	return string(byteSlice), nil
}
