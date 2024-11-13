package main

import (
	"net/http"

	"github.com/LuisHDantas/F1-Database-App/api/handler"
	"github.com/LuisHDantas/F1-Database-App/api/middleware"
	"github.com/LuisHDantas/F1-Database-App/api/model"
	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: Check later
	// gin.SetMode(gin.ReleaseMode) //optional to not get warning
	// route.SetTrustedProxies([]string{"192.168.1.2"}) //to trust only a specific value

	router := gin.Default()
	database.ConnectDatabase()
	// Set up session middleware
	store := cookie.NewStore([]byte("super-secret-key"))
	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("server/pages/*.html")

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/overview", func(c *gin.Context) {
		c.HTML(http.StatusOK, "overview.html", nil)
	})

	router.GET("/report", func(c *gin.Context) {
		c.HTML(http.StatusOK, "relatorios.html", nil)
	})

	router.GET("/driver/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_driver.html", nil)
	})

	router.GET("/constructor/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_constructor.html", nil)
	})

	router.POST("/login", handler.Login)

	protected := router.Group("/")
	protected.Use(middleware.Auth_middleware())
	{
		protected.GET("/countries", func(c *gin.Context) {
			model.GetCountries()
			c.JSON(200, "OK")
		})

		protected.GET("/user/original_id", func(c *gin.Context) {
			c.JSON(200, gin.H{"Original_ID": middleware.Original_ID(c)})
		})

		protected.GET("/user/role", func(c *gin.Context) {
			c.JSON(200, gin.H{"Role": middleware.User_role(c)})
		})

		protected.GET("/constructor/:name/drivers/count", handler.Constructor_NDrivers)

		protected.GET("/drivers/count", handler.Total_drivers)

		protected.GET("/constructors/drivers/count", handler.Constructors_drivers_count)

		protected.GET("/constructors/count", handler.NConstructors)

		protected.GET("/races/count", handler.Total_races)

		protected.GET("/circuits/overview", handler.Circuits_overview)

		protected.GET("/seasons/races/count", handler.Season_races_count)

		protected.GET("/constructor/:name/victories/count", handler.Constructor_victories)

		protected.GET("/constructor/:name/data_range", handler.Constructor_data_range)

		protected.GET("/driver/:name/data_range", handler.Driver_data_range)

		protected.GET("/driver/:name/performances/year", handler.Driver_performances_by_year)

		protected.GET("/driver/:name/performances/circuit", handler.Driver_performances_by_circuit)

		protected.GET("/status/count", handler.Status_count)

		protected.GET("/constructor/:name/drivers/victories", handler.Driver_victories_for_constructor)

		protected.GET("/constructor/:name/status/count", handler.Status_count_for_constructor)

		protected.GET("/driver/:name/victories/summary", handler.Driver_wins_summary)

		protected.GET("/driver/:name/results/summary", handler.Driver_results_summary)

		protected.POST("/driver/add", handler.Driver_add)

		protected.POST("/constructor/add", handler.Constructor_add)

		protected.POST("/logout", handler.Logout)
	}

	router.Run(":3000")
}
