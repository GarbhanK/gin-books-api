package main

import (
	"context"
	
	"github.com/gin-gonic/gin"

	"github.com/garbhank/gin-api-test/controllers"
	"github.com/garbhank/gin-api-test/models"
	"github.com/garbhank/gin-api-test/db"
)

func main() {

	ctx := context.Background()
	r := gin.Default()

	models.ConnectDatabase()
	db.CreateFirestoreClient(ctx)

	r.GET("/ping", controllers.Ping)
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	r.GET("/books/:id", controllers.FindBook)
	r.PATCH("books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	err := r.Run()
	if err != nil {
		return
	}
}
