package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetParams(c *gin.Context, queryParam string) (string, error) {
	// get the value from the url parameters
	log.Printf("Query Param: %v", c.Query(queryParam))

	param, err := c.GetQuery(queryParam)
	if !err {
		log.Printf("No title provided...")
		return "", errors.New("empty name")
	}

	return param, nil
}

func SetupLogging(logFilename string) error {
	// logrus config
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	switch gin.Mode() {
	case "debug":
		log.Printf("Debug - Writing logs to stdout...\n")
	case "release":
		log.Printf("Release - writing logs to %s\n", logFilename)
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

func GetField(obj any, fieldName string) (any, error) {
	// use reflection to grab struct field value by name
	val := reflect.ValueOf(obj)

	// If pointer, resolve
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s", fieldName)
	}

	return field.Interface(), nil
}
