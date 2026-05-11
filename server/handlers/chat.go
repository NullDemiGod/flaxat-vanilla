package handlers

import (
	"net/http"
	"encoding/json"
	"flaxat/server/models"
	"flaxat/server/middleware"
)

func CreateChat(response http.ResponseWriter, request *http.Request) {
	var bodyStructure struct {
		Member1 int		`json:"member_1"`
		Member2 int 	`json:"member_2"`
	}

	err := json.NewDecoder(request.Body).Decode(&bodyStructure)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid Request Body",
		})
		return
	}

	if bodyStructure.Member1 == 0 || bodyStructure.Member2 == 0 {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "IDs Of Both Users Are Necessary To Create Chat",
		})
		return
	}

	if bodyStructure.Member1 == bodyStructure.Member2 {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "You Can't Start A Conversation With Yourself",
		})
		return
	}

	chat, err := models.CreateChat(bodyStructure.Member1, bodyStructure.Member2)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Chat Between The Given Users Already Exists",
		})
		return
	}

	writeJSON(response, http.StatusCreated, chat)
}

func GetUserChat(response http.ResponseWriter, request *http.Request) {
	userID, ok := middleware.GetUserID(request)
	if !ok || userID == 0 {
		writeJSON(response, http.StatusUnauthorized, map[string]string{
			"error": "Invalid User ID",
		})
		return
	}

	chatList, err := models.GetAllUserChats(userID)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Chat List Is Empty",
		})
		return
	}

	writeJSON(response, http.StatusOK, chatList)
}
