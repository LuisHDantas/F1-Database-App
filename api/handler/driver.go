package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

func Total_drivers(c *gin.Context) {
	// Get the number of drivers
	query := "SELECT COUNT(DISTINCT nome) AS Numero_de_Pilotos FROM pilotos;"
	row := database.DB.QueryRow(query)

	var nDrivers int
	err := row.Scan(&nDrivers)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"count": nDrivers})

}

func Driver_data_range(c *gin.Context) {
	pilotName := c.Param("name")

	query := "SELECT get_driver_year_range($1);"

	row := database.DB.QueryRow(query, pilotName)

	var yearRange string
	err := row.Scan(&yearRange)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var startYear, endYear string
	_, err = fmt.Sscanf(yearRange, "(%4s,%4s)", &startYear, &endYear)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"start_year": startYear, "end_year": endYear})
}

func Driver_performances_by_year(c *gin.Context) {
	pilotName := c.Param("name")

	query := "SELECT * FROM get_driver_performance_by_year($1);"

	rows, err := database.DB.Query(query, pilotName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var year int
	var totalPoints float64 // Use float64 for NUMERIC type
	var totalVictories int  // Use int for victories (INT type)

	performances := []gin.H{}
	for rows.Next() {
		// Scan the appropriate fields based on the result set
		if err := rows.Scan(&year, &totalPoints, &totalVictories); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Append the result to the performances slice
		performances = append(performances, gin.H{
			"year":            year,
			"total_points":    totalPoints,
			"total_victories": totalVictories,
		})
	}

	// Check if there were any errors during the iteration
	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the performances data as JSON
	c.JSON(200, performances)
}

func Driver_performances_by_circuit(c *gin.Context) {
	pilotName := c.Param("name")

	// Correct query to call the set-returning function for circuit-based grouping
	query := "SELECT circuit_name, total_points, total_victories FROM get_driver_performance_by_circuit($1);"

	rows, err := database.DB.Query(query, pilotName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var circuitName string
	var totalPoints float64 // Use float64 for NUMERIC type
	var totalVictories int  // Use int for victories (INT type)

	performances := []gin.H{}
	for rows.Next() {
		// Scan the appropriate fields based on the result set
		if err := rows.Scan(&circuitName, &totalPoints, &totalVictories); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Append the result to the performances slice
		performances = append(performances, gin.H{
			"circuit_name":    circuitName,
			"total_points":    totalPoints,
			"total_victories": totalVictories,
		})
	}

	// Check if there were any errors during the iteration
	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the performances data as JSON
	c.JSON(200, performances)
}

func Driver_wins_summary(c *gin.Context) {
	driverName := c.Param("name")

	query := "SELECT * FROM get_driver_victories_summary($1);"

	rows, err := database.DB.Query(query, driverName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	driverWins := []gin.H{}
	for rows.Next() {
		var year sql.NullInt32
		var circuit sql.NullString
		var victories int

		if err := rows.Scan(&year, &circuit, &victories); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		yearValue := 0
		if year.Valid {
			yearValue = int(year.Int32)
		}

		circuitName := ""
		if circuit.Valid {
			circuitName = circuit.String
		}

		driverWins = append(driverWins, gin.H{"year": yearValue, "circuit": circuitName, "victories": victories})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, driverWins)
}

func Driver_results_summary(c *gin.Context) {
	driverName := c.Param("name")

	query := "SELECT * FROM get_driver_results_by_status($1);"

	rows, err := database.DB.Query(query, driverName)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	driverResults := []gin.H{}

	for rows.Next() {
		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		driverResults = append(driverResults, gin.H{"status": status, "count": count})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, driverResults)
}

func Driver_add(c *gin.Context) {
	var driver struct {
		Nome           string `form:"nome"`
		Sobrenome      string `form:"sobrenome"`
		Numero         int    `form:"numero"`
		Codigo         string `form:"codigo"`
		DataNascimento string `form:"datanascimento"`
		NomePais       string `form:"nomepais"`
	}

	if err := c.ShouldBind(&driver); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fullName := driver.Nome + " " + driver.Sobrenome

	query := `
		INSERT INTO pilotos (nome, numero, codigo, datanascimento, nomepais)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := database.DB.Exec(query, fullName, driver.Numero, driver.Codigo, driver.DataNascimento, driver.NomePais)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Driver added successfully"})

}
