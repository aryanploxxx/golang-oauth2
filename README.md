# OAuth Authentication with Gin-Gonic and Gothic

## Overview
This project demonstrates how to implement OAuth-based user authentication using the `Gin-Gonic` web framework, `Goth` for handling multiple providers, and `Gin-contrib/sessions` for managing user sessions. The primary authentication provider in this example is Google.

## Features

- Google OAuth 2.0 authentication
- Session-based user data storage
- Dynamic provider routing
- Protected routes accessible only after login
- User logout functionality

## Setup

### 1. Clone the Repository
```bash
git clone https://github.com/aryanploxxx/golang-oauth2.git
cd golang-outh2
```

### 2. Create a `.env` File
Add the following environment variables to a `.env` file in the root directory:
```env
CLIENT_ID = your-google-client-id
CLIENT_SECRET = your-google-client-secret
CLIENT_CALLBACK_URL = http://localhost:5000/auth/google/callback
SESSION_SECRET = your-session-secret
```

- **CLIENT_ID**: Your Google OAuth client ID.
- **CLIENT_SECRET**: Your Google OAuth client secret.
- **CLIENT_CALLBACK_URL**: The redirect URI set in your Google Cloud project.
- **SESSION_SECRET**: A secure random key for session management. You can generate one using:
  
  ```bash
  openssl rand -base64 32
  ```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run the Application
```bash
go run main.go
```
The server will start on `http://localhost:5000`.

## Routes

### `GET /`
- Displays the home page with a link to initiate authentication via Google.
- HTML file: `templates/index.html`

### `GET /auth/:provider`
- Initiates the OAuth flow for the specified provider (e.g., `google`).

### `GET /auth/:provider/callback`
- Callback URL for the provider to redirect back after authentication.
- Stores user information (e.g., name, email, access token) in the session.
- Redirects to `/success` on successful login.

### `GET /success`
- Displays the user information stored in the session.
- HTML file: `templates/success.html`

### `GET /logout`
- Clears the session and logs the user out.
- Redirects back to the home page.

## File Structure
```plaintext
.
├── main.go            # Main application file
├── templates
│   ├── index.html     # Home page
│   └── success.html   # Success page displaying user information
├── .env               # Environment variables (not tracked in version control)
├── go.mod             # Go module file
└── go.sum             # Dependencies checksum file
```

## Code Highlights

### Authentication Flow
1. **Sign-In**: The `GET /auth/:provider` route initiates the OAuth flow using `gothic.BeginAuthHandler`.
2. **Callback**: The `GET /auth/:provider/callback` route handles the callback, retrieves user information, and stores it in the session.
3. **Protected Routes**: The `/success` route is protected using middleware to ensure only authenticated users can access it.

### Session Management
- Sessions are stored using `gin-contrib/sessions` with cookie-based storage.
- Cookies are configured with:
  - `MaxAge`: Session expiration time (default 30 days).
  - `HttpOnly`: Prevents client-side scripts from accessing cookies.
  - `Secure`: Ensures cookies are only sent over HTTPS (set to `true` in production).

### Environment Configuration
- The application uses the `github.com/joho/godotenv` package to load environment variables from a `.env` file.

## How to Customize

### Add More Providers
1. Import the desired provider from `github.com/markbates/goth/providers`.
2. Add the provider configuration in the `goth.UseProviders` section:
   ```go
   goth.UseProviders(
       google.New(clientID, clientSecret, clientCallbackURL),
       facebook.New(fbClientID, fbClientSecret, fbCallbackURL),
   )
   ```

### Use a Database for Sessions
Replace the `cookie` store with a database-backed session store like Redis or MySQL.

### Secure the Application for Production
1. Set `Secure: true` in cookie options.
2. Use HTTPS for the application.
