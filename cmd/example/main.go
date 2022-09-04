package main

import (
	"example/task"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := task.MustProvideDatabase("sqlite3", ":memory:")
	defer db.Close()
	repo := task.ProvideRepository(db)
	svc := task.ProvideService(repo)
	hdl := task.ProvideHandler(svc)

	log.Println("listening at :8000")
	http.ListenAndServe(":8000", hdl)
}
