package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/telegram_clone/websocket_service/genproto/chat_service"
	grpcPkg "gitlab.com/telegram_clone/websocket_service/pkg/grpc_client"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	grpcClient grpcPkg.GrpcClientI
}

func newHub(grpcClient grpcPkg.GrpcClientI) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
		grpcClient: grpcClient,
	}
}

type Message struct {
	UserID   int64  `json:"user_id"`
	ChatType string `json:"chat_type"`
	ChatID   int64  `json:"chat_id"`
	Message  string `json:"message"`
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
			fmt.Println("New client connected", client.userID)
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			fmt.Println("Client disconnected", client.userID)
		case data := <-h.broadcast:
			var message Message
			err := json.Unmarshal(data, &message)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Create message
			_, err = h.grpcClient.MessageService().Create(context.Background(), &chat_service.ChatMessage{
				Message: message.Message,
				ChatId:  message.ChatID,
				UserId:  message.UserID,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Get Chat members
			result, err := h.grpcClient.ChatService().GetChatMembers(context.Background(), &chat_service.GetChatMembersParams{
				Limit:  1000,
				Page:   1,
				ChatId: message.ChatID,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, user := range result.Users {
				if message.UserID == user.Id {
					continue
				}

				client, ok := h.clients[user.Id]
				if ok {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(h.clients, client.userID)
					}
				}
			}
		}
	}
}
