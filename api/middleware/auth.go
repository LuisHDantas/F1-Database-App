package middleware

import (
	"fmt"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Start_session(c *gin.Context, id int) error {
	session := sessions.Default(c)
	session.Set("user_id", id)
	err := session.Save()
	return err

}

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

func user_ID(c *gin.Context) int {
	userID, _ := c.Get("user_id")
	return userID.(int)
}
