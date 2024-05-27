package utils

import (
	"fmt"
	"os"
	"io"
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

)

func GetParams(c *gin.Context, queryParam string) (string, error) {
	log.Printf("Query Param: %v", c.Query(queryParam))

	param, err := c.GetQuery(queryParam)
	if err == false {
		log.Printf("No title provided...")
		return "", errors.New("empty name")
	}	

	return param, nil
}


func SetupLogging() error {
	// logrus config
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	
	switch gin.Mode() {
	case "debug":
		fmt.Println("It's debug")
	case "release":
		fmt.Println("release muh")
		// Disable Console Color when running in 'release' mode
		gin.DisableConsoleColor()

		// Logging to a file.
		f, _ := os.Create("books.log")
		gin.DefaultWriter = io.MultiWriter(f)	
	default:
		fmt.Printf("%s is not a recognised mode...", gin.Mode())
	}
	
	return nil
}
