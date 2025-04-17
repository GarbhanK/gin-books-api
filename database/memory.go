package database

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/garbhank/gin-books-api/models"
	"github.com/garbhank/gin-books-api/utils"
)

// fake in memory db for demo/testing
type MemoryDB struct {
	data map[string][]models.Book
	mu   sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{}
}

func (m *MemoryDB) Conn(ctx context.Context) error {
	m.data = make(map[string][]models.Book)
	return nil
}

func (m *MemoryDB) Insert(ctx context.Context, table string, data models.InsertBookInput) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[table] = append(m.data[table], models.Book(data))
	return nil
}

func (m *MemoryDB) Get(ctx context.Context, table, key, val string) ([]models.Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	books, ok := m.data[table]
	if !ok {
		return []models.Book{}, errors.New("not found")
	}

	// filter books array
	matchingBooks := []models.Book{}
	for _, book := range books {
		// use reflect to get struct field by string
		fieldValue, err := utils.GetField(models.Book{}, key)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		if fieldValue == val {
			matchingBooks = append(matchingBooks, book)
		}
	}

	return matchingBooks, nil
}

func (m *MemoryDB) Drop(ctx context.Context, table, key, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryDB) Close() error {
	return nil
}
