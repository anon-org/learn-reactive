package task

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type repository struct {
	db *sql.DB
}

func (r *repository) Fetch(ctx context.Context, in <-chan string) <-chan string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)

	out := make(chan string)
	go func() {
		defer func() {
			log.Println("Fetch out close and cancel ctx")
			cancel()
			close(out)
		}()

		ids := make([]string, 0)
		for id := range in {
			select {
			case <-ctx.Done():
				log.Println("ctx Done")
				return
			default:
				log.Println("receive", id)
				ids = append(ids, id)
			}
		}

		if len(ids) == 0 {
			log.Println("ids empty")
			return
		}

		rows, err := r.db.QueryContext(ctx, `SELECT name FROM tasks WHERE id IN($1, $2)`, ids[0], ids[1])
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Println(err)
				return
			}
			log.Println("sending", name)
			out <- name
		}
	}()

	return out
}
