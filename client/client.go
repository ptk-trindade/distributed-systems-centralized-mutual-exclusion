package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

const REQUEST = 1
const RELEASE = 3
const CLOSE = 9
const FILEPATH = "resultado.txt"

func Client(r, k int) {
	pid := os.Getpid()

	// connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for i := 0; i < r; i++ {
		fmt.Println("Request: ", i)
		// REQUEST
		msg := fmt.Sprintf("%d|%06d|000", REQUEST, pid)
		err = sendToServer(conn, msg)
		if err != nil {
			fmt.Println("Error sending message to server:", err)
			return
		}

		// RECEIVE GRANT
		// read the response from the server
		msg, err = readFromServer(conn)
		fmt.Println("Response from server:", msg)
		if err != nil {
			fmt.Println("Error reading response from server:", err)
			return
		}

		// WRITE TO FILE
		_, err = os.Stat(FILEPATH)
		if os.IsNotExist(err) {
			fmt.Println("creating file")
			createFile(FILEPATH)
		}

		fmt.Println("writing to file")
		currentTime := time.Now()
		txt := fmt.Sprintf("%06d | %s\n", pid, currentTime.Format("15:04:05.000"))
		appendToFile(FILEPATH, txt)

		time.Sleep(time.Duration(k) * time.Second)

		fmt.Println("releasing...")
		// RELEASE
		msg = fmt.Sprintf("%d|%06d|000", RELEASE, pid)
		err = sendToServer(conn, msg)
		if err != nil {
			fmt.Println("Error sending message to server:", err)
			return
		}
	}

	msg := fmt.Sprintf("%d|%06d|000", CLOSE, pid)
	err = sendToServer(conn, msg)

}

// ----- LOCAL FUNCTIONS ------

func sendToServer(conn net.Conn, message string) error {
	if len(message) != 12 {
		return fmt.Errorf("error: text must be 12 characters long")
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

	return string(byteSlice), err
}
