// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Stop running
	done chan bool
}

// CommandChannelFactory creates and returns a running CommandChannel
type CommandChannelFactory func(listID string) CommandChannel

// CommandChannel is an alias for chan TodoCommand
type CommandChannel chan TodoCommand

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
		done:       make(chan bool),
	}
}

func (h *Hub) run(pipelineFactory CommandChannelFactory) {
	stores := make(map[string]CommandChannel)

	for {
		select {
		case <-h.done:
			for _, m := range h.clients {
				for c := range m {
					close(c.Send)
				}
			}
			for _, c := range stores {
				close(c)
			}
			return
		case client := <-h.register:
			if _, ok := h.clients[client.ListID]; !ok {
				h.clients[client.ListID] = make(map[*Client]bool)
			}
			h.clients[client.ListID][client] = true
			if _, ok := stores[client.ListID]; !ok {
				stores[client.ListID] = pipelineFactory(client.ListID)
			}
		case client := <-h.unregister:
			if m, ok := h.clients[client.ListID]; ok {
				delete(m, client)
				h.clients[client.ListID] = m
				close(client.Send)
				if len(m) == 0 {
					if s, ok := stores[client.ListID]; ok {
						close(s)
						delete(stores, client.ListID)
					}
				}
			}
		case message := <-h.broadcast:

			command := TodoCommand{}
			decoder := json.NewDecoder(bytes.NewReader(message))
			err := decoder.Decode(&command)
			if err == nil {
				if store, ok := stores[command.ListID]; ok {
					store <- command
				}

				if clients, ok := h.clients[command.ListID]; ok {
					for client := range clients {
						select {
						case client.Send <- message:
						default:
							close(client.Send)
							delete(clients, client)
							h.clients[command.ListID] = clients
						}
					}
				}
			}
		}
	}
}
