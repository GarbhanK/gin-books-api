package main

import (
	// "context"
	"fmt"
	
	"github.com/gin-gonic/gin"

	"github.com/garbhank/gin-api-test/controllers"
	"github.com/garbhank/gin-api-test/models"
	// "github.com/garbhank/gin-api-test/db"
)

func main() {

	r := gin.Default()
	models.ConnectDatabase()
	
	// fsClient := db.CreateFirestoreClient(ctx)
	// defer fsClient.Close()

	fmt.Println("Firestore Client created...")

	r.GET("/", controllers.Root)
	r.GET("/ping", controllers.Ping())
	r.GET("/books", controllers.FindBooks())
	r.POST("/books", controllers.CreateBook)
	r.GET("/books/:id", controllers.FindBook)
	r.PATCH("books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	err := r.Run()
	if err != nil {
		return
	}
}
