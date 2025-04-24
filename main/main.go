package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-books-api/controllers"
	"github.com/garbhank/gin-books-api/database"
	"github.com/garbhank/gin-books-api/utils"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)

var db database.Database

func init() {
	err := utils.SetupLogging("books.log")
	if err != nil {
		log.Fatalf("Failed to set up logging, %v\n", err)
	}
}

func setupRouter(handler controllers.Handler, enableCacheing bool) *gin.Engine {
	r := gin.Default()

	// cache endpoints which calls the Firestore db
	store := persistence.NewInMemoryStore(time.Second)
	ttl := time.Minute * 1 // todo: os.GetEnv

	// logic to toggle caching on specific pages
	handleGetAllBooks := handler.GetAllBooks
	handleFindAuthor := handler.FindAuthor
	handleFindBook := handler.FindBook

	if enableCacheing {
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
	// parse command line flags
	dbType := flag.String("db", "memory", "Database type: 'memory', 'firestore', or 'postgres'")
	flag.Parse()

	switch *dbType {
	case "firestore":
		db = database.NewFirestore()
	case "memory":
		db = database.NewMemoryDB(nil)
	case "postgres":
		db = database.NewPostgres()
	default:
		log.Fatalf("Unknown DB type: %s", *dbType)
	}

	handler := controllers.NewHandler(db)
	r := setupRouter(*handler, true)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
