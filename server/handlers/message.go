package handlers

import (
	"encoding/json"
	"flaxat/server/middleware"
	"flaxat/server/models"
	"net/http"
	"strconv"
)

func CreateMessage(response http.ResponseWriter, request *http.Request) {
	senderID, ok := middleware.GetUserID(request)
	if !ok || senderID == 0 {
		writeJSON(response, http.StatusUnauthorized, map[string]string{
			"error": "Can't Fetch User ID",
		})
		return
	}

	var bodyStructure struct {
		ChatID int			`json:"chat_id"`
		Content string		`json:"content"`
	}

	err := json.NewDecoder(request.Body).Decode(&bodyStructure)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid Request Body",
		})
		return
	}

	if bodyStructure.ChatID == 0 || bodyStructure.Content == "" {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Missing SenderID, ChatID or Content",
		})
		return
	}

	chat, err := models.GetChatByID(bodyStructure.ChatID)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Chat Not Found",
		})
		return
	}

	if senderID != chat.Member1 && senderID != chat.Member2 {
		writeJSON(response, http.StatusForbidden, map[string]string{
			"error": "Forbidden",
		})
		return
	}

	newMessage, err := models.CreateMessage(senderID, bodyStructure.ChatID, bodyStructure.Content)
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, map[string]string{
			"error": "Something Went Wrong",
		})
		return
	}

	writeJSON(response, http.StatusCreated, newMessage)
}

func GetChatMessages(response http.ResponseWriter, request *http.Request) {
	userID, ok := middleware.GetUserID(request)
	if !ok || userID == 0 {
		writeJSON(response, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	stringChatID := request.PathValue("chatID")
	
	chatID, err := strconv.Atoi(stringChatID)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid ChatID",
		})
		return
	}

	chat, err := models.GetChatByID(chatID)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Chat Is Not Found",
		})
		return
	}

	if userID != chat.Member1 && userID != chat.Member2 {
		writeJSON(response, http.StatusForbidden, map[string]string{
			"error": "Forbidden",
		})
		return
	}

	messageList, err := models.GetChatMessages(chatID)
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, map[string]string{
			"error": "Failed To Fetch Chat Messages",
		})
		return
	}

	writeJSON(response, http.StatusOK, messageList)
}
