// This file is a display driver, doesn't talk to the backend, just handles rendering
// Specifically, rendering the chats sidebar

// Get a reference for the scrollable HTML element responsible for the chat list 
const chatsList = document.getElementById("scrollable_sidebar_chats")

// This function handles attaching the callback to the HTML element we imported earlier
// Because this is a display driver, it doesn't know what should happen if a chat is clicked on
// The behavior is delegated to chat.js, where all the callbacks/handlers are defined,
// and subsequently passed to the init functions of each display driver
export function initChatList(onChatClick) {
    chatsList._onChatClick = onChatClick
}

// This function handles rendering the chat list, and all its inner components
export function renderChatList(userChats, currentUserID, allUsers, onlineUsersID) {
    // Each time we render the list, empty it first
    chatsList.innerHTML = ""

    // Iterate through each chat
    userChats.forEach((chat) => {
        // If for some reason, a chat exists between a user and themselves, return
        if (chat.member_1 === chat.member_2)
            return

        // Determine who the recipient is
        const recipientID = (currentUserID === chat.member_1) ? chat.member_2 : chat.member_1

        // Determine if the recipient is online or not
        const isOnline = onlineUsersID.includes(recipientID)
        const onlineClassRef = isOnline ? "online" : ""
        
        // Get the recipient object from the users list using the ID we extracted earlier
        const recipientObject = allUsers.find(user => user.id === recipientID)
        if (!recipientObject)
            return

        // Create a list item for the chat and its subelements
        const newChatItem = document.createElement("li")
        newChatItem.className = "chat_list_item"
        // Store the recipientID so we can use it later for updating online status
        newChatItem.dataset.recipientID = recipientObject.id

        const bottomDiv = document.createElement("div")
        bottomDiv.className = "chat_item_bottom"

        const usernameSpan = document.createElement("span")
        usernameSpan.className = "chat_item_name"
        usernameSpan.textContent = recipientObject.username

        const previewSpan = document.createElement("span")
        previewSpan.className = "chat_item_preview"
        newChatItem.dataset.chatID = chat.id

        let previewText = chat.last_message ? chat.last_message : ""

        if (chat.last_message && String(chat.last_message_sender_id) === String(currentUserID)) {
            previewText = "You: " + previewText
        }

        previewSpan.textContent = previewText
        
        const onlineIndicatorSpan = document.createElement("span")
        onlineIndicatorSpan.className = `online_indicator ${onlineClassRef}`

        // Stitch the element together
        bottomDiv.appendChild(previewSpan)
        bottomDiv.appendChild(onlineIndicatorSpan)
        newChatItem.appendChild(usernameSpan)
        newChatItem.appendChild(bottomDiv)

        // Attach the callback to each chat item
        // And also make the chat highlighted once clicked
        newChatItem.addEventListener("click", () => {
            document.querySelectorAll(".chat_list_item").forEach((item) => {
                item.classList.remove("active")
            })
            newChatItem.classList.add("active")
            chatsList._onChatClick(chat, recipientObject)
        })

        // Append the finished element to the list
        chatsList.appendChild(newChatItem)
    })
}

// This function updates the online status of each chat element
export function updateOnlineIndicatorChats(onlineUsersID) {
    // Select all the chat items in the list
    document.querySelectorAll(".chat_list_item").forEach((item) => {
        // Extract the recipientID
        const recipientID = parseInt(item.dataset.recipientID)
        if (!recipientID)
            return

        // Check if the user is online or not
        const isOnline = onlineUsersID.includes(recipientID)

        // Select the dot so we can color it
        const indicator = item.querySelector(".online_indicator")
        if (!indicator)
            return

        // Color it based on the online status
        if (isOnline)
            indicator.classList.add("online")
        else
            indicator.classList.remove("online")
    })
}

export function updateChatPreview(chatID, text) {
    document.querySelectorAll(".chat_list_item").forEach((item) => {
        const itemChatID = parseInt(item.dataset.chatID)
        if (itemChatID === chatID) {
            const preview = item.querySelector(".chat_item_preview")
            if (preview) {
                preview.textContent = text
            }
        }
    })
}
