package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

// Total_drivers handles the request to get the total number of distinct drivers.
// @Summary Get total number of drivers
// @Description Retrieves the total number of distinct drivers from the database.
// @Tags drivers
// @Produce json
// @Success 200 {object} map[string]int "count"
// @Failure 500 {object} map[string]string "error"
// @Router /drivers/total [get]
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

// Driver_data_range handles the request to get the range of years the data is available.
// @Summary Get driver's active year range
// @Description Retrieves the range of years (start and end) a driver has data stored inside the database.
// @Tags drivers
// @Produce json// @Param name path string true "Driver Name"00 {object} map[string]int "count"
func Driver_data_range(c *gin.Context) {
	driverName := c.Param("name")

	query := "SELECT get_driver_year_range($1);"

	row := database.DB.QueryRow(query, driverName)

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

// Driver_performances_by_year handles the request to get a driver's list of performance by year.
// @Summary Get driver's performance by year
// @Description Retrieves the performance of a driver for each year from the database, including total points and total victories.
// @Tags drivers
// @Produce json// @Param name path string true "Driver Name"00 {object} map[string]int "count"
func Driver_performances_by_year(c *gin.Context) {
	driverName := c.Param("name")

	query := "SELECT * FROM get_driver_performance_by_year($1);"

	rows, err := database.DB.Query(query, driverName)
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

// Driver_performances_by_circuit handles the request to get driver performances by circuit.
//
// @param c *gin.Context - The Gin context containing the request parameters and context.
//
// This function retrieves the performance data of a driver grouped by circuit from the database.
// It expects a driver name as a URL parameter and queries the database using a set-returning function.
// The results are returned as a JSON array of objects, each containing the circuit name, total points, and total victories.
//
// @response 200 - JSON array of driver performances by circuit.
// @response 500 - JSON object with an error message if there is an issue with the database query or iteration.
func Driver_performances_by_circuit(c *gin.Context) {
	driverName := c.Param("name")

	query := "SELECT circuit_name, total_points, total_victories FROM get_driver_performance_by_circuit($1);"

	rows, err := database.DB.Query(query, driverName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var circuitName string
	var totalPoints float64
	var totalVictories int

	performances := []gin.H{}
	for rows.Next() {
		if err := rows.Scan(&circuitName, &totalPoints, &totalVictories); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		performances = append(performances, gin.H{
			"circuit_name":    circuitName,
			"total_points":    totalPoints,
			"total_victories": totalVictories,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the performances data as JSON
	c.JSON(200, performances)
}

// Driver_wins_summary handles the HTTP request to get the summary of a driver's wins.
// @param c *gin.Context - The Gin context which contains the request and response objects.
//
// This function retrieves the driver's name from the URL parameter, executes a SQL query to get the summary
// of the driver's victories, and returns the results in JSON format. If an error occurs during the query
// execution or while scanning the rows, it returns a JSON error response with a 500 status code.
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

// Driver_results_summary handles the HTTP request to get a summary of driver results by status.
// @param c *gin.Context - The context for the request, which includes parameters and other request data.
//
// This function retrieves the driver's name from the URL parameter, executes a database query to get the
// results summary by status, and returns the results as a JSON response. If an error occurs during the
// database query or while scanning the rows, it returns a 500 status code with the error message.
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

// Driver_add handles the addition of a new driver to the database.
//
// @Summary Add a new driver
// @Description This endpoint allows the addition of a new driver to the database by providing the necessary details.
// @Tags drivers
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param nome formData string true "First Name"
// @Param sobrenome formData string true "Last Name"
// @Param numero formData int true "Driver Number"
// @Param codigo formData string true "Driver Code"
// @Param datanascimento formData string true "Date of Birth"
// @Param nomepais formData string true "Country Name"
// @Success 200 {object} map[string]string "message: Driver added successfully"
// @Failure 400 {object} map[string]string "error: Error message"
// @Failure 500 {object} map[string]string "error: Error message"
// @Router /drivers/add [post]
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
