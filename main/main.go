package main

import (
	"os"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-books-api/controllers"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)



func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	debugEnabled := gin.IsDebugging()

	if (!debugEnabled) {
		// Disable Console Color when running in 'release' mode
		gin.DisableConsoleColor()

		// Logging to a file.
		f, _ := os.Create("books.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	// cache endpoints which calls the Firestore db
	store := persistence.NewInMemoryStore(time.Second)
	ttl := time.Minute * 1 // todo: os.GetEnv

	v1 := r.Group("/v1")
	{
		v1.GET("/", controllers.Root)
		v1.GET("/ping", controllers.Ping)
		v1.GET("/books", cache.CachePage(store, ttl, controllers.FindBooks))
		v1.GET("/books/author/", cache.CachePage(store, ttl, controllers.FindAuthor))
		v1.GET("/books/title/", cache.CachePage(store, ttl, controllers.FindBook))
		v1.POST("/books", controllers.CreateBook)
		v1.DELETE("/books/", controllers.DeleteBook)
		// v1.PATCH("books/:id", controllers.UpdateBook)
	}


	return r
}

func main() {

	r := setupRouter()

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
