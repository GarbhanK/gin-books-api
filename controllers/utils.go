package controllers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

func getParams(c *gin.Context, queryParam string) (string, error) {
	log.Printf("Query Param: %v", c.Query(queryParam))

	param, err := c.GetQuery(queryParam)
	if err == false {
		log.Printf("No title provided...")
		return "", errors.New("empty name")
	}	

	return param, nil
}
