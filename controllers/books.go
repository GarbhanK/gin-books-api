package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/garbhank/gin-api-test/db"
	"github.com/garbhank/gin-api-test/models"

	"github.com/gin-gonic/gin"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am root"})
}

// GET /ping
// get server status
func Ping() func(c *gin.Context) {
	currentTime := time.Now()
	connectToFirestore := "unable to connect"

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()
	if client != nil {
		connectToFirestore = "ok"
	}

	// put stuff here to ping firestore db
	var status = models.Status{
		Timestamp:       currentTime.Format("2006-01-02 15:04:05"),
		APIStatus:       "ok",
		FirestoreStatus: connectToFirestore,
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": status})
	}

}

// GET /books
// Get all books
func FindBooks() func(c *gin.Context) {

	// GORM local db
	// var books []models.Book
	// models.DB.Find(&books)
	// c.JSON(http.StatusOK, gin.H{"data": books})

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()

	// array of books to return
	var firestoreBooks []models.Book

	// iterator over books collection in firestore
	iter := client.Collection("books").Documents(ctx)
	defer iter.Stop() // add to clean up resources

	// loop until all documents are added to books array
	for {
		var fsBookBuffer models.Book
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate:\n%v", err)
		}
		log.Println(doc.Data())
		if err := doc.DataTo(&fsBookBuffer); err != nil {
			log.Fatalf("can't cast docsnap to Book:\n%v", err)
		}

		// append record to array
		firestoreBooks = append(firestoreBooks, fsBookBuffer)
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": firestoreBooks})
	}
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()


	// Validate input
	var input models.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a DocumentReference
	_, _, err := client.Collection("books").Add(ctx, input)
	if err != nil {
        log.Fatalf("Failed adding document:\n%v", err)
	}

	book := models.Book{Title: input.Title, Author: input.Author}
	// models.DB.Create(&book)

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
