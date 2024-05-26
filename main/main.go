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
	v1 := r.Group("/v1")
	{
		v1.GET("/", controllers.Root)
		v1.GET("/ping", controllers.Ping)
		v1.GET("/books", controllers.FindBooks)
		v1.POST("/books", controllers.CreateBook)
		v1.GET("/books/author/", controllers.FindAuthor)
		v1.GET("/books/title/", controllers.FindBook)
		v1.DELETE("/books/", controllers.DeleteBook)
		// v1.PATCH("books/:id", controllers.UpdateBook)
	}

	return r
}

func main() {

	r := setupRouter()

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
