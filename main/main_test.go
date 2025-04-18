package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/garbhank/gin-books-api/controllers"
	"github.com/garbhank/gin-books-api/database"
	"github.com/garbhank/gin-books-api/models"
)

type PingStatusTest struct {
	Data models.Status `json:"data"`
}

func TestPingRoute(t *testing.T) {
	handler := controllers.NewHandler(&database.MemoryDB{})
	router := setupRouter(*handler)
	currentTime := time.Now()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	pingResponse := models.Status{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
		APIStatus: "ok",
		DBStatus:  "ok",
	}

	mockResponse := &PingStatusTest{Data: pingResponse}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}
