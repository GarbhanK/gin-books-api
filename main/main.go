package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-books-api/controllers"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", controllers.Root)
	r.GET("/ping", controllers.Ping)
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	r.GET("/books/author/", controllers.FindAuthor)
	r.GET("/books/title/", controllers.FindBook)
	r.DELETE("/books/", controllers.DeleteBook)
	// r.PATCH("books/:id", controllers.UpdateBook)
	return r
}

func main() {

	r := setupRouter()

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
