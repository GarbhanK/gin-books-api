package main

import (
	"github.com/gin-gonic/gin"

	"github.com/garbhank/gin-api-test/controllers"
	"github.com/garbhank/gin-api-test/models"
)

func main() {

	r := gin.Default()

	models.ConnectDatabase()

	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)

	r.Run()
}
