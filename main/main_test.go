package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"encoding/json"

	"github.com/stretchr/testify/assert"

	"github.com/garbhank/gin-books-api/models"
)

type PingStatusTest struct {
	Data  models.Status `json:"data"`
}

func TestPingRoute(t *testing.T) {
	router := setupRouter()
	currentTime := time.Now()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	pingResponse := models.Status{
		Timestamp:       currentTime.Format("2006-01-02 15:04:05"),
		APIStatus:       "ok",
		FirestoreStatus: "ok",
	}

	mockResponse := &PingStatusTest{Data: pingResponse}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}
