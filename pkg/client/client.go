package client

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Block struct {
	ID         int
	Date       time.Time
	Size       int
	Downloaded bool
}

type Client struct {
	Blocks      []Block
	Hosts       []string
	DownloadWin int
	Current     int
	windowStart int
	windowEnd   int
	mu          sync.Mutex
}

func NewClient(blocks []Block, hosts []string, downloadWin int) *Client {
	return &Client{
		Blocks:      blocks,
		Hosts:       hosts,
		DownloadWin: downloadWin,
		windowStart: 0,
		windowEnd:   downloadWin,
	}
}

func (c *Client) StartDownload() {
	var wg sync.WaitGroup

	for _, host := range c.Hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			c.downloadBlocks(host)
		}(host)
	}

	wg.Wait()
}

func (c *Client) downloadBlocks(host string) {
	for {
		id, ok := c.getNextBlockID()
		if !ok {
			break
		}

		err := c.requestBlock(id, host)
		if err != nil {
			log.Printf("Failed to download block %d from host %s, switching host\n", id, host)
			// Switch to the next host
			host = c.getNextHost(host)
			// Retry the same block with the new host
			c.mu.Lock()
			c.Current--
			c.mu.Unlock()
			continue
		}
	}
}

func (c *Client) getNextBlockID() (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Current >= len(c.Blocks) {
		return 0, false
	}

	id := c.Blocks[c.Current].ID
	c.Current++

	return id, true
}

func (c *Client) requestBlock(id int, host string) error {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("Failed to connect to host %s: %v\n", host, err)
		return err
	}
	defer conn.Close()

	// Convert the block ID to a byte slice
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, uint32(id))

	// Send the block ID to the server
	_, err = conn.Write(idBytes)
	if err != nil {
		log.Printf("Failed to send block %d ID to host %s: %v\n", id, host, err)
		return err
	}

	// Handle the response
	c.handleResponse(conn, id)

	return nil
}

func (c *Client) handleResponse(conn net.Conn, id int) {
	// Create a buffer to hold the block data
	var buffer bytes.Buffer

	// Read the block data from the connection
	conn.SetReadDeadline(time.Now().Add(20 * time.Second)) // Set a 5 second timeout
	_, err := io.Copy(&buffer, conn)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("Read operation timed out")
		} else {
			log.Println("Read operation failed:", err)
		}
	}
	if err != nil {
		log.Printf("Failed to read block %d data from host %s: %v\n", id, conn.RemoteAddr(), err)
		return
	}

	// Mark the block as downloaded
	for i := range c.Blocks {
		if c.Blocks[i].ID == id {
			c.Blocks[i].Downloaded = true
			break
		}
	}

	// Log the block data
	//log.Printf("Received block %d data: %v\n", id, buffer.Bytes())

	// Log the size of the incoming data
	log.Printf("block id:%d size: %d bytes\n", id, buffer.Len())

	// Move the window forward if the first block in the window has been downloaded
	c.mu.Lock()
	for c.windowStart < len(c.Blocks) && c.Blocks[c.windowStart].Downloaded && c.windowStart < c.windowEnd {
		c.windowStart++
		if c.windowEnd < len(c.Blocks) {
			c.windowEnd++
		}
	}
	c.mu.Unlock()
}

func (c *Client) getNextHost(currentHost string) string {
	for i, host := range c.Hosts {
		if host == currentHost {
			return c.Hosts[(i+1)%len(c.Hosts)]
		}
	}
	return c.Hosts[0]
}
