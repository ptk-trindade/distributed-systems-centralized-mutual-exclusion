package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createFile(filename string) {
	err := os.WriteFile(filename, []byte(""), 0644)
	check(err)
}

func appendToFile(filePath, content string) error {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Failed to write to file:", err)
		return err
	}

	return nil
}
