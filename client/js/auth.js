export function saveSession(user, token) {
    localStorage.setItem("user", JSON.stringify(user))
    localStorage.setItem("token", token)
}

export function getToken() {
    return localStorage.getItem("token")
}

export function getUser() {
    const stringifiedUser = localStorage.getItem("user")
    if (!stringifiedUser)
        return null

    return JSON.parse(stringifiedUser)
}

export function isLoggedIn() {
    if (getUser() !== null && getToken() !== null)
        return true
    
    return false
}

export function logOut() {
    localStorage.removeItem("user")
    localStorage.removeItem("token")

    window.location.href = "/login.html"
}

export function redirectIfLoggedIn() {
    if (isLoggedIn())
        window.location.href = "/chat.html"
}

export function redirectIfNotLoggedIn() {
    if (!isLoggedIn())
        window.location.href = "/login.html"
}
