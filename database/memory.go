package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/garbhank/gin-books-api/models"
	"github.com/garbhank/gin-books-api/utils"
)

// fake in memory db for demo/testing
type MemoryDB struct {
	Client map[string][]models.Book
	mu     sync.RWMutex
}

func NewMemoryDB(data map[string][]models.Book) *MemoryDB {
	log.Println("Creating new memoryBD...")
	var memoryMap map[string][]models.Book

	// if no seed data provided, init a new map
	if data == nil {
		memoryMap = make(map[string][]models.Book)
	} else {
		memoryMap = data
	}

	return &MemoryDB{
		Client: memoryMap,
	}
}

func (m *MemoryDB) Conn(ctx context.Context) error {
	if m.Client == nil {
		return errors.New("no in-memory database found")
	}
	fmt.Printf("Connected to MemoryDB! :: %v\n", m.Client)
	return nil
}

func (m *MemoryDB) Close() error {
	// Map will be cleaned up by the GC, no manual `clear()` needed
	return nil
}

func (m *MemoryDB) Insert(ctx context.Context, table string, data models.InsertBookInput) (models.Book, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Printf("insert map before: %v\n", m.Client)

	// create new book struct
	newBook := models.Book(data)

	// append new book to the 'table' array
	m.Client[table] = append(m.Client[table], newBook)

	log.Printf("insert map after: %v\n", m.Client)
	return newBook, nil
}

func (m *MemoryDB) Get(ctx context.Context, table, key, val string) ([]models.Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	log.Printf("map: %v\n", m.Client)

	books, ok := m.Client[table]
	if !ok {
		return []models.Book{}, fmt.Errorf("data not found for: %v", key)
	}

	matchingBooks := []models.Book{}

	// filter books array
	for _, book := range books {
		// use reflect to get struct field by string
		fieldValue, err := utils.GetField(book, key)
		if err != nil {
			return nil, fmt.Errorf("error getting field value: %v", err)
		}

		if fieldValue == val {
			matchingBooks = append(matchingBooks, book)
		}
	}

	return matchingBooks, nil
}

func (m *MemoryDB) Drop(ctx context.Context, table, key, val string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var filteredBooks []models.Book
	booksFound := 0

	log.Printf("pre drop map: %v\n", m.Client[table])

	for _, book := range m.Client[table] {
		// get book field value
		fieldValue, err := utils.GetField(book, key)
		if err != nil {
			return booksFound, fmt.Errorf("error: %v", err)
		}

		// if value matches, don't append to the output array
		if fieldValue == val {
			log.Printf("Book to delete: %v\n", book)
			booksFound += 1
			continue
		}

		// if it's not to be deleted, add it to the re-created array
		filteredBooks = append(filteredBooks, book)
	}

	m.Client[table] = filteredBooks
	log.Printf("Post-drop post-loop map: %v\n", m.Client[table])
	return booksFound, nil
}

func (m *MemoryDB) All(ctx context.Context, table string) ([]models.Book, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	allRecords := m.Client[table]

	return allRecords, nil
}

func (m *MemoryDB) IsConnected(ctx context.Context) bool {
	return m.Client != nil
}

func (f *MemoryDB) Type() string { return "memorydb" }
