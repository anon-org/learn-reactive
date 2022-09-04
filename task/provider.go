package task

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func MustProvideDatabase(driver, dsn string) *sql.DB {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec(`
CREATE TABLE IF NOT EXISTS tasks(
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
)`)
	db.Exec(`INSERT INTO tasks VALUES
                      ('id1', 'name1'),
                      ('id2', 'name2'),
                      ('id3', 'name3'),
                      ('id4', 'name4')
`)

	return db
}

func ProvideRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func ProvideService(repo *repository) *service {
	return &service{
		repo: repo,
	}
}

func ProvideHandler(svc *service) *handler {
	return &handler{
		svc: svc,
	}
}
