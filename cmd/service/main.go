// Public Routes:
// - GET /login: Renders the login HTML page.
// - POST /login: Handles user login via the Login handler.

// Protected Routes:
// - GET /overview: Renders the overview HTML page.
// - GET /report: Renders the report HTML page.
// - GET /driver/create: Renders the add driver HTML page.
// - GET /constructor/create: Renders the add constructor HTML page.
// - GET /driver/search: Renders the driver search HTML page.

// API Routes (Protected):
// - GET /user/original_id: Returns the original ID of the authenticated user.
// - GET /user/role: Returns the role of the authenticated user.
// - GET /constructor/:name/drivers/count: Returns the number of drivers for a specific constructor.
// - GET /drivers/count: Returns the total number of drivers.
// - GET /constructors/drivers/count: Returns the count of drivers for each constructor.
// - GET /constructors/count: Returns the total number of constructors.
// - GET /races/count: Returns the total number of races.
// - GET /circuits/overview: Returns an overview of all circuits.
// - GET /seasons/races/count: Returns the count of races per season.
// - GET /constructor/:name/victories/count: Returns the number of victories for a specific constructor.
// - GET /constructor/:name/data_range: Returns the data range for a specific constructor.
// - GET /driver/:name/data_range: Returns the data range for a specific driver.
// - GET /driver/:name/performances/year: Returns the performance of a specific driver by year.
// - GET /driver/:name/performances/circuit: Returns the performance of a specific driver by circuit.
// - GET /status/count: Returns the count of different statuses.
// - GET /constructor/:name/drivers/victories: Returns the victories of drivers for a specific constructor.
// - GET /constructor/:name/status/count: Returns the count of statuses for a specific constructor.
// - GET /driver/:name/victories/summary: Returns a summary of victories for a specific driver.
// - GET /driver/:name/results/summary: Returns a summary of results for a specific driver.
// - GET /constructor/:name/driver/search: Searches for drivers of a specific constructor.
// - GET /airports/close_to/:city: Returns airports close to a specific city.

// API Routes (Protected, POST):
// - POST /driver/add: Adds a new driver via the Driver_add handler.
// - POST /constructor/add: Adds a new constructor via the Constructor_add handler.
// - POST /logout: Handles user logout via the Logout handler.

// This file sets up the main HTTP server using the Gin framework, connects to the database,
// and defines the public and protected routes for the application.
//
// Functions:
// - main: Initializes the Gin router, connects to the database, sets up session middleware,
//         loads HTML templates, and defines public and protected routes.
// - setupPublicRoutes: Defines the public routes for the application, including HTML and API routes.
// - setupProtectedRoutes: Defines the protected routes for the application, including HTML and API routes.
//                         These routes require authentication middleware.

package main

import (
	"net/http"

	"github.com/LuisHDantas/F1-Database-App/api/handler"
	"github.com/LuisHDantas/F1-Database-App/api/middleware"
	"github.com/LuisHDantas/F1-Database-App/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()

	// Set up session middleware
	store := cookie.NewStore([]byte("super-secret-key"))
	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("server/pages/*.html")

	// Public routes
	setupPublicRoutes(router)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.Auth_middleware())
	setupProtectedRoutes(protected)

	router.Run(":3000")
}

func setupPublicRoutes(router *gin.Engine) {
	// HTML routes
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// API routes
	router.POST("/login", handler.Login)
}

func setupProtectedRoutes(router *gin.RouterGroup) {
	// HTML routes
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

	router.GET("/driver/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "driver_search.html", nil)
	})

	// API routes
	router.GET("/user/original_id", func(c *gin.Context) {
		c.JSON(200, gin.H{"Original_ID": middleware.Original_ID(c)})
	})

	router.GET("/user/role", func(c *gin.Context) {
		c.JSON(200, gin.H{"Role": middleware.User_role(c)})
	})

	router.GET("/constructor/:name/drivers/count", handler.Constructor_NDrivers)
	router.GET("/drivers/count", handler.Total_drivers)
	router.GET("/constructors/drivers/count", handler.Constructors_drivers_count)
	router.GET("/constructors/count", handler.NConstructors)
	router.GET("/races/count", handler.Total_races)
	router.GET("/circuits/overview", handler.Circuits_overview)
	router.GET("/seasons/races/count", handler.Season_races_count)
	router.GET("/constructor/:name/victories/count", handler.Constructor_victories)
	router.GET("/constructor/:name/data_range", handler.Constructor_data_range)
	router.GET("/driver/:name/data_range", handler.Driver_data_range)
	router.GET("/driver/:name/performances/year", handler.Driver_performances_by_year)
	router.GET("/driver/:name/performances/circuit", handler.Driver_performances_by_circuit)
	router.GET("/status/count", handler.Status_count)
	router.GET("/constructor/:name/drivers/victories", handler.Driver_victories_for_constructor)
	router.GET("/constructor/:name/status/count", handler.Status_count_for_constructor)
	router.GET("/driver/:name/victories/summary", handler.Driver_wins_summary)
	router.GET("/driver/:name/results/summary", handler.Driver_results_summary)
	router.GET("/constructor/:name/driver/search", handler.Constructor_driver_search)
	router.GET("/airports/close_to/:city", handler.Airports_close_to)

	router.POST("/driver/add", handler.Driver_add)
	router.POST("/constructor/add", handler.Constructor_add)
	router.POST("/logout", handler.Logout)
}
