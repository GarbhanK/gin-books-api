package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"strings"

	"github.com/garbhank/gin-books-api/db"
	"github.com/garbhank/gin-books-api/models"

	"github.com/gin-gonic/gin"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am root"})
}

// GET /ping
// get server status
func Ping(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"data": status})
}

// GET /books
// Get all books
func FindBooks(c *gin.Context) {

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

	c.JSON(http.StatusOK, gin.H{"data": firestoreBooks})
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()

	// Validate input
	var newBook models.CreateBookInput
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a DocumentReference
	_, _, err := client.Collection("books").Add(ctx, newBook)
	if err != nil {
        log.Fatalf("Failed adding document:\n%v", err)
	}

	book := models.Book{Title: newBook.Title, Author: newBook.Author}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// GET /books/:title
// Find a book
func FindBook(c *gin.Context) {
	
	// parse out author name in query params
	log.Printf("title query param %v", c.Query("title"))
	title, err := c.GetQuery("title")
	if err == false {
		log.Printf("No title provided...")
		return
	}

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()

	// array of books to return
	var bookDocs []models.Book

	// iterator over books collection in firestore
	iter := client.Collection("books").Where("Title", "==", title).Documents(ctx)
	defer iter.Stop() // add to clean up resources

	// loop until all documents matching title are added to books array
	for {
		var booksBuffer models.Book

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate:\n%v", err)
		}

		log.Println(doc.Data())
		if err := doc.DataTo(&booksBuffer); err != nil {
			log.Fatalf("can't cast docsnap to Book:\n%v", err)
		}

		bookDocs = append(bookDocs, booksBuffer)
	}

	c.JSON(http.StatusOK, gin.H{"data": bookDocs})
}

// 	// Validate Input
// 	var input models.UpdateBookInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	models.DB.Model(&book).Updates(input)

// 	c.JSON(http.StatusOK, gin.H{"data": book})
// }

func FindAuthor(c *gin.Context) {

	// parse out author name in query params
	log.Printf("author query param %v", c.Query("name"))
	author, err := c.GetQuery("name")
	if err == false {
		log.Printf("No name provided...")
		return
	}

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()

	// array of books to return
	var authorBooks []models.Book

	iter := client.Collection("books").Where("Author", "==", author).Documents(ctx)
	defer iter.Stop() // add to clean up resources

	// loop until all documents are added to books array
	for {
		log.Println("starting the loop...")
		var authorBooksBuffer models.Book

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate:\n%v", err)
		}

		log.Println(doc.Data())

		if err := doc.DataTo(&authorBooksBuffer); err != nil {
			log.Fatalf("can't cast docsnap to Book:\n%v", err)
		}

		// append record to return array
		authorBooks = append(authorBooks, authorBooksBuffer)
	}

	c.JSON(http.StatusOK, gin.H{"data": authorBooks})
}


// Delete a book
func DeleteBook(c *gin.Context) {

	// parse out author name in query params
	log.Printf("params: %v", c.Query("title"))

	title, err := c.GetQuery("title")
	if err == false {
		log.Printf("No title provided...")
		return
	}

	// create client
	ctx := context.Background()
	client := db.CreateFirestoreClient(ctx)
	defer client.Close()

	bulkwriter := client.BulkWriter(ctx)

	for {
		iter := client.Collection("books").Where("Title", "==", title).Documents(ctx)
		numDeleted := 0

		for {
			var booksBuffer models.Book
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate:\n%v", err)
			}

			log.Println(doc.Data())

			if err := doc.DataTo(&booksBuffer); err != nil {
				log.Fatalf("can't cast docsnap to Book:\n%v", err)
			}

			titleLower := strings.ToLower(title)
			parsedFirebaseTitle := strings.ToLower(booksBuffer.Title)

			// append record to array
			if (parsedFirebaseTitle == titleLower) {
				bulkwriter.Delete(doc.Ref)
				numDeleted++
			}
		}

		// if there are no docs to delete, process over
		if numDeleted == 0 {
			bulkwriter.End()
			return
		}
		bulkwriter.Flush()
	}

	log.Printf("Deleted record: %s", title)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
