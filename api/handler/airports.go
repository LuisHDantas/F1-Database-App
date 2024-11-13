package handler

import (
	"database/sql"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

func Airports_close_to(c *gin.Context) {
	query := `
	SELECT 
		c.Nome AS Nome_Cidade,
		a.iata AS Codigo_IATA,
		a.Nome AS Nome_Aeroporto,
		a.nomepais AS Cidade_Aeroporto,
		(6371 * acos(cos(radians(c.latitude)) * cos(radians(a.Latitude)) * cos(radians(a.Longitude) - radians(c.longitude)) 
		+ sin(radians(c.latitude)) * sin(radians(a.Latitude)))) AS Distancia_km
	FROM 
		CIDADES c
	JOIN 
		AEROPORTOS a ON a.nomepais = 'BR'
	WHERE 
		c.Nome = $1
		AND (6371 * acos(cos(radians(c.latitude)) * cos(radians(a.Latitude)) * cos(radians(a.Longitude) - radians(c.longitude)) 
		+ sin(radians(c.latitude)) * sin(radians(a.Latitude)))) <= 100
	ORDER BY
		Distancia_km ASC;
	`

	city := c.Param("city")
	rows, err := database.DB.Query(query, city)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		var nomeCidade, codigoIATA, nomeAeroporto, cidadeAeroporto sql.NullString
		var distanciaKm float64
		if err := rows.Scan(&nomeCidade, &codigoIATA, &nomeAeroporto, &cidadeAeroporto, &distanciaKm); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		results = append(results, gin.H{
			"Nome_Cidade":      nomeCidade.String,
			"Codigo_IATA":      codigoIATA.String,
			"Nome_Aeroporto":   nomeAeroporto.String,
			"Cidade_Aeroporto": cidadeAeroporto.String,
			"Distancia_km":     distanciaKm,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, results)
}
