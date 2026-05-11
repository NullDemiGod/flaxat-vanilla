// The Backend API endpoint responsible for handling WebSocket connections
const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
const BASE_URL = `${protocol}//${window.location.host}/ws`

// We define a global socket object so our chat function can target it
let socket = null

// This function creates a socket object for a specific logged in user
export function connectSocket(userID, token, onMessageReceived) {
    // We send the token as well to the backend to authenticate the user
    socket = new WebSocket(`${BASE_URL}?token=${token}`)
    
    // The connection handshake, our backend expects the first message to be the userID
    socket.addEventListener("open", () => {
        console.log("WebSocket Connection Established")

        socket.send(JSON.stringify({ user_id: userID }))
    })

    // This listener handles what happens when a message is received on the pipe
    // That's why we passed a callback to the function, it will be defined in chat.js
    socket.addEventListener("message", (message) => {
        const data = JSON.parse(message.data)

        onMessageReceived(data)
    })

    // Handles what happens when a connection is gracefully closed
    // Debug statements only, nothing too fancy
    socket.addEventListener("close", () => {
        console.log("WebSocket Connection Closed")
    })

    // Handles what happens when an error occurs
    // Debug statements only, nothing too fancy
    socket.addEventListener("error", (err) => {
        console.error("Network Error: ", err)
    })
}

// This function provides a way for our chat handlers to send a message through the socket
export function sendSocketMessage(message) {
    // Only send the message if the socket is connected and ready
    if (socket && socket.readyState === WebSocket.OPEN)
        socket.send(JSON.stringify(message))
}
