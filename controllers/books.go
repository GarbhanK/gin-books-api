package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/sirupsen/logrus"

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

	// TODO: properly check connection
	if h.db != nil {
		connectToDatabase = "ok"
	}

	// put stuff here to ping the db
	var status = models.Status{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
		APIStatus: "ok",
		DBStatus:  connectToDatabase,
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}

// GET /books/?table=books
// Get all books
func (h *Handler) GetAllBooks(c *gin.Context) {
	ctx := context.Background()

	// parse out author name in query params
	table, err := utils.GetParams(c, "table")
	fmt.Printf("table: %s\n", table)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No 'title' parameter provided"})
		return
	}

	data, err := h.db.All(ctx, table)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// POST /books
// Create new book
func (h *Handler) CreateBook(c *gin.Context) {
	ctx := context.Background()

	// Validate input
	var newBook models.InsertBookInput
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.db.Insert(ctx, "books", newBook)
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
	fmt.Printf("bookTitle: %s\n", bookTitle)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No 'title' parameter provided"})
		return
	}

	// array of books to return
	bookDocs, err := h.db.Get(ctx, "books", "title", bookTitle)
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

	// array of books to return
	authorBooks, err := h.db.Get(ctx, "books", "author", author)
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

	booksDeleted, err := h.db.Drop(ctx, "books", "Title", title)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": booksDeleted})
}
