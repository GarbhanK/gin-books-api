package database

import (
	"context"

	"github.com/garbhank/gin-books-api/models"
)

// interface for multiple Database backends
type Database interface {
	Conn(ctx context.Context) error
	Close() error
	All(ctx context.Context, table string) ([]models.Book, error)
	Get(ctx context.Context, table, key, val string) ([]models.Book, error)
	Drop(ctx context.Context, table, key, val string) (int, error)
	Insert(ctx context.Context, table string, data models.InsertBookInput) (models.Book, error)
	// Setup() (could be used to create table if not exists, etc)
	IsConnected(ctx context.Context) bool // (test db connection, currently ping just checks for nil)
}
