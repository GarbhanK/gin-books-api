package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-books-api/controllers"
	"github.com/garbhank/gin-books-api/database"
	"github.com/garbhank/gin-books-api/utils"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)

func init() {
	err := utils.SetupLogging("books.log")
	if err != nil {
		log.Fatalf("Failed to set up logging, %v\n", err)
	}
}

func setupRouter(handler *controllers.Handler, noCache bool) *gin.Engine {
	r := gin.Default()

	// cache endpoints which calls the Firestore db
	store := persistence.NewInMemoryStore(time.Second)
	ttl := time.Minute * time.Duration(utils.GetEnvInt("CACHE_TTL_MIN", 1))

	var (
		handleGetAllBooks,
		handleFindAuthor,
		handleFindBook func(c *gin.Context)
	)

	// logic to toggle caching on specific pages
	if noCache {
		log.Info("Setting up router with caching disabled...")
		handleGetAllBooks = handler.GetAllBooks
		handleFindAuthor = handler.FindAuthor
		handleFindBook = handler.FindBook
	} else {
		log.Info("Setting up router with caching enabled...")
		handleGetAllBooks = cache.CachePage(store, ttl, handler.GetAllBooks)
		handleFindAuthor = cache.CachePage(store, ttl, handler.FindAuthor)
		handleFindBook = cache.CachePage(store, ttl, handler.FindBook)
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", handler.Root)
		v1.GET("/ping", handler.Ping)
		v1.GET("/books/", handleGetAllBooks)
		v1.GET("/books/author/", handleFindAuthor)
		v1.GET("/books/title/", handleFindBook)
		v1.POST("/books", handler.CreateBook)
		v1.DELETE("/books/", handler.DeleteBook)
	}

	return r
}

func main() {
	// parse database environment variables
	var primary, secondary string // 'memory', 'firestore', or 'postgres'
	if v := os.Getenv("PRIMARY_DB"); v != "" {
		primary = v
	}
	if v := os.Getenv("SECONDARY_DB"); v != "" {
		secondary = v
	}

	// parse cache environment variable
	var noCache bool // 'true' or 'false'
	if v := os.Getenv("ENABLE_CACHE"); v != "" {
		if strings.ToLower(v) == "true" {
			noCache = true
		} else {
			noCache = false
		}
	}

	var primaryDB, secondaryDB database.Database

	// create Database structs based on input name
	primaryDB = database.GetDB(primary)
	err := primaryDB.Conn(context.Background())
	if err != nil {
		log.Errorf("Unable to connect to primary database: %v\n", err)
	}

	// if set, also get the database type of the secondary
	if secondary != "" {
		secondaryDB = database.GetDB(secondary)
		err := secondaryDB.Conn(context.Background())
		if err != nil {
			log.Errorf("Unable to connect to primary database: %v\n", err)
		}
		if primaryDB.Type() == secondaryDB.Type() {
			log.Warnf("Primary and Secondary databases are of the same type: %v, %v", primaryDB.Type(), secondaryDB.Type())
		}
	}
	log.Infof("Primary Database: %v\n", primaryDB.Type())
	log.Infof("Secondary Database: %v\n", secondaryDB.Type())

	handler := controllers.NewHandler(primaryDB, secondaryDB)
	r := setupRouter(handler, noCache)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
