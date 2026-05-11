package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var connUpgrader = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

type Client struct {
	ClientID int
	connectionPipe *websocket.Conn
}

type Message struct {
	Type string			`json:"type"`
	SenderID int		`json:"sender_id"`
	RecipientID int		`json:"recipient_id"`
	ChatID int			`json:"chat_id"`
	Content string		`json:"content"`
}

var onlineClients = make(map[int]*Client)
var mutex sync.Mutex

func getOnlineClientsIDs() []int {
	mutex.Lock()
	defer mutex.Unlock()

	var idList = make([]int, 0, len(onlineClients))
	for id := range onlineClients {
		idList = append(idList, id)
	}
	return idList
}

func broadcastOnlineClientsIDs() {
	idList := getOnlineClientsIDs()

	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range onlineClients {
		err := client.connectionPipe.WriteJSON(map[string]any{
			"type": "online_users",
			"online_users": idList,
		})
		if err != nil {
			log.Printf("Failed to broadcast online users list to client %d\n", client.ClientID)
		}
	}
}

func WebSocketHandler(response http.ResponseWriter, request *http.Request) {
	connection, err := connUpgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println("Failed To Upgrade Connection")
		return
	}

	var registration struct {
		UserID int		`json:"user_id"`
	}
	err = connection.ReadJSON(&registration)
	if err != nil || registration.UserID == 0 {
		log.Println("Failed To Fetch Registration ID")
		return
	}

	newClient := &Client{
		ClientID: registration.UserID,
		connectionPipe: connection,
	}
	
	mutex.Lock()
	onlineClients[registration.UserID] = newClient
	mutex.Unlock()
	log.Printf("Client %d is successfully connected\n", registration.UserID)
	broadcastOnlineClientsIDs()

	defer func() {
		mutex.Lock()
		delete(onlineClients, registration.UserID)
		mutex.Unlock()

		connection.Close()
		log.Printf("Client %d is successfully disconnected\n", registration.UserID)
		broadcastOnlineClientsIDs()
	}()

	for {
		var currentMsg Message
		err := connection.ReadJSON(&currentMsg)
		if err != nil {
			log.Printf("Client %d is about to disconnect...\n", registration.UserID)
			break;
		}
		
		mutex.Lock()
		recipientClient, isOnline := onlineClients[currentMsg.RecipientID]
		mutex.Unlock()

		if isOnline {
			err := recipientClient.connectionPipe.WriteJSON(map[string]any{
				"type": "new_message",
				"sender_id": currentMsg.SenderID,
				"chat_id": currentMsg.ChatID,
				"content": currentMsg.Content,
			})
			if err != nil {
				log.Printf("Failed To Forward Message To Client %d\n", currentMsg.RecipientID)
			}

			err = recipientClient.connectionPipe.WriteJSON(map[string]any{
				"type": "notification",
				"sender_id": currentMsg.SenderID,
				"chat_id": currentMsg.ChatID,
			})
			if err != nil {
				log.Printf("Failed To Forward Notification To Client %d\n", currentMsg.RecipientID)
			}
		}
	}
}
