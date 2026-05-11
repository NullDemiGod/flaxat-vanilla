package main

import (
	"log"
	"os"
	"net/http"

	"flaxat/server/db"
	"flaxat/server/handlers"
	"flaxat/server/middleware"
	
	"github.com/joho/godotenv"
)

// This function wraps all the handlers and adds CORS headers to every response
// Without these headers, the browser will refuse to accept responses sent from different ports
func corsMiddleware(mux http.Handler) http.Handler {
	// A HandlerFunc is a function that has a ServeHTTP method attached to it
	// ServeHTTP calls the anonymous function itself
	// Because it does so, we can use this to inject necessary CORS headers
	// And that is exactly why we cast our anonymous function to HandlerFunc
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Make sure that the browser knows that CORS is allowed
		if request.Method == "OPTIONS" {
			response.WriteHeader(http.StatusOK)
			return
		}
		
		// After injecting the headers, we call the actual ServeHTTP method on the routing table
		// This way, we set CORS headers, and the routing table calls the needed handler with the modified response
		mux.ServeHTTP(response, request)
	})
}

func main() {
	// Read and Load the .env file and attach the read variables to ENV of this process
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment file")
	}

	// Initiate the database connection pool using the custom Connect() function
	db.Connect()

	// Create the routing table, can be thought of as a dictionary:
	// The key would be the api endpoint, the value is its specified handler function
	routingTable := http.NewServeMux()

	// We implement our routing table and start adding key:value pairs to it

	// User routes
	routingTable.HandleFunc("POST /api/users/register", handlers.Register)
	routingTable.HandleFunc("POST /api/users/login", handlers.Login)
	routingTable.HandleFunc("GET /api/users", middleware.RequireAuth(handlers.GetAllUsers))
	routingTable.HandleFunc("GET /api/users/{userID}", middleware.RequireAuth(handlers.GetUser))

	// Chat routes
	routingTable.HandleFunc("POST /api/chats", middleware.RequireAuth(handlers.CreateChat))
	routingTable.HandleFunc("GET /api/chats/", middleware.RequireAuth(handlers.GetUserChat))

	// Message routes
	routingTable.HandleFunc("POST /api/messages", middleware.RequireAuth(handlers.CreateMessage))
	routingTable.HandleFunc("GET /api/messages/{chatID}", middleware.RequireAuth(handlers.GetChatMessages))

	// WebSocket route
	routingTable.HandleFunc("GET /ws", middleware.RequireAuthWS(handlers.WebSocketHandler))

	// Use corsMiddleware() to wrap the routing table
	handlersWrapped := corsMiddleware(routingTable)

	// Retrieve the PORT from the env variables
	port := os.Getenv("PORT")
	log.Println("Server is listening on port ", port, "...")

	// Now that everything is ready, we configure the server to listen on PORT
	// Then map each connection to its specified handler
	err = http.ListenAndServe(":" + port, handlersWrapped)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
