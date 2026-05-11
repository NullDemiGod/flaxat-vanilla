// This file is a display driver, doesn't talk directly to the backend, just handles rendering
// Specifically, rendering the users sidebar

// Select all the elements we need to render the list
const usersSidebar = document.getElementById("users_sidebar")
const usersList = document.getElementById("scrollable_sidebar_users")
const openSidebarButton = document.getElementById("open_users_sidebar")
const closeSidebarButton = document.getElementById("close_users_sidebar")

// This function initializes the button behaviors, and attaches the callback to the list
// Since this is a display driver, it doesn't know what to do if a user element is clicked
// The chat.js file does, and passes that function as a callback to the init function
export function initUsersList(onUserClick) {
    // If we click on the + button, it opens the sidebar
    openSidebarButton.addEventListener("click", () => {
        usersSidebar.classList.remove("hidden")
    })

    // If we click on the X button, it closes the sidebar
    closeSidebarButton.addEventListener("click", () => {
        usersSidebar.classList.add("hidden")
    })

    // Attach the callback to the list, for later use
    usersList._onUserClick = onUserClick
}


// This function handles rendering the users list
export function renderUsersList(users, currentID, onlineUsersID) {
    // Clear the list each time we render
    usersList.innerHTML = ""

    // Iterate through each user
    users.forEach((user) => {
        // No reason to render the current user, since no one can chat with themselves
        if (user.id === currentID)
            return

        // Determine if the user is online or not
        const isOnline = onlineUsersID.includes(user.id)
        const onlineClassRef = isOnline ? "online" : ""

        // Construct the list element and its children
        const newListItem = document.createElement("li")
        newListItem.className = "user_list_item"

        const usernameSpan = document.createElement("span")
        usernameSpan.className = "user_list_name"
        usernameSpan.textContent = user.username

        const onlineIndicatorSpan = document.createElement("span")
        onlineIndicatorSpan.className = `online_indicator ${onlineClassRef}`

        // Stitch the children to the list element
        newListItem.appendChild(usernameSpan)
        newListItem.appendChild(onlineIndicatorSpan)
        
        // Attach behavior to each user item, this the reason we passed the callback to the init()
        newListItem.addEventListener("click", () => {
            usersList._onUserClick(user)
            usersSidebar.classList.add("hidden")
        })
        
        // Append the finished list to the user list
        usersList.appendChild(newListItem)
    })
}
