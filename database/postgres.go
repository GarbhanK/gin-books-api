package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/garbhank/gin-books-api/models"
	"github.com/garbhank/gin-books-api/utils"
	_ "github.com/lib/pq"
)

const (
	host     = "postgres" // hostname set for docker-compose internal dns
	port     = "5432"
	user     = "gin"
	password = "ginpass"
	dbname   = "books"
)

type Postgres struct {
	Client   *sql.DB
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewPostgres() *Postgres {
	// verify port string is a valid integer
	if _, err := strconv.Atoi(utils.GetenvDefault("PGSQL_PORT", port)); err != nil {
		log.Fatal(err)
	}

	// if running as a container, use the docker network name instead of localhost
	var hostname string = "localhost"
	if useContainerNetworking := os.Getenv("CONTAINER_NETWORKING"); useContainerNetworking == "true" {
		hostname = host
	}
	// log.Infof("hostname: %v\n", hostname)

	return &Postgres{
		port:     utils.GetenvDefault("PGSQL_PORT", port),
		host:     utils.GetenvDefault("PGSQL_HOST", hostname),
		user:     utils.GetenvDefault("PGSQL_USER", user),
		password: utils.GetenvDefault("PGSQL_PASSWORD", password),
		dbname:   utils.GetenvDefault("PGSQL_DBNAME", dbname),
	}
}

func (p *Postgres) Type() string { return "postgres" }

func (p *Postgres) IsConnected(ctx context.Context) bool {
	if err := p.Client.PingContext(ctx); err != nil {
		log.Printf("DB ping failed: %v\n", err)
		return false
	}

	return true
}

func (p *Postgres) Setup(ctx context.Context) error {
	// create the table if not present
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
		id	   TEXT PRIMARY KEY,
		title  VARCHAR(255),
		author VARCHAR(255)
	);`, "books")

	_, err := p.Client.ExecContext(ctx, createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func (p *Postgres) Conn(ctx context.Context) error {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.host, p.port, p.user, p.password, p.dbname,
	)
	log.Printf("psqlInfo: %v\n", psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Error connecting to Postgres: %v", err)
		return err
	}

	log.Info("Successfully connected!")
	p.Client = db

	return nil
}

func (p *Postgres) Close() error {
	err := p.Client.Close()
	if err != nil {
		return fmt.Errorf("error closing database connection: %v", err)
	}
	return nil
}

func (p *Postgres) Get(ctx context.Context, table, key, val string) ([]models.Book, error) {
	// filter based on the selected column and value
	selectQuery := fmt.Sprintf(`SELECT title, author FROM "%s" WHERE "%s" = $1`, table, strings.ToLower(key))
	rows, err := p.Client.QueryContext(ctx, selectQuery, val)
	if err != nil {
		return nil, fmt.Errorf("error while performing query: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v\n", err)
		}
	}()

	// create a slice with 0 elements
	books := []models.Book{}

	log.Printf("Iterating through rows...")
	for rows.Next() {
		var b models.Book

		if err := rows.Scan(&b.Title, &b.Author); err != nil {
			return books, err
		}
		books = append(books, b)
	}
	if err = rows.Err(); err != nil {
		return books, err
	}

	return books, nil
}

func (p *Postgres) Drop(ctx context.Context, table, key, val string) (int, error) {
	// TODO: make sure casting int64 to int isn't causing any trouble

	// delete data from the table based on the input table/key/value
	deleteQuery := fmt.Sprintf(`DELETE FROM "%s" WHERE "%s" = $1`, table, strings.ToLower(key))
	res, err := p.Client.ExecContext(ctx, deleteQuery, val)
	if err != nil {
		return 0, fmt.Errorf("error while performing query: %v", err)
	}

	// get the number of rows deleted by the query
	n, err := res.RowsAffected()
	if err != nil {
		return int(n), fmt.Errorf("error while getting the number of rows affected by the DELETE command: %v", err)
	}

	return int(n), nil
}

func (p *Postgres) All(ctx context.Context, table string) ([]models.Book, error) {

	// filter based on the selected column and value
	selectQuery := fmt.Sprintf(`SELECT title, author FROM "%s" LIMIT 100`, table)
	rows, err := p.Client.QueryContext(ctx, selectQuery)
	if err != nil {
		return nil, fmt.Errorf("error while performing query: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v\n", err)
		}
	}()

	// create a slice with 0 elements
	books := []models.Book{}

	log.Printf("Iterating through rows...")
	for rows.Next() {
		var b models.Book

		if err := rows.Scan(&b.Title, &b.Author); err != nil {
			return books, err
		}
		books = append(books, b)
	}
	if err = rows.Err(); err != nil {
		return books, err
	}

	return books, nil
}

func (p *Postgres) Insert(ctx context.Context, table string, data models.InsertBookInput) (models.Book, error) {
	book := models.Book{
		Id:     utils.UUID(),
		Title:  data.Title,
		Author: data.Author,
	}

	if p.Client == nil {
		return models.Book{}, fmt.Errorf("Database client is not initialised")
	}

	if !utils.IsSafeIdentifier(table) {
		return models.Book{}, fmt.Errorf("invalid table name: %v", table)
	}

	// insert new book into db table
	insertQuery := fmt.Sprintf(`INSERT INTO "%s" (id, title, author) VALUES ($1, $2, $3)`, table)

	_, err := p.Client.ExecContext(ctx, insertQuery, book.Id, book.Title, book.Author)
	if err != nil {
		return book, fmt.Errorf("error while performing query: %v", err)
	}

	return book, nil
}
