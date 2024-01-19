package server

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	Port      string
	BlockFile string
}

func NewServer(port string, blockFile string) *Server {
	log.Println("Creating a new server")
	return &Server{
		Port:      port,
		BlockFile: blockFile,
	}
}

func (s *Server) Start() {
	log.Println("Starting the server")
	listener, err := net.Listen("tcp", "0.0.0.0:"+s.Port)
	if err != nil {
		log.Println("Failed to start server:", err) // Change fmt.Println to log.Println
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err) // Change fmt.Println to log.Println
			continue
		}

		s.handleConnection(conn)

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Handling connection from %s", conn.RemoteAddr().String())
	for {
		// Read the block ID from the connection
		idBytes := make([]byte, 4)
		_, err := io.ReadFull(conn, idBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Failed to read block ID:", err)
			return
		}

		id := binary.BigEndian.Uint32(idBytes)

		// Look up the block size from the text file
		size, err := s.getBlockSize(int(id))
		if err != nil {
			fmt.Println("Failed to get block size:", err)
			return
		}

		// Generate some random data of the appropriate size
		data := make([]byte, size)
		rand.Read(data)
		fmt.Println(data)
		// Send the block data to the client
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Failed to send block data:", err)
			return
		}

		fmt.Printf("Sent block %d data to client\n", id)
		log.Printf("Sent block %d data to client\n", id)
		break
	}
}

func (s *Server) getBlockSize(id int) (int, error) {
	file, err := os.Open(s.BlockFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 3 {
			continue
		}

		blockID, err1 := strconv.Atoi(parts[0])
		size, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			continue
		}

		if blockID == id {
			return size, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("block %d not found", id)
}
