package controllers

import (
	"net/http"
	"time"
	"context"
	// "fmt"
	// "log"

	"github.com/garbhank/gin-api-test/models"
	"github.com/gin-gonic/gin"
	"cloud.google.com/go/firestore"
	// "google.golang.org/api/iterator"
)


func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am root"})
}


// GET /ping
// get server status
func Ping(ctx context.Context, client *firestore.Client) func(c *gin.Context) {
	currentTime := time.Now()

	// put stuff here to ping firestore db
	var status = models.Status{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
		APIStatus: "ok",
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": status})
	}

}

// GET /books
// Get all books
func FindBooks(ctx context.Context, client *firestore.Client) func(c *gin.Context) {

	// GORM local db
	var books []models.Book
	models.DB.Find(&books)

	// c.JSON(http.StatusOK, gin.H{"data": books})

	// iter := client.Collection("books-api/books").Documents(ctx)
	// defer iter.Stop() // add to clean up resources

	// for {
	// 	log.Println("starting iterator loop")
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Failed to iterate: %v", err)
	// 	}
	// 	fmt.Println(doc.Data())
	// }

	return func(c *gin.Context) {
		// iter := client.Collection("books-api/books").Documents(ctx)
		// defer iter.Stop() // add to clean up resources
		// for {
		// 	log.Println("starting iterator loop")
		// 	doc, err := iter.Next()
		// 	if err == iterator.Done {
		// 		break
		// 	}
		// 	if err != nil {
		// 		log.Fatalf("Failed to iterate: %v", err)
		// 	}
		// 	fmt.Println(doc.Data())
		// }
		c.JSON(http.StatusOK, gin.H{"data": books})
	}
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {
	// Validate input
	var input models.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.Create(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	var book models.Book

	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate Input
	var input models.UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&book).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
