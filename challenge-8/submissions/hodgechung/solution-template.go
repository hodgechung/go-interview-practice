// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	// Add any other necessary imports
)

// Client represents a connected chat client
type Client struct {
	// TODO: Implement this struct
	// Hint: username, message channel, mutex, disconnected flag
	username string
  disconnected bool
	Incoming chan string
}

// Send sends a message to the client
func (c *Client) Send(message string) error {
	// TODO: Implement this method
	// Hint: thread-safe, non-blocking send
	select {
	case c.Incoming <- message:
	default:
    return ErrClientDisconnected
	}
  return nil
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	return <-c.Incoming
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	clients map[string]*Client
	mu      sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		make(map[string]*Client),
		sync.Mutex{},
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	if len(username) == 0 {
		return nil, ErrInvalidUsername
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.clients[username]
	if ok {
		return nil, ErrUsernameAlreadyTaken
	}

	c := Client{username, false, make(chan string, 1024)}
	s.clients[username] = &c

	return &c, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels
	s.mu.Lock()
	defer s.mu.Unlock()
  client.disconnected = true
	close(client.Incoming)
	delete(s.clients, client.username)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range s.clients {
		if sender.username != c.username {
			c.Send(message)
		}
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	s.mu.Lock()
	defer s.mu.Unlock()

  if sender.disconnected {
    return ErrClientDisconnected
  }

	receiver, ok := s.clients[recipient]
	if !ok {
		return ErrRecipientNotFound
	}
	return receiver.Send(message)
}

// Common errors that can be returned by the Chat Server
var (
	ErrInvalidUsername      = errors.New("invalid username")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
