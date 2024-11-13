package model

import "github.com/LuisHDantas/F1-Database-App/database"

func GetCountries() {
	query := "SELECT * FROM paises"
	rows, err := database.DB.Query(query)

	if err != nil {
		panic(err)
	}

	//Print the result
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
	}

}
