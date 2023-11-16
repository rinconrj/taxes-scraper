package main

// Import necessary libraries
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/rinconrj/golang-scraper/internal/contaja"
	"github.com/rinconrj/golang-scraper/internal/google"

	"github.com/joho/godotenv"
)

// Initialize global variables
var (
	client *http.Client  // HTTP client to manage our requests
	docs   []contaja.Doc // Variable to store documents fetched from Contaja
)

// Initialize function to setup our cookie-enabled HTTP client
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jar, _ := cookiejar.New(nil) // Create cookie jar
	client = &http.Client{       // Initialize HTTP client with cookie jar
		Jar: jar,
	}
}

// Main function runs when the program starts
func main() {
	http.HandleFunc("/", handleMain)             // Redirect main site endpoint to handleMain() function
	http.HandleFunc("/callback", handleCallback) // Redirect OAuth2 callback endpoint to handleCallback() function
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server on port 8080
}

// Function to handle web requests on main site
func handleMain(w http.ResponseWriter, r *http.Request) {
	cookies, csrfToken := contaja.GetTokens(client)            // Get authentication tokens from Contaja
	logErr := contaja.ContajaLogin(client, csrfToken, cookies) // Log into Contaja
	if logErr != nil {                                         // Check for login errors
		fmt.Println("login error:", logErr)
		return
	}
	docs, err := contaja.GetFiles(client) // Fetch files from Contaja
	if err != nil {                       // Check for file fetching errors
		fmt.Println("get files from portal error:", err)
		return
	}
	google.CreateEventFromDocs(docs, w, r) // Create events from documents fetched from Contaja
}

// Function to handle web requests on Oauth2 callback endpoint
func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()            // Get current context
	oauthConfig := google.GetOauthConfig() // Get Google OAuth2 credentials
	code := r.URL.Query().Get("code")      // Fetch auth code from URL query
	fmt.Println("google code:", code)
	tok, err := oauthConfig.Exchange(ctx, code) // Exchange auth code for token
	if err != nil {                             // Check for OAuth2 code exchange errors
		fmt.Println("error: oauthConfig.Exchange")
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return Internal Server Error if exchange fails
		return
	}
	google.SaveToken(tok)                  // Save received auth token for further use
	google.CreateEventFromDocs(docs, w, r) // Create events from documents previously fetched from Contaja
}
