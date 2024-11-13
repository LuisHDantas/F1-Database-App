package handler

import (
	"fmt"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

func Total_races(c *gin.Context) {
	// Get the number
	query := "SELECT COUNT(DISTINCT id) AS Numero_de_Corridas FROM corridas;"
	row := database.DB.QueryRow(query)

	var nRaces int
	err := row.Scan(&nRaces)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(nRaces)
	c.JSON(200, gin.H{"count": nRaces})
}
