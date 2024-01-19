package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/abdulmeLINK/bitcoin-core-ibd-simulation/pkg/client"
)

func main() {
	// Parse block ID range from command-line arguments
	var blockIDs []client.Block
	for _, arg := range os.Args[1:] {
		ids := strings.Split(arg, "-")
		if len(ids) != 2 {
			fmt.Println("Invalid block ID range:", arg)
			os.Exit(1)
		}
		startID, err1 := strconv.Atoi(ids[0])
		endID, err2 := strconv.Atoi(ids[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid block ID range:", arg)
			os.Exit(1)
		}
		for id := startID; id <= endID; id++ {
			blockIDs = append(blockIDs, client.Block{ID: id})
		}
	}

	// Define some hosts for testing
	var hosts []string
	startPort := 5000
	endPort := 5015
	for port := startPort; port <= endPort; port++ {
		hosts = append(hosts, fmt.Sprintf(":%d", port))
	}

	// Create a new client
	c := client.NewClient(blockIDs, hosts, 16)

	// Open log file
	logFile, err := os.OpenFile("logs/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Set log output to file
	log.SetOutput(logFile)

	// Start the download
	c.StartDownload()

	log.Println("Download complete")
}
