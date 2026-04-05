package main

type room struct {
	// hold all current clients in this room
	clients map[*client]bool

	// join is a channel for all clients wishing to join this room
	join chan *client

	// leave is a channel for all clients wishing to leave the room
	leave chan *client

	// forward is a channel that holds incoming messages that should be forwarded to other clients
	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// each room is a separate thread that should be run indepedently (but as long as the main server is running)
func (r *room) run() {
	for {
		select {
		// adding a user to a channel
		case client := <-r.join:
			r.clients[client] = true

		// removing a user from a channel
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)

			// send a message to all clients in the room
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}

	}
}
