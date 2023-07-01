package main

import (
	"fmt"
	"sync"
)

func main() {

	mutexLog := CreateMutexLog()
	var wg sync.WaitGroup
	// go Controller(monitor)
	go Server(mutexLog, &wg)

	// wait for a keypress to close the program
	var input string
	for input != "3" {
		fmt.Println("Select an option:")
		fmt.Println("1. Current Requests Queue")
		fmt.Println("2. How many requests per process have been attended")
		fmt.Println("3. Exit")
		fmt.Scanln(&input)

		if input == "3" {
			break
		}

		if input == "1" {
			fmt.Println("Current Requests Queue")
			requestQueue := mutexLog.GetQueue()

			for _, v := range requestQueue {
				fmt.Printf("Process %s\n", v.(string))
			}

		} else if input == "2" {
			fmt.Println("How many requests per process have been attended")
			history := mutexLog.GetHistory()

			attdCount := make(map[string]int)

			for _, v := range history {
				attdCount[v.(string)]++
			}

			for k, v := range attdCount {
				fmt.Printf("Process %s: %d\n", k, v)
			}

		} else {
			fmt.Println("Invalid option")
		}

	}

	mutexLog.Cancel()
	fmt.Println("Ctx cancelled")
	wg.Wait()
}
