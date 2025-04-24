package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/garbhank/gin-books-api/models"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "db-demo"
)

type Postgres struct {
	Client   *sql.DB
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func NewPostgres() *Postgres {
	// get port separately because it's an integer
	port, err := strconv.Atoi(os.Getenv("PGSQL_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	return &Postgres{
		host:     os.Getenv("PGSQL_HOST"),
		port:     port,
		user:     os.Getenv("PGSQL_USER"),
		password: os.Getenv("PGSQL_PASSWORD"),
		dbname:   os.Getenv("PGSQL_DBNAME"),
	}
}

func (p Postgres) Conn(ctx context.Context) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Error connecting to Postgres: %v", err)
		return err
	}
	defer db.Close()

	fmt.Println("Successfully connected!")
	p.Client = db

	return nil
}

func (p Postgres) Close() error {
	err := p.Client.Close()
	if err != nil {
		return fmt.Errorf("Error closing database connection: %v\n", err)
	}
	return nil
}

func (p Postgres) Get(ctx context.Context, table, key, val string) ([]models.Book, error) {
	return []models.Book{}, nil
}

func (p Postgres) Drop(ctx context.Context, table, key, val string) (int, error) {
	return 0, nil
}

func (p Postgres) All(ctx context.Context, table string) ([]models.Book, error) {
	return []models.Book{}, nil
}

func (p Postgres) Insert(ctx context.Context, table string, data models.InsertBookInput) (models.Book, error) {
	return models.Book{}, nil
}
