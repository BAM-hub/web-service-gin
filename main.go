package main

import (
	dbConnection "example/web-service-gin/db"
	useerRouter "example/web-service-gin/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := dbConnection.Connect()
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	router.GET("/user/list", useerRouter.GetUsers)
	router.GET("/user", useerRouter.GetUserByName)
	router.POST("/user", useerRouter.CreateUser)
	router.Run("localhost:8000")
}
