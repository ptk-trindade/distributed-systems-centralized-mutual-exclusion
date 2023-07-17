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

	pid := os.Getpid()

	// create a new mutex server
	mutexServer, err := createMutexServer("localhost:8080", pid)
	if err != nil {
		fmt.Println("Error creating mutex server:", err)
		return
	}

	defer mutexServer.Close()

	for i := 0; i < r; i++ {
		fmt.Println("Request: ", i)

		err = mutexServer.Lock()
		if err != nil {
			fmt.Println("Error locking mutex:", err)
			return
		}

		// WRITE TO FILE
		const filepath = "resultado.txt"
		_, err = os.Stat(filepath)
		if os.IsNotExist(err) {
			fmt.Println("creating file")
			createFile(filepath)
		}

		fmt.Println("writing to file")
		currentTime := time.Now()
		txt := fmt.Sprintf("%06d | %s\n", pid, currentTime.Format("15:04:05.000"))
		appendToFile(filepath, txt)

		time.Sleep(time.Duration(k) * time.Second)

		fmt.Println("releasing...")
		// RELEASE
		err = mutexServer.Unlock()
		if err != nil {
			fmt.Println("Error unlocking mutex:", err)
			return
		}
	}
}
