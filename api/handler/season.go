package handler

import (
	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

func Season_races_count(c *gin.Context) {
	query := `
	SELECT
		C.Temporada AS season,          -- Season (year or identifier)
		COUNT(C.ID) AS corridas_count   -- Number of races in each season
	FROM
		CORRIDAS C
	GROUP BY
		C.Temporada
	ORDER BY
		C.Temporada;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	seasons := []gin.H{}
	for rows.Next() {
		var season string
		var racesCount int

		if err := rows.Scan(&season, &racesCount); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		seasons = append(seasons, gin.H{
			"season":      season,
			"races_count": racesCount,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, seasons)
}
