import { getToken } from "./auth.js"

const BASE_URL = ""

function getAuthorizationHeader() {
    return { "Authorization": `Bearer ${getToken()}` }
}

function getContentTypeHeader() {
    return { "Content-Type": "application/json" }
}

export async function registerUser(username, email, password) {
    try {
        const response = await fetch(`${BASE_URL}/api/users/register`, {
            method: "POST",
            headers: { ...getContentTypeHeader() },
            body: JSON.stringify( { username, email, password } )
        })

        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }
        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function loginUser(email, password) {
    try {
        const response = await fetch(`${BASE_URL}/api/users/login`, {
            method: "POST",
            headers: { ...getContentTypeHeader() },
            body: JSON.stringify({ email, password })
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }
        
        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function getUserByID(userID) {
    try {
        const response = await fetch(`${BASE_URL}/api/users/${userID}`, {
            method: "GET",
            headers: { ...getContentTypeHeader(), ...getAuthorizationHeader() }
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }
        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function getAllUsers() {
    try {
        const response = await fetch(`${BASE_URL}/api/users`, {
            method: "GET",
            headers: { ...getContentTypeHeader(), ...getAuthorizationHeader() }
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }

        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function createChat(member_1, member_2) {
    try {
        const response = await fetch(`${BASE_URL}/api/chats`, {
            method: "POST",
            headers: { ...getContentTypeHeader(), ...getAuthorizationHeader() },
            body: JSON.stringify( { member_1, member_2 } )
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }

        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function getUserChats() {
    try {
        const response = await fetch(`${BASE_URL}/api/chats`, {
            method: "GET",
            headers: { ...getContentTypeHeader(), ...getAuthorizationHeader() }
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }

        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function createMessage(chat_id, content) {
    try {
        const response = await fetch(`${BASE_URL}/api/messages`, {
            method: "POST",
            headers: { ...getContentTypeHeader(), ...getAuthorizationHeader() },
            body: JSON.stringify( { chat_id, content } )
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }

        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}

export async function getChatMessages(chatID) {
    try {
        const response = await fetch(`${BASE_URL}/api/messages/${chatID}`, {
            method: "GET",
            headers: { ...getContentTypeHeader(), ...getAuthorizationHeader() }
        })
        const data = await response.json()
        if (!response.ok) {
            return { error: true, message: data.error }
        }

        return data
    }
    catch (err) {
        return { error: true, message: "Network Error" }
    }
}
