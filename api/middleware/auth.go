// Package middleware provides authentication and session management functionalities
// for the F1-Database-App API.

package middleware

import (
	"fmt"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Initializes a new session for the user with the given ID.
// It sets the "user_id" in the session and saves it.
// Parameters:
// - c: *gin.Context - the Gin context
// - id: int - the user ID to be set in the session
// Returns:
// - error: an error if the session could not be saved

func Start_session(c *gin.Context, id int) error {
	session := sessions.Default(c)
	session.Set("user_id", id)
	err := session.Save()
	return err

}

// Gin middleware function that checks if a user session exists.
// If the session does not exist, it returns a 401 Unauthorized response and aborts the request.
// If the session exists, it sets the "user_id" in the context and proceeds to the next handler.
func Auth_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		fmt.Println(userID)

		if userID == nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// Retrieves the role of the user from the database based on the user ID stored in the context.
// Parameters:
// - c: *gin.Context - the Gin context
// Returns:
// - string: the role of the user
func User_role(c *gin.Context) string {
	var role string
	query := "SELECT tipo FROM users WHERE user_id = $1"
	row := database.DB.QueryRow(query, user_ID(c))
	err := row.Scan(&role)

	if err != nil {
		panic(err)
	}

	return role
}

// Retrieves the original ID of the user from the database based on the user ID stored in the context.
// Parameters:
// - c: *gin.Context - the Gin context
// Returns:
// - string: the original ID of the user
func Original_ID(c *gin.Context) string {
	var originalID string

	query := "SELECT id_original FROM users WHERE user_id = $1"
	row := database.DB.QueryRow(query, user_ID(c))

	err := row.Scan(&originalID)

	if err != nil {
		panic(err)
	}

	return originalID
}

// Helper function that retrieves the user ID from the context.
// Parameters:
// - c: *gin.Context - the Gin context
// Returns:
// - int: the user ID
func user_ID(c *gin.Context) int {
	userID, _ := c.Get("user_id")
	return userID.(int)
}
