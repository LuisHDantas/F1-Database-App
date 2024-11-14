package handler

import (
	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

// Status_count handles the HTTP request to retrieve the status counts from the database.
// It executes the SQL query to fetch the status counts and returns the result as a JSON response.
//
// @Summary Retrieve status counts
// @Description Retrieves the counts of different statuses from the database by executing the admin_report_status_counts function.
// @Tags status
// @Produce json
// @Success 200 {array} struct{Status string "json:\"status\""; Count int "json:\"count\""}
// @Failure 500 {object} gin.H{"error": string}
// @Router /status/count [get]
func Status_count(c *gin.Context) {
	query := "SELECT * FROM admin_report_status_counts();"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var status_count []struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}

	for rows.Next() {
		var sc struct {
			Status string `json:"status"`
			Count  int    `json:"count"`
		}
		err := rows.Scan(&sc.Status, &sc.Count)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		status_count = append(status_count, sc)
	}

	c.JSON(200, status_count)
}
