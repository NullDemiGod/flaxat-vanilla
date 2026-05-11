import { loginUser } from "./services.js"
import { redirectIfLoggedIn, saveSession } from "./auth.js"

redirectIfLoggedIn()

const emailInput = document.getElementById("email")
const passwordInput = document.getElementById("password")
const formInput = document.getElementById("login_form")
const errorMessage = document.getElementById("error_message")

formInput.addEventListener("submit", async (event) => {
    errorMessage.textContent = ""
    event.preventDefault()
    const email = emailInput.value.trim()
    const password = passwordInput.value

    const response = await loginUser(email, password)
    if (response.error) {
        errorMessage.textContent = response.message
        return
    }

    saveSession(response.user, response.token)
    redirectIfLoggedIn()
})
