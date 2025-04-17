package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "calhounio_demo"
)

type Postgres struct {
	Client    *sql.DB
	projectId string
}

func (p Postgres) Conn() {
	psqlInfo := nil
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

}
