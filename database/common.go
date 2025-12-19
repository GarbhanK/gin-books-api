package database

import (
	"context"

	log "github.com/sirupsen/logrus"

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
	IsConnected(ctx context.Context) bool // (test db connection, currently ping just checks for nil)
	Setup(ctx context.Context) error
	Type() string
}

func GetDB(dbName string) Database {
	var db Database

	switch dbName {
	case "firestore":
		db = NewFirestore()
	case "memorydb":
		db = NewMemoryDB(nil)
	case "postgres":
		db = NewPostgres()
	default:
		log.Fatalf("Unknown DB type: %s", dbName)
	}

	return db
}
