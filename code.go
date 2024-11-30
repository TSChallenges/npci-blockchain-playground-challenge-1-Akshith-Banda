package main

import (
	"blockchain-api/server"
	"fmt"
)

func main() {
	server := server.NewServer()
	if err := server.Start("8080"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
