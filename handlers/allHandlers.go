package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func Home(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	// template.ParseFiles to parse the HTML file once it's loaded
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(c.Writer, gin.H{})
	// tmpl.Execute to render the HTML file
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func SignInWithProvider(ctx *gin.Context) {

	provider := ctx.Param("provider")
	// Retrieves the provider parameter from the URL path.
	// For example, if the route is /auth/google, ctx.Param("provider") will return "google".

	q := ctx.Request.URL.Query()
	// Extracts the query parameters from the current HTTP request URL as a url.Values object

	q.Add("provider", provider)
	// Adds a new query parameter provider with the value obtained in step 1 (e.g., "google")
	ctx.Request.URL.RawQuery = q.Encode()
	// Encodes the updated query parameters (q) into a raw query string and assigns it back to the request URL.
	// For example, if q contains provider=google, this line updates ctx.Request.URL.RawQuery to ?provider=google

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	// initiates the OAuth authentication flow for the specified provider.
	// gothic.BeginAuthHandler initiates the OAuth2 flow, redirecting the user to the provider's login/authorization page.
	// The provider prompts the user to authenticate and consent to the application accessing their information.

}

func CallbackHandler(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	// After the user consents, the provider redirects the user to /auth/:provider/callback with an authorization code included in the query parameters.
	// The callbackHandler function is triggered and uses gothic.CompleteUserAuth to exchange the authorization code for an access token (and possibly a refresh token).

	fmt.Println("user", user)
	// We are discarding the user data we might have got here

	// Store user data in session as this will persits, storing data in gin context will not persist
	session := sessions.Default(c)
	session.Set("userName", user.Name)
	session.Set("userEmail", user.Email)
	session.Set("accessToken", user.AccessToken)
	session.Set("userPicture", user.AvatarURL)
	session.Set("refreshToken", user.RefreshToken) // Store it in the session
	session.Save()
	// either set the user data in session this way, and get it in the success page

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Essentially, after a successful sign-in, it retrieves authentication credentials including user details and tokens such as access and refresh tokens ( this will vary depending on the provider's implementation). You can then store and further tailor the authentication process according to your application's needs.
	c.Redirect(http.StatusTemporaryRedirect, "/success")
}

func Success(c *gin.Context) {

	session := sessions.Default(c)

	// Retrieve user data from session
	// Create data structure for template
	data := gin.H{
		"Name":         session.Get("userName"),
		"Email":        session.Get("userEmail"),
		"AccessToken":  session.Get("accessToken"),
		"Picture":      session.Get("userPicture"),
		"RefreshToken": session.Get("refreshToken"),
	}

	// Render the success template with data
	c.HTML(http.StatusOK, "success.html", data)

}

// Add this middleware function at the top level
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		accessToken := session.Get("accessToken")
		if accessToken == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Add logout handler
func Logout(c *gin.Context) {
	session := sessions.Default(c)

	// Clear all session data
	session.Clear()
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1, // This will delete the cookie
	})
	session.Save()

	// Clear Gothic session
	gothic.Logout(c.Writer, c.Request)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
