package task

import (
	"context"
	"log"
)

type service struct {
	repo *repository
}

func (s *service) Fetch(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)

	fetchCh := s.repo.Fetch(ctx, in)

	go func() {
		defer func() {
			log.Println("Fetch out close")
			close(out)
		}()

		for ch := range fetchCh {
			select {
			case <-ctx.Done():
				log.Println("ctx done")
			case out <- ch:
				log.Println("receiving", ch)
			}
		}
	}()

	return out
}
