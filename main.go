package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sidx04/go-jwt/controllers"
	"github.com/sidx04/go-jwt/initialisers"
	"github.com/sidx04/go-jwt/middleware"
	"github.com/sidx04/go-jwt/migrations"
)

func init() {
	initialisers.LoadEnvironmentVariables()
	initialisers.ConnectDB()
	migrations.MigrateDatabase()
}

func main() {
	app := gin.Default()

	app.POST("/signup", controllers.UserSignUp)
	app.POST("/login", controllers.LoginUser)
	app.GET("/validate", middleware.RequireAuth, controllers.ValidateToken)

	app.Run()
}
