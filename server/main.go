package main

import (
	"fmt"
)

type QueueElement struct {
	processId    string
	handlerGrant chan bool
}

func main() {

	monitor := createMonitor()
	go Controller(monitor)

	// wait for a keypress to close the program
	var input string
	for input != "3" {
		fmt.Println("Select an option:")
		fmt.Println("1. Current Requests Queue")
		fmt.Println("2. How many requests per process have been attended")
		fmt.Println("3. Exit")
		fmt.Scanln(&input)

		switch input {
		case "1":
			fmt.Println("Current Requests Queue")
			requestQueue := monitor.getRequestQueue()

			for _, v := range requestQueue {
				fmt.Printf("Process %s\n", v.processId)
			}

		case "2":
			fmt.Println("How many requests per process have been attended")
			attendanceCount := monitor.getAttendanceCount()
			for k, v := range attendanceCount {
				fmt.Printf("Process %s: %d\n", k, v)
			}

		case "3":
			fmt.Println("Exiting...")
			monitor.exit()

		default:
			fmt.Println("Invalid option")
		}
	}

}
