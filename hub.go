package main

import (
	"fmt"
)

type message struct {
	data []byte
	room string
}

type subscription struct {
	conn            *connection
	room            string
	connectionType  string
	connectionLimit int64
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	clients map[subscription]bool

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
	clients:    make(map[subscription]bool),
}

func (h *hub) run() {
	for {
		select {

		case s := <-h.register:
			connections := h.rooms[s.room]
			fmt.Println(s.connectionType)
			if len(connections) <= int(s.connectionLimit)-1 {
				if connections == nil {
					if s.connectionType == "host" {
						connections = make(map[*connection]bool)
						h.rooms[s.room] = connections
					} else {
						delete(connections, s.conn)
						close(s.conn.send)
						break
					}
				}
				h.rooms[s.room][s.conn] = true
			} else {
				delete(connections, s.conn)
				close(s.conn.send)
			}

		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
