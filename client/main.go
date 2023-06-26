package main

import (
	"fmt"
	"os"
	"strconv"
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

	Client(r, k)

}
