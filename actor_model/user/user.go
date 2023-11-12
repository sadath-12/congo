package user

import (
	"fmt"
)

// User is an actor that represents a chat user
type User struct {
	Name  string       // the name of the user
	room  chan Message // the channel to communicate with the room
	inbox chan Message // the channel to receive messages from other actors
}

// Message is a struct that represents a message sent between actors
type Message struct {
	Sender      string           // the name of the sender
	Text        string           // the text of the message
	SenderInbox chan Message // the inbox channel of the sender
}

// NewUser creates a new user actor with a given name and room
func NewUser(name string, room chan Message) *User {
	return &User{
		Name:  name,
		room:  room,
		inbox: make(chan Message, 10), // create a buffered channel for the inbox
	}
}

// Start starts the user actor by listening on its inbox channel
func (u *User) Start() {
	for {
		select {
		case msg := <-u.inbox: // receive a message from the inbox
			u.handleMessage(msg) // handle the message
		}
	}
}

// JoinRoom sends a message to the room to join it
func (u *User) JoinRoom() {
	u.room <- Message{u.Name, "join", u.inbox} // send a join message to the room
}

// LeaveRoom sends a message to the room to leave it
func (u *User) LeaveRoom() {
	u.room <- Message{u.Name, "leave", u.inbox} // send a leave message to the room
}

// SendMessage sends a message to the room with a given text
func (u *User) SendMessage(text string) {
	u.room <- Message{u.Name, text, u.inbox} // send a text message to the room
}

// handleMessage handles a message received from the inbox
func (u *User) handleMessage(msg Message) {
	switch msg.Text {
	case "join": // if the message is a join message
		fmt.Printf("%s joined the room\n", msg.Sender) // print that the sender joined the room
	case "leave": // if the message is a leave message
		fmt.Printf("%s left the room\n", msg.Sender) // print that the sender left the room
	default: // otherwise
		fmt.Printf("%s: %s\n", msg.Sender, msg.Text) // print the sender and the text of the message
	}
}
