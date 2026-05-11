// This file is a display driver, doesn't talk to the backend, just handles rendering
// Specifically, rendering the chatbox

// Store the callback and the user id for later use
let onSendCallback = null
let currentUserID = null

// All these elements will be needed to handle the chatbox dynamically
// It would've been easier if Angular or React was used
const emptyChatState = document.getElementById("empty_chat")
const activeChatState = document.getElementById("active_chat_container")
const recipientHeader = document.getElementById("recipient_info")
const messageList = document.getElementById("message_history")
const messageForm = document.getElementById("message_form")
const messageInput = document.getElementById("message_input")

// Initializes the chatbox, and stores the callback coming from chat.js
// This display driver doesn't know what to do or how to send a message
// It just knows that a user clicked send, and calls the callback which is responsible for it
export function initChatBox(userID, onSend) {
    onSendCallback = onSend
    currentUserID = userID

    // Listen for clicks on the send button
    messageForm.addEventListener("submit", (event) => {
        event.preventDefault()
        const content = messageInput.value.trim()

        if (content === "")
            return

        // Use the callback to send the message
        onSendCallback(content)

        // Reset the input so the user can send again
        messageInput.value = ""
    })
}

// This function is responsible for rendering the chatbox after a user clicked on a chat
export function openChatBox(recipient, isOnline) {
    // Transition from the empty chat state to the active one
    emptyChatState.classList.add("hidden")
    activeChatState.classList.remove("hidden")

    // Reset the header and store the recipientID so we can use it for updates later
    recipientHeader.innerHTML = ""
    recipientHeader.dataset.recipientID = recipient.id

    // Construct the header element, and its children
    const usernameSpan = document.createElement("span")
    usernameSpan.className = "recipient_name"
    usernameSpan.textContent = recipient.username

    const onlineIndicatorSpan = document.createElement("span")
    onlineIndicatorSpan.className = `online_indicator ${isOnline ? "online" : ""}`

    // Stitch the header together
    recipientHeader.appendChild(usernameSpan)
    recipientHeader.appendChild(onlineIndicatorSpan)
}


// This function is responsible for rendering a single message bubble
export function appendMessage(message, userID) {
    // Determine if the message we are rendering is sent or received
    const messageType = (String(message.sender_id) === String(userID)) ? "sent" : "received"

    // Construct the message bubble element and its children
    const messageBubbleDiv = document.createElement("div")
    messageBubbleDiv.className = `message_bubble ${messageType}`

    const messageContentDiv = document.createElement("div")
    messageContentDiv.className = "message_bubble_message"
    messageContentDiv.textContent = message.content

    const messageTimestampDiv = document.createElement("div")
    messageTimestampDiv.className = "message_bubble_timestamp"
    const time = new Date(message.created_at).toLocaleTimeString([], {
        hour: "2-digit", minute: "2-digit"
    })
    messageTimestampDiv.textContent = time

    // Stitch the bubble together and append it to the list of messages
    messageBubbleDiv.appendChild(messageContentDiv)
    messageBubbleDiv.appendChild(messageTimestampDiv)
    messageList.appendChild(messageBubbleDiv)

    // Scroll the list down automatically
    messageList.scrollTop = messageList.scrollHeight
}

// This function uses the above function to render the whole list of messages
export function renderMessages(messages, userID) {
    messageList.innerHTML = ""

    messages.forEach((message) => {
        appendMessage(message, userID)
    })
}

// This function uses the ID we stashed earlier to update the recipient's status
export function updateRecipientStatus(onlineIDs) {
    const indicator = recipientHeader.querySelector(".online_indicator")
    if (!indicator)
        return

    const recipientID = parseInt(recipientHeader.dataset.recipientID)
    
    const isOnline = onlineIDs.includes(recipientID)
    
    if (isOnline)
        indicator.classList.add("online")
    else
        indicator.classList.remove("online")
}

