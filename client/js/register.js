import { registerUser } from "./services.js"
import { redirectIfLoggedIn, saveSession } from "./auth.js"
import { validateUsername, validateEmail } from "./utilities.js"

const redColor = "#FF0054"
const greenColor = "#6BE37F"
const customColor = "#DEDEE0"

redirectIfLoggedIn()

const usernameInput = document.getElementById("username")
const emailInput = document.getElementById("email")
const passwordInput = document.getElementById("password")
const confirmPasswordInput = document.getElementById("confirmPassword")
const formInput = document.getElementById("register_auth")
const errorMessage = document.getElementById("error_message")

usernameInput.addEventListener("input", () => {
    const value = usernameInput.value
    if (value === "") {
        usernameInput.style.borderColor = customColor
        usernameInput.setCustomValidity("")
    }
    else if (validateUsername(value)) {
        usernameInput.style.borderColor = greenColor
        usernameInput.setCustomValidity("")
    }
    else if (!validateUsername(value)) {
        usernameInput.style.borderColor = redColor
        usernameInput.setCustomValidity("Username must be at least 6 characters long, max 20, alphanumeric only except for '_'")
    }
})

emailInput.addEventListener("input", () => {
    const value = emailInput.value

    if (value === "") {
        emailInput.style.borderColor = customColor
        emailInput.setCustomValidity("")
    }
    else if (validateEmail(value)) {
        emailInput.style.borderColor = greenColor
        emailInput.setCustomValidity("")
    }
    else if (!validateEmail(value)) {
        emailInput.style.borderColor = redColor
        emailInput.setCustomValidity("Email Invalid: Follow the pattern 'johndoe@example.com'")
    }
})

passwordInput.addEventListener("input", () => {
    const value = passwordInput.value

    if (value === "") {
        passwordInput.style.borderColor = customColor
        passwordInput.setCustomValidity("")
    }
    else if (value.length >= 6) {
        passwordInput.style.borderColor = greenColor
        passwordInput.setCustomValidity("")
    }
    else if (value.length < 6) {
        passwordInput.style.borderColor = redColor
        passwordInput.setCustomValidity("Password must be at least 6 characters long")
    }
})

confirmPasswordInput.addEventListener("input", () => {
    const confirmPasswordValue = confirmPasswordInput.value
    const passwordValue = passwordInput.value

    if (confirmPasswordValue === "") {
        confirmPasswordInput.style.borderColor = customColor
        confirmPasswordInput.setCustomValidity("")
    }
    else if (confirmPasswordValue === passwordValue) {
        confirmPasswordInput.style.borderColor = greenColor
        confirmPasswordInput.setCustomValidity("")
    }
    else if (confirmPasswordValue !== passwordValue) {
        confirmPasswordInput.style.borderColor = redColor
        confirmPasswordInput.setCustomValidity("Passwords Must Match!")
    }
})

formInput.addEventListener("submit", async (event) => {
    event.preventDefault()
    errorMessage.textContent = ""

    const username = usernameInput.value.trim()
    const email = emailInput.value.trim()
    const password = passwordInput.value

    const response = await registerUser(username, email, password)
    if (response.error) {
        errorMessage.textContent = response.message
        return
    }
    
    saveSession(response.user, response.token)
    redirectIfLoggedIn()
})
