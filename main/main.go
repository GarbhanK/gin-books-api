package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-api-test/controllers"
	"github.com/garbhank/gin-api-test/models"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {

	r := gin.Default()

	models.ConnectDatabase()
	
	r.GET("/", controllers.Root)
	r.GET("/ping", controllers.Ping)
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	r.GET("/books/author/", controllers.FindAuthor)
	r.GET("/books/title/", controllers.FindBook)
	// r.PATCH("books/:id", controllers.UpdateBook)
	// r.DELETE("/books/:id", controllers.DeleteBook)

	err := r.Run()
	if err != nil {
		return
	}
}
