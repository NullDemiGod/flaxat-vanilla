// This acts like our main() function, it imports everything, and stitches behavior with UI
import { getUser, getToken, redirectIfNotLoggedIn, logOut } from "../auth.js"
import { createChat, createMessage, getAllUsers, getChatMessages, getUserChats } from "../services.js"
import { connectSocket, sendSocketMessage } from "./socket.js"
import { initUsersList, renderUsersList } from "./userList.js"
import { initChatList, renderChatList, updateOnlineIndicatorChats, updateChatPreview } from "./chatList.js"
import { appendMessage, initChatBox, openChatBox, renderMessages, updateRecipientStatus } from "./chatBox.js"

redirectIfNotLoggedIn()

const AppState = {
    userID: getUser().id,
    userToken: getToken(),
    currentOpenedChat: null,
    currentRecipient: null,
    allUsers: [],
    allChats: [],
    onlineUsersID: []
}

async function main() {
    connectSocket(AppState.userID, AppState.userToken, handleSocketMessage)

    const [usersResult, chatsResult] = await Promise.all([getAllUsers(), getUserChats()])

    if (usersResult && !usersResult.error) {
        AppState.allUsers = usersResult
    }
    else {
        AppState.allUsers = []
    }

    if (chatsResult && !chatsResult.error) {
        AppState.allChats = chatsResult
    }
    else {
        AppState.allChats = []
    }

    initUsersList(handleUserClick)
    initChatList(handleChatClick)
    initChatBox(AppState.userID, handleSendingMessage)

    renderChatList(AppState.allChats, AppState.userID, AppState.allUsers, AppState.onlineUsersID)
    renderUsersList(AppState.allUsers, AppState.userID, AppState.onlineUsersID)

    const logoutButton = document.getElementById("logout_button")
    logoutButton.addEventListener("click", () => {
    logOut()
    })
}

async function handleUserClick(clickedOnUser) {
    if (!clickedOnUser)
        return
    const result = await createChat(AppState.userID, clickedOnUser.id)

    const chatsResult = await getUserChats()
    if (chatsResult && !chatsResult.error) {
        AppState.allChats = chatsResult
    }
    else {
        AppState.allChats = []
    }
    renderChatList(AppState.allChats, AppState.userID, AppState.allUsers, AppState.onlineUsersID)
}

async function handleChatClick(chat, recipient) {
    if (!chat || !recipient)
        return

    AppState.currentRecipient = recipient
    AppState.currentOpenedChat = chat
    
    openChatBox(recipient, AppState.onlineUsersID.includes(recipient.id))

    const allMessages = await getChatMessages(chat.id)
    if (allMessages && !allMessages.error) {
        renderMessages(allMessages, AppState.userID)
    }
    else {
        renderMessages([], AppState.userID)
    }
}

async function handleSendingMessage(content) {
    if (!content || !AppState.currentOpenedChat || !AppState.currentRecipient)
        return
    
    updateChatPreview(AppState.currentOpenedChat.id, "You: " + content)
    const response = await createMessage(AppState.currentOpenedChat.id, content)
    if (response.error)
        return

    appendMessage({
        sender_id: AppState.userID,
        content: content,
        created_at: new Date().toISOString()
    }, AppState.userID)

    sendSocketMessage({
        type: "new_message",
        chat_id: AppState.currentOpenedChat.id,
        sender_id: AppState.userID,
        recipient_id: AppState.currentRecipient.id,
        content: content
    })

}

async function handleSocketMessage(message) {
    switch (message.type) {
        case "new_message":
            handleIncomingMessage(message);
            break;
        case "online_users":
            handleOnlineUsersUpdate(message.online_users);
            break;
    }
}

async function handleOnlineUsersUpdate(onlineIDs) {
    AppState.onlineUsersID = onlineIDs

    updateOnlineIndicatorChats(AppState.onlineUsersID)
    renderUsersList(AppState.allUsers, AppState.userID, AppState.onlineUsersID)

    if (AppState.currentOpenedChat && AppState.currentRecipient) {
        updateRecipientStatus(AppState.onlineUsersID)
    }
}

async function handleIncomingMessage(message) {
    updateChatPreview(message.chat_id, message.content)
    // First Scenario: The chat is currently open on the screen
    if (AppState.currentOpenedChat && AppState.currentOpenedChat.id === message.chat_id) {
        // Construct a message object to match what appendMessage expects
        const messageObject = {
            sender_id: message.sender_id,
            content: message.content,
            created_at: new Date().toISOString()
        }
        // Append the message directly to the screen
        appendMessage(messageObject, AppState.userID)
    } 
    // Second Scenrio: The chat is in the background or it's a brand new chat
    else {
        // Check if we know this chat exists
        const chatExists = AppState.allChats.some(c => c.id === message.chat_id)

        // If we don't know this chat, 
        // Someone just started a brand new conversation with us.
        if (!chatExists) {
            
            // Re-fetch our chats
            const chatsResult = await getUserChats()
            if (chatsResult && !chatsResult.error) {
                AppState.allChats = chatsResult
            }

            // Do we know who the sender is? If not, update our phonebook.
            const senderExists = AppState.allUsers.some(u => u.id === message.sender_id)
            if (!senderExists) {
                const usersResult = await getAllUsers()
                if (usersResult && !usersResult.error) {
                    AppState.allUsers = usersResult
                    renderUsersList(AppState.allUsers, AppState.userID, AppState.onlineUsersID)
                }
            }

            // Finally, redraw the sidebar so the new chat appears
            renderChatList(AppState.allChats, AppState.userID, AppState.allUsers, AppState.onlineUsersID)
        } else {
            // Background Chat
            console.log(`New background message received for chat ${message.chat_id}`)
        }
    }
}


main()
