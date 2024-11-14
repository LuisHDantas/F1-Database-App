package handler

import (
	"fmt"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

// @Summary Get total number of races
// @Description Handles the HTTP request to get the total number of races.
// @Tags races
// @Produce json
// @Success 200 {object} map[string]int "count"
// @Failure 500 {object} map[string]string "error"
// @Router /races/total [get]
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
