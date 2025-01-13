package main

import (
	"fmt"
	"golang-oauth2/handlers"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// first function to always run
// will load environemnt variables from .env file, so that they are accessible in the application via os.Getenv
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}
}

func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
	sessionKey := os.Getenv("SESSION_SECRET")

	// If SESSION_SECRET is not set, generate a random key
	if sessionKey == "" {
		sessionKey = "your-32-byte-secret-key-here-12345"
		log.Println("Warning: Using default session key. Consider setting SESSION_SECRET environment variable")
	}
	// For production, generate a secure random key. You can generate one using:
	// openssl rand -base64 32

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	goth.UseProviders(
		// The goth.UseProviders method registers the identity providers that your application will support for user authentication. It allows you to specify which providers your application will integrate with, making their authentication mechanisms available to your users.
		google.New(clientID, clientSecret, clientCallbackURL),
		// google.New method configures Goth to specifically use Google as one of the supported identity providers.
	)

	/*
		store := cookie.NewStore([]byte(sessionKey))
		// Configure session settings
		store.Options(sessions.Options{
			Path:     "/",   // Path for which cookie is valid
			MaxAge:   3600,  // Max age in seconds (1 hour)
			HttpOnly: true,  // Prevent JavaScript access
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
		})

		router.Use(sessions.Sessions("mysession", store))
	*/

	// Set up Gothic session store
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true in production
	store := cookie.NewStore([]byte(sessionKey))
	store.Options(sessions.Options{
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
	})
	gothic.Store = store
	// store ensures Gothic can manage its own session state during OAuth flows.

	router.Use(sessions.Sessions("mysession", store))
	// This cookie store is set up for Gin's session middleware
	// Using ginStore ensures that your app's session logic and data are separate from Gothic’s internal session handling.

	router.LoadHTMLGlob("templates/*")
	// this line of code to make sure the template is loaded from the templates directory

	router.GET("/", handlers.Home)
	router.GET("/auth/:provider", handlers.SignInWithProvider)
	router.GET("/auth/:provider/callback", handlers.CallbackHandler)
	// Add middleware to protect the success route
	router.GET("/success", handlers.AuthRequired(), handlers.Success)
	router.GET("/logout", handlers.Logout)

	router.Run(":5000")

}
