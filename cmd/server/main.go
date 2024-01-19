package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abdulmeLINK/bitcoin-core-ibd-simulation/pkg/server"
)

func main() {
	// Parse port and block file path from command-line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: server <port> <block_file>")
		os.Exit(1)
	}
	port := os.Args[1]
	blockFile := os.Args[2]

	// Open log file
	logFile, err := os.OpenFile(fmt.Sprintf("logs/servers/server_%s.log", port), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Set log output to file
	log.SetOutput(logFile)

	// Create a new server
	s := server.NewServer(port, blockFile)
	if s == nil {
		os.Exit(1)
	}

	// Start the server
	s.Start()

	fmt.Println("Server stopped")
}
