package handler

import (
	"fmt"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

func Circuits_overview(c *gin.Context) {
	query := `
	SELECT
		CI.Nome AS circuito_nome,             -- Circuit name
		COUNT(DISTINCT (C.ID, C.nomecircuito)) AS corridas_por_circuito, -- Number of races per circuit
		COUNT(DISTINCT (V.idcorrida, V.nomecircuito, V.nomepiloto, V.nrvolta)) AS voltas_total,          -- Total laps for each circuit
		MIN(V.nrvolta) AS voltas_min,    -- Minimum lap count per race in the circuit
		AVG(V.nrvolta) AS voltas_media,  -- Average lap count per race in the circuit
		MAX(V.nrvolta) AS voltas_max     -- Maximum lap count per race in the circuit
	FROM
		CORRIDAS C
	JOIN
		CIRCUITOS CI ON C.nomecircuito = CI.nome  -- Join CIRCUITOS to get circuit details
	LEFT JOIN
		VOLTAS V ON C.ID = V.IDCorrida            -- Join VOLTAS to get lap details
	WHERE
		V.nrvolta IS NOT NULL                     -- Ignore null values in lap count
	GROUP BY
		CI.Nome;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Return data in JSON
	defer rows.Close()

	circuits := []gin.H{}
	for rows.Next() {
		var circuito_nome string
		var corridas_por_circuito int
		var voltas_total int
		var voltas_min int
		var voltas_media float64
		var voltas_max int

		err := rows.Scan(&circuito_nome, &corridas_por_circuito, &voltas_total, &voltas_min, &voltas_media, &voltas_max)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		circuit := gin.H{
			"circuito_nome":         circuito_nome,
			"corridas_por_circuito": corridas_por_circuito,
			"voltas_total":          voltas_total,
			"voltas_min":            voltas_min,
			"voltas_media":          voltas_media,
			"voltas_max":            voltas_max,
		}
		circuits = append(circuits, circuit)
	}

	if err = rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(circuits)
	c.JSON(200, circuits)
}
