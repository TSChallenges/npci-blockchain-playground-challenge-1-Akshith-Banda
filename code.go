package main

import (
	"blockchain-api/server"
	"fmt"
)

// // Client represents a chat user
// type Client struct {
// 	conn     net.Conn
// 	username string
// }

// // Server holds all active clients
// type Server struct {
// 	clients map[string]*Client
// 	mutex   sync.RWMutex
// }

// // NewServer creates a new chat server
// func NewServer() *Server {
// 	return &Server{
// 		clients: make(map[string]*Client),
// 	}
// }

// // broadcast sends a message to all clients
// func (s *Server) broadcast(message string, sender *Client) {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()

// 	for _, client := range s.clients {
// 		if client != sender { // Don't send to the sender
// 			fmt.Fprintf(client.conn, "> %s: %s\n", sender.username, message)
// 		}
// 	}
// }

// // listUsers sends the list of connected users to the requesting client
// func (s *Server) listUsers(requestor *Client) {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()

// 	fmt.Fprintf(requestor.conn, "Connected users:\n")
// 	for username := range s.clients {
// 		fmt.Fprintf(requestor.conn, "- %s\n", username)
// 	}
// }

// // handleClient manages one client's connection
// func (s *Server) handleClient(conn net.Conn) {
// 	defer conn.Close()

// 	// Get username
// 	fmt.Fprintf(conn, "Enter your username: ")
// 	scanner := bufio.NewScanner(conn)
// 	if !scanner.Scan() {
// 		return
// 	}
// 	username := strings.TrimSpace(scanner.Text())

// 	// Check if username is taken
// 	s.mutex.Lock()
// 	if _, exists := s.clients[username]; exists {
// 		s.mutex.Unlock()
// 		fmt.Fprintf(conn, "Username already taken\n")
// 		return
// 	}

// 	// Create and store new client
// 	client := &Client{
// 		conn:     conn,
// 		username: username,
// 	}
// 	s.clients[username] = client
// 	s.mutex.Unlock()

// 	// Announce new user
// 	s.broadcast(fmt.Sprintf("joined the chat"), client)

// 	// Handle messages
// 	for scanner.Scan() {
// 		message := strings.TrimSpace(scanner.Text())

// 		// Handle commands
// 		if message == "/quit" {
// 			break
// 		} else if message == "/users" {
// 			s.listUsers(client)
// 		} else if message != "" {
// 			s.broadcast(message, client)
// 		}
// 	}

// 	// Cleanup on disconnect
// 	s.mutex.Lock()
// 	delete(s.clients, username)
// 	s.mutex.Unlock()
// 	s.broadcast(fmt.Sprintf("left the chat"), client)
// }

// // Start runs the server
// func (s *Server) Start(port string) error {
// 	listener, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		return fmt.Errorf("failed to start server: %v", err)
// 	}
// 	defer listener.Close()

// 	fmt.Printf("Server listening on port %s...\n", port)
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Printf("Error accepting connection: %v\n", err)
// 			continue
// 		}
// 		go s.handleClient(conn)
// 	}
// }

func main() {
	server := server.NewServer()
	if err := server.Start("8080"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
