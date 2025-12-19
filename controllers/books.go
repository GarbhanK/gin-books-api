package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-books-api/database"
	"github.com/garbhank/gin-books-api/models"
	"github.com/garbhank/gin-books-api/utils"
)

type Handler struct {
	primaryDB   database.Database
	secondaryDB database.Database
}

func NewHandler(primary database.Database, secondary database.Database) *Handler {
	err := primary.Setup(context.Background())
	if err != nil {
		log.Errorf("Failed to setup %s database: %v", primary.Type(), err)
	}
	if secondary != nil {
		err := secondary.Setup(context.Background())
		if err != nil {
			log.Errorf("Failed to setup %s database: %v", secondary.Type(), err)
		}
	}

	return &Handler{
		primaryDB:   primary,
		secondaryDB: secondary,
	}
}

// GET /
func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am root"})
}

// GET /ping
// get server status
func (h *Handler) Ping(c *gin.Context) {
	currentTime := time.Now()
	connectionSuccess := "ok"
	connectionFailed := "unable to connect!"
	dbStatus := []models.DBStatus{}

	// check the connection status of the primary db
	var primaryConn string
	if h.primaryDB.IsConnected(context.Background()) {
		primaryConn = connectionSuccess
	} else {
		primaryConn = connectionFailed
	}

	primaryStatus := models.DBStatus{
		Tier:       "primary",
		Type:       h.primaryDB.Type(),
		Connection: primaryConn,
	}
	dbStatus = append(dbStatus, primaryStatus)

	// if set, also check secondary DB connection
	if h.secondaryDB != nil {
		var secondaryConn string
		if h.secondaryDB.IsConnected(context.Background()) {
			secondaryConn = connectionSuccess
		} else {
			secondaryConn = connectionFailed
		}
		secondaryStatus := models.DBStatus{
			Tier:       "secondary",
			Type:       h.secondaryDB.Type(),
			Connection: secondaryConn,
		}
		dbStatus = append(dbStatus, secondaryStatus)
	}

	var status = models.APIStatus{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
		APIStatus: "ok",
		DBStatus:  dbStatus,
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}

// GET /books/?table=books
// Get all books
func (h *Handler) GetAllBooks(c *gin.Context) {
	ctx := context.Background()

	// parse out author name in query params
	table, err := utils.GetParams(c, "table")
	log.Infof("table: %s\n", table)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No 'title' parameter provided"})
		return
	}

	data, err := h.primaryDB.All(ctx, table)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Unable to complete query"})
		return
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

	// struct to carry goroutine db insert results info
	type insertRes struct {
		Book models.Book // insert response book
		Err  error       // inert error
		DB   string      // "primary" or "secondary"
	}

	// create buffered channel for up to 2 goroutines
	results := make(chan insertRes, 2)

	// always insert primary synchronously
	go func(b models.InsertBookInput) {
		book, err := h.primaryDB.Insert(ctx, "books", newBook)
		results <- insertRes{Book: book, Err: err, DB: "primary"}
	}(newBook)

	// insert to the secondary DB if provided
	insertedToSecondary := false
	if h.secondaryDB != nil {
		insertedToSecondary = true
		go func(b models.InsertBookInput) {
			book, err := h.secondaryDB.Insert(ctx, "books", newBook)
			results <- insertRes{Book: book, Err: err, DB: "secondary"}
		}(newBook)
	}

	numResults := 1
	if insertedToSecondary {
		numResults = 2
	}

	var respBook = models.Book{}
	for i := 0; i < numResults; i++ {
		res := <-results
		if res.Err != nil {
			log.Errorf("Database (%s) insert failed: %v", res.DB, res.Err)
		}
		if res.DB == "primary" {
			respBook = res.Book
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": respBook})
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
	bookDocs, err := h.primaryDB.Get(ctx, "books", "Title", bookTitle)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Unable to complete query"})
		return
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
	authorBooks, err := h.primaryDB.Get(ctx, "books", "Author", author)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Unable to complete query"})
		return
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

	booksDeleted, err := h.primaryDB.Drop(ctx, "books", "Title", title)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Unable to complete query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booksDeleted})
}
