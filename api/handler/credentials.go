package handler

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/LuisHDantas/F1-Database-App/api/middleware"
	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func hash_password(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func store_login_log(id int) error {
	query := "INSERT INTO log_table (user_id, login_date, login_time) VALUES ($1, CURRENT_DATE, CURRENT_TIME)"

	_, err := database.DB.Exec(query, id)

	return err
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	hashedPassword := hash_password(password)

	query := "SELECT user_id FROM users WHERE login = $1 AND password = $2"
	rows, err := database.DB.Query(query, username, hashedPassword)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if rows.Next() {
		var id int
		err = rows.Scan(&id)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = store_login_log(id)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		middleware.Start_session(c, id)

		c.JSON(200, gin.H{
			"message":  "Login successful",
			"redirect": "/overview",
		})
	} else {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(200, gin.H{"message": "Logout successful"})
}
