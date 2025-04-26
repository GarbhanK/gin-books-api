package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetParams(c *gin.Context, queryParam string) (string, error) {
	// get the value from the url parameters
	log.Infof("Query Param: %v", c.Query(queryParam))

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
		log.Printf("%s is not a recognised mode...", gin.Mode())
	}

	return nil
}

func GetField(obj any, fieldName string) (any, error) {
	// use reflection to grab struct field value by name
	v := reflect.ValueOf(obj)
	log.Printf("GetField val: %v\n", v)

	// get the value from the reflected value points to
	v = reflect.Indirect(v)

	fieldVal := v.FieldByName(fieldName)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("no such field: %s", fieldName)
	}

	return fieldVal.Interface(), nil
}

func GetenvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func IsSafeIdentifier(id string) bool {
	// allows only letters, numbers, and underscores
	re := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return re.MatchString(id)
}
