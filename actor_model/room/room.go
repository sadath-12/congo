package room

import (
	"actor_model/user"
)

// Room is an actor that represents a chat room
type Room struct {
	users map[string]chan user.Message // a map of user names to their Inbox channels
	Inbox chan user.Message            // the channel to receive messages from other actors
}

// NewRoom creates a new room actor
func NewRoom() *Room {
	return &Room{
		users: make(map[string]chan user.Message), // create an empty map of users
		Inbox: make(chan user.Message),            // create an unbuffered channel for the Inbox
	}
}

// Start starts the room actor by listening on its Inbox channel
func (r *Room) Start() {
	for {
		select {
		case msg := <-r.Inbox: // receive a message from the Inbox
			r.handleMessage(msg) // handle the message
		}
	}
}

// AddUser adds a user to the room with a given name and Inbox channel
func (r *Room) AddUser(name string, Inbox chan user.Message) {
	r.users[name] = Inbox // add the user to the map of users
}

// RemoveUser removes a user from the room with a given name
func (r *Room) RemoveUser(name string) {
	delete(r.users, name) // delete the user from the map of users
}

// Broadcast sends a message to all the users in the room
func (r *Room) Broadcast(msg user.Message) {
	for _, Inbox := range r.users { // iterate over the user Inbox channels
		Inbox <- msg // send the message to each Inbox
	}
}

// handleMessage handles a message received from the Inbox
func (r *Room) handleMessage(msg user.Message) {
	switch msg.Text {
	case "join": // if the message is a join message
		r.AddUser(msg.Sender, msg.SenderInbox) // add the sender to the room
		r.Broadcast(msg)                       // broadcast the join message to all the users
	case "leave": // if the message is a leave message
		r.RemoveUser(msg.Sender) // remove the sender from the room
		r.Broadcast(msg)         // broadcast the leave message to all the users
	default: // otherwise
		r.Broadcast(msg) // broadcast the text message to all the users
	}
}
