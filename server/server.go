package server

import (
	c "blockchain-api/client"
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Iserver interface {
	Broadcast(message string, sender *c.Client)
	ListUsers(requestor *c.Client)
}

// Server holds all active c.Clients
type Server struct {
	Clients map[string]*c.Client
	Mutex   sync.RWMutex
}

// NewServer creates a new chat server
func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*c.Client),
	}
}

// broadcast sends a message to all c.Clients
func (s *Server) Broadcast(message string, sender *c.Client) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	for _, Client := range s.Clients {
		if Client != sender { // Don't send to the sender
			_, err := fmt.Fprintf(Client.Conn, "> %s: %s\n", sender.Username, message)
			if err != nil {
				fmt.Println("err broadcasting message to client : ", Client.Username)
			}
		}
	}
}

// listUsers sends the list of connected users to the requesting c.Client
func (s *Server) ListUsers(requestor *c.Client) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, err := fmt.Fprintf(requestor.Conn, "Connected users:\n")
	if err != nil {
		fmt.Println("err connecting to client : ", requestor.Username)
	}
	for username := range s.Clients {
		_, err := requestor.Conn.Write([]byte(fmt.Sprintf("- %s\n", username)))
		if err != nil {
			fmt.Println("err listing all users to client : ", requestor.Username)
		}
	}
}

func (s *Server) RegisterUser(conn net.Conn) (string, error) {
	// defer conn.Close()

	// Get username
	var username string
	fmt.Fprintf(conn, "Enter your username: ")
	scanner := bufio.NewScanner(conn)
	for username == "" {
		if scanner.Scan() {
			username = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error scanning for input during register user")
		}
	}

	// Check if username is taken
	s.Mutex.Lock()
	if _, exists := s.Clients[username]; exists {
		s.Mutex.Unlock()
		conn.Write([]byte("Username already taken\n"))
		// if err != nil {
		// 	return "", fmt.Errorf("err %+v writing to client : %s", err, username)
		// }
		// return "", fmt.Errorf("username already exists")
		return "", fmt.Errorf("user %s already exists",username)
	}

	// Create and store new client
	client := &c.Client{
		Conn:     conn,
		Username: username,
	}
	s.Clients[username] = client
	s.Mutex.Unlock()

	// Announce new user
	s.Broadcast("joined the chat", client)
	return username, nil
}

func (s *Server) RemoveUser(username string) {
	// Cleanup on disconnect
	var client *c.Client
	s.Mutex.Lock()
	client = s.Clients[username]
	delete(s.Clients, username)
	s.Mutex.Unlock()
	s.Broadcast("left the chat", client)
}

// handleClient manages one client's connection
func (s *Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	username, err := s.RegisterUser(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := s.Clients[username]

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}

	// Handle messages
	var retryCount int
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			if retryCount >= 5 {
				fmt.Printf("err reading client [%s] messages : %+v\n", username, err)
				break
			}
			retryCount++
			continue
		}

		message := strings.TrimSpace(scanner.Text())

		// Handle commands
		if message == "/quit" {
			break
		} else if message == "/users" {
			s.ListUsers(client)
		} else if message != "" {
			s.Broadcast(message, client)
		}
	}

	s.RemoveUser(username)

}

// Start runs the server
func (s *Server) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s...\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.HandleClient(conn)
	}
}
