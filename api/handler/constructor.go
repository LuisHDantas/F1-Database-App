package handler

import (
	"fmt"
	"net/http"

	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-gonic/gin"
)

// func Constructor_NDrivers(c *gin.Context) {
// 	constructorName := c.Param("name")

// 	// Get the number of drivers
// 	query := "SELECT COUNT(DISTINCT nomepiloto) AS Numero_de_Pilotos FROM QUALIFICA WHERE idconstrutor = (SELECT ID FROM CONSTRUTORES WHERE Nome = $1);"
// 	row := database.DB.QueryRow(query, constructorName)

// 	var nDrivers int
// 	err := row.Scan(&nDrivers)
// 	fmt.Println(nDrivers)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(200, gin.H{"count": nDrivers})
// }

func Constructor_NDrivers(c *gin.Context) {
	// use get_unique_driver_count(constructor_name VARCHAR)
	constructorName := c.Param("name")

	query := "SELECT get_unique_driver_count($1);"
	row := database.DB.QueryRow(query, constructorName)

	var nDrivers int
	err := row.Scan(&nDrivers)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"count": nDrivers})
}

func NConstructors(c *gin.Context) {
	// Get the number of constructors
	query := "SELECT COUNT(DISTINCT nome) AS Numero_de_Construtores FROM construtores;"
	row := database.DB.QueryRow(query)

	var nConstructors int
	err := row.Scan(&nConstructors)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"count": nConstructors})
}

// quantidade de escuderias cadastradas, cada uma com a respectiva quantidade de pilotos;
func Constructors_drivers_count(c *gin.Context) {
	query := `
		SELECT c.Nome, COUNT(DISTINCT q.nomepiloto) AS Numero_de_Pilotos
		FROM CONSTRUTORES c
		LEFT JOIN QUALIFICA q ON c.ID = q.idconstrutor
		GROUP BY c.Nome;
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	constructors := []gin.H{}
	for rows.Next() {
		var name string
		var nDrivers int
		if err := rows.Scan(&name, &nDrivers); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		constructors = append(constructors, gin.H{"constructor": name, "drivers_count": nDrivers})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(constructors)
	c.JSON(200, constructors)
}

func Constructor_victories(c *gin.Context) {
	constructorName := c.Param("name")

	query := "SELECT get_constructor_victories($1);"

	row := database.DB.QueryRow(query, constructorName)

	var victories int
	err := row.Scan(&victories)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"victories": victories})
}

func Constructor_data_range(c *gin.Context) {
	constructorName := c.Param("name")

	query := "SELECT get_constructor_year_range($1);"

	row := database.DB.QueryRow(query, constructorName)

	var dataRange string
	err := row.Scan(&dataRange)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var startYear, endYear string
	_, err = fmt.Sscanf(dataRange, "(%4s,%4s)", &startYear, &endYear)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"start_year": startYear, "end_year": endYear})
}

func Driver_victories_for_constructor(c *gin.Context) {
	constructorName := c.Param("name")

	query := "SELECT * FROM get_constructor_driver_wins($1);"

	rows, err := database.DB.Query(query, constructorName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	driverWins := []gin.H{}
	for rows.Next() {
		var driver string
		var wins int
		if err := rows.Scan(&driver, &wins); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		driverWins = append(driverWins, gin.H{"driver": driver, "wins": wins})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, driverWins)
}

func Status_count_for_constructor(c *gin.Context) {
	constructorName := c.Param("name")

	query := "SELECT * FROM get_constructor_status_count($1);"

	rows, err := database.DB.Query(query, constructorName)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	statusCount := []gin.H{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		statusCount = append(statusCount, gin.H{"status": status, "count": count})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, statusCount)
}

func Constructor_add(c *gin.Context) {
	var constructor struct {
		Nome          string `form:"nome"`
		Nacionalidade string `form:"nacionalidade"`
	}

	if err := c.ShouldBind(&constructor); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO construtores (id, nome, nomepais)
		VALUES ($1, $2, $3)
	`

	var id int
	err := database.DB.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM construtores").Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = database.DB.Exec(query, id, constructor.Nome, constructor.Nacionalidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Constructor added successfully"})
}

func Constructor_driver_search(c *gin.Context) {
	var form struct {
		Surname string `form:"surname"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	constructorName := c.Param("name")
	driverSurname := form.Surname

	query := `
		SELECT p.nome AS nome_completo,
			   p.datanascimento AS data_nascimento,
			   p.nomepais AS nacionalidade
		FROM PILOTOS p
		JOIN RESULTADOS r ON p.nome = r.nomepiloto
		WHERE p.nome LIKE '%' || $1 || '%'
		  AND r.idconstrutor = (SELECT id FROM construtores WHERE nome = $2)
		GROUP BY p.nome, p.datanascimento, p.nomepais;
	`

	rows, err := database.DB.Query(query, driverSurname, constructorName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	drivers := []gin.H{}
	for rows.Next() {
		var name, birthDate, nationality string
		if err := rows.Scan(&name, &birthDate, &nationality); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		drivers = append(drivers, gin.H{"name": name, "birth_date": birthDate, "nationality": nationality})
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, drivers)
}
