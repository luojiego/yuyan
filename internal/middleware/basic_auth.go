package middleware

import "github.com/gin-gonic/gin"

// Accounts stores the user credentials
// You should move these credentials to a configuration file in production
var Accounts = gin.Accounts{
	"admin": "admin", // Change this to a secure password
}

// BasicAuth returns the basic auth middleware
func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(Accounts)
}
