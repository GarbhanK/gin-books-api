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

func NewMemoryDB() *MemoryDB {
	fmt.Println("Creating new memoryBD...")
	memoryMap := make(map[string][]models.Book)
	return &MemoryDB{
		Client: memoryMap,
	}
}

func (m *MemoryDB) Conn(ctx context.Context) error {
	if m.Client == nil {
		return errors.New("No in-memory database found!")
	}
	fmt.Printf("Connected to MemoryDB! :: %v\n", m.Client)
	return nil
}

func (m *MemoryDB) Close() error {
	// clear(m.Client)
	// if len(m.Client) != 0 {
	// 	return errors.New("Failed to close DB connection")
	// }

	return nil
}

func (m *MemoryDB) Insert(ctx context.Context, table string, data models.InsertBookInput) (models.Book, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Printf("insert map before: %v\n", m.Client)

	// create new book struct
	newBook := models.Book(data)

	// append new book to the 'table' array
	m.Client[table] = append(m.Client[table], newBook)

	fmt.Printf("insert map after: %v\n", m.Client)
	return newBook, nil
}

func (m *MemoryDB) Get(ctx context.Context, table, key, val string) ([]models.Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Printf("map: %v\n", m.Client)

	books, ok := m.Client[table]
	if !ok {
		return []models.Book{}, fmt.Errorf("Data not found for: %v", key)
	}

	// filter books array
	var matchingBooks []models.Book

	for _, book := range books {
		// use reflect to get struct field by string
		fieldValue, err := utils.GetField(book, key)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		fmt.Printf("fieldValue: %s\n", fieldValue)
		if fieldValue == val {
			matchingBooks = append(matchingBooks, book)
		}
	}

	return matchingBooks, nil
}

func (m *MemoryDB) Drop(ctx context.Context, table, key, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Client, key)
	return nil
}
