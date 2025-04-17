package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/garbhank/gin-books-api/database"
	"github.com/garbhank/gin-books-api/models"
	"github.com/garbhank/gin-books-api/utils"
)

type Handler struct {
	db database.Database
}

func NewHandler(dbType database.Database) *Handler {
	return &Handler{db: dbType}
}

// GET /
func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am root"})
}

// GET /ping
// get server status
func (h *Handler) Ping(c *gin.Context) {
	currentTime := time.Now()
	connectToDatabase := "unable to connect!"

	// create client
	db := database.NewFirestore()
	err := db.Conn(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	// TODO: properly check connection
	connectToDatabase = "ok"

	// put stuff here to ping the db
	var status = models.Status{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
		APIStatus: "ok",
		DBStatus:  connectToDatabase,
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}

// GET /books
// Get all books
func (h *Handler) FindBooks(c *gin.Context) {
	ctx := context.Background()

	// create client
	db := database.NewFirestore()
	err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	// array of books to return
	var firestoreBooks []models.Book

	// iterator over books collection in firestore
	iter := db.Client.Collection("books").Documents(ctx)
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

	c.JSON(http.StatusOK, gin.H{"data": firestoreBooks})
}

// POST /books
// Create new book
func (h *Handler) CreateBook(c *gin.Context) {
	ctx := context.Background()

	// create client
	db := database.NewFirestore()
	err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	// Validate input
	var newBook models.InsertBookInput
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := db.Insert(ctx, "books", newBook)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// GET /books/title/
// Find a specific book
func (h *Handler) FindBook(c *gin.Context) {
	ctx := context.Background()

	// parse out author name in query params
	bookTitle, err := utils.GetParams(c, "title")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No 'title' parameter provided"})
		return
	}

	// create client
	db := database.NewFirestore()
	err = db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	// array of books to return
	bookDocs, err := db.Get(ctx, "books", "title", bookTitle)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": bookDocs})
}

func (h *Handler) FindAuthor(c *gin.Context) {
	ctx := context.Background()

	// parse out author name in query params
	author, err := utils.GetParams(c, "name")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No 'name' parameter provided"})
		return
	}

	// create client
	db := database.NewFirestore()
	err = db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	// array of books to return
	authorBooks, err := db.Get(ctx, "books", "Author", author)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": authorBooks})
}

// Delete a book by title
func (h *Handler) DeleteBook(c *gin.Context) {
	ctx := context.Background()

	// parse out author name in query params
	title, err := utils.GetParams(c, "title")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No 'title' parameter provided"})
		return
	}

	// create client
	db := database.NewFirestore()
	err = db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	err = db.Drop(ctx, "books", "Title", title)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
