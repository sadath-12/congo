package main

import (
	"actor_model/room"
	"actor_model/user"
	"time"
)

func main() {
	room := room.NewRoom() // create a new room actor
	go room.Start()        // start the room actor in a goroutine

	user1 := user.NewUser("Alice", room.Inbox) // create a new user actor named Alice
	go user1.Start()                           // start the user actor in a goroutine
	user1.JoinRoom()                           // Alice joins the room

	user2 := user.NewUser("Bob", room.Inbox) // create a new user actor named Bob
	go user2.Start()                    // start the user actor in a goroutine
	user2.JoinRoom()                    // Bob joins the room

	time.Sleep(time.Second) // wait for a second

	user1.SendMessage("Hello, Bob!") // Alice sends a message to Bob
	user2.SendMessage("Hi, Alice!")  // Bob sends a message to Alice

	time.Sleep(time.Second) // wait for a second

	user1.LeaveRoom() // Alice leaves the room
	user2.LeaveRoom() // Bob leaves the room
}
