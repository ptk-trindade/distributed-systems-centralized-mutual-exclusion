package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Error: Missing parameters, add parameters 'r' (amount of requests) and 'k' (time to sleep)")
		os.Exit(1)
	}

	r_str := args[0]
	r, err := strconv.Atoi(r_str)

	if err != nil {
		fmt.Println("Parameter 'r' must be a number, got: ", r_str)
		os.Exit(1)
	}

	k_str := args[1]
	k, err := strconv.Atoi(k_str)
	if err != nil {
		fmt.Println("Parameter 'k' must be a number, got: ", k_str)
		os.Exit(1)
	}

	// create a new mutex server
	mutexServer, err := createMutexServer()
	if err != nil {
		fmt.Println("Error creating mutex server:", err)
		return
	}

	pid := os.Getpid()
	for i := 0; i < r; i++ {
		fmt.Println("Request: ", i)

		err := mutexServer.Lock()
		if err != nil {
			fmt.Println("Error locking mutex:", err)
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
		err = mutexServer.Unlock()
		if err != nil {
			fmt.Println("Error unlocking mutex:", err)
			return
		}
	}

	// CLOSE
	mutexServer.Close()
}
