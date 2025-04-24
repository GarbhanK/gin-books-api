package main

import (
	"bytes"
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

type pingStatusTest struct {
	Data models.Status `json:"data"`
}

type postBookTest struct {
	Data models.Book `json:"data"`
}

type getBookTitleTest struct {
	Data []models.Book `json:"data"`
}

type deleteBookTest struct {
	Data int `json:"data"`
}

var seedDataSingle map[string][]models.Book = map[string][]models.Book{
	"books": {
		{Author: "Jorge Luis Borges", Title: "Fictions"},
	},
}

var seedDataMultiple map[string][]models.Book = map[string][]models.Book{
	"books": {
		{Author: "Jorge Luis Borges", Title: "Fictions"},
		{Author: "Jorge Luis Borges", Title: "The Aleph"},
		{Author: "John Smith", Title: "Fictions"},
	},
}

func TestGetPingRoute(t *testing.T) {
	handler := controllers.NewHandler(database.NewMemoryDB(nil))
	router := setupRouter(*handler, false)
	currentTime := time.Now()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	pingResponse := models.Status{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
		APIStatus: "ok",
		DBStatus:  "ok",
	}

	mockResponse := &pingStatusTest{Data: pingResponse}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}

func TestPostBookRoute(t *testing.T) {
	handler := controllers.NewHandler(database.NewMemoryDB(nil))
	router := setupRouter(*handler, false)

	w := httptest.NewRecorder()

	jsonBody := []byte(`{"Author":"Jorge Luis Borges","Title":"Fictions"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/books",
		bodyReader,
	)
	router.ServeHTTP(w, req)

	// create the expected response
	mockResponse := &postBookTest{Data: models.Book{Title: "Fictions", Author: "Jorge Luis Borges"}}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}

func TestGetBookTitleSingleRoute(t *testing.T) {
	// create memoryDB with seed data
	handler := controllers.NewHandler(database.NewMemoryDB(seedDataSingle))
	router := setupRouter(*handler, false)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/books/title/?title=Fictions", nil)
	router.ServeHTTP(w, req)

	// create the expected response
	mockResponse := &getBookTitleTest{
		Data: []models.Book{
			{Title: "Fictions", Author: "Jorge Luis Borges"},
		},
	}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}

func TestGetBookTitleMultipleRoute(t *testing.T) {
	// create memoryDB with seed data
	handler := controllers.NewHandler(database.NewMemoryDB(seedDataMultiple))
	router := setupRouter(*handler, false)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/books/title/?title=Fictions", nil)
	router.ServeHTTP(w, req)

	// create the expected response
	mockResponse := &getBookTitleTest{
		Data: []models.Book{
			{Title: "Fictions", Author: "Jorge Luis Borges"},
			{Title: "Fictions", Author: "John Smith"},
		},
	}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}

func TestGetBookAuthorSingleRoute(t *testing.T) {
	handler := controllers.NewHandler(database.NewMemoryDB(seedDataSingle))
	router := setupRouter(*handler, false)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/books/author/?name=Jorge+Luis+Borges", nil)
	router.ServeHTTP(w, req)

	// create the expected response
	mockResponse := &getBookTitleTest{
		Data: []models.Book{
			{Title: "Fictions", Author: "Jorge Luis Borges"},
		},
	}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}

func TestGetBookAuthorMultipleRoute(t *testing.T) {
	handler := controllers.NewHandler(database.NewMemoryDB(seedDataMultiple))
	router := setupRouter(*handler, false)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/books/author/?name=Jorge+Luis+Borges", nil)
	router.ServeHTTP(w, req)

	// create the expected response
	mockResponse := &getBookTitleTest{
		Data: []models.Book{
			{Author: "Jorge Luis Borges", Title: "Fictions"},
			{Author: "Jorge Luis Borges", Title: "The Aleph"},
		},
	}
	b, _ := json.Marshal(mockResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}

// DELETE api/v1/books/?title=Fictions
func TestDeleteBookPositive(t *testing.T) {
	handler := controllers.NewHandler(database.NewMemoryDB(seedDataMultiple))
	router := setupRouter(*handler, false)

	w1 := httptest.NewRecorder()
	deleteReq, _ := http.NewRequest(http.MethodDelete, "/api/v1/books/?title=Fictions", nil)
	router.ServeHTTP(w1, deleteReq)

	// create the expected response from delete
	mockDeleteResponse := &deleteBookTest{2}
	bDelete, _ := json.Marshal(mockDeleteResponse)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, string(bDelete), w1.Body.String())

	// query the db to see if title has been removed
	w2 := httptest.NewRecorder()
	getReq, _ := http.NewRequest(http.MethodGet, "/api/v1/books/title/?title=Fictions", nil)
	router.ServeHTTP(w2, getReq)
	mockGetResponse := &getBookTitleTest{Data: []models.Book{}}
	bGet, _ := json.Marshal(mockGetResponse)

	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, string(bGet), w2.Body.String())
}

// DELETE /api/v1/books/?title=Fictions"
func TestDeleteBookNegative(t *testing.T) {
	// create memoryDB with seed data
	handler := controllers.NewHandler(database.NewMemoryDB(seedDataSingle))
	router := setupRouter(*handler, false)

	w1 := httptest.NewRecorder()
	deleteReq, _ := http.NewRequest(http.MethodDelete, "/api/v1/books/?title=NoSuchBook", nil)
	router.ServeHTTP(w1, deleteReq)

	// create the expected response from delete
	mockDeleteResponse := &deleteBookTest{0}
	bDelete, _ := json.Marshal(mockDeleteResponse)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, string(bDelete), w1.Body.String())

	// query the db to see if title has been removed
	w2 := httptest.NewRecorder()
	getReq, _ := http.NewRequest(http.MethodGet, "/api/v1/books/title/?title=NoSuchBook", nil)
	router.ServeHTTP(w2, getReq)
	mockGetResponse := &getBookTitleTest{Data: []models.Book{}}
	bGet, _ := json.Marshal(mockGetResponse)

	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, string(bGet), w2.Body.String())
}
