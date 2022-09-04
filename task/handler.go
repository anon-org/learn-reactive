package task

import (
	"context"
	"encoding/json"
	"example/domain"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type handler struct {
	svc *service
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out := make(chan string)
	ctx, cancel := context.WithCancel(r.Context())
	defer func() {
		log.Println("ServeHTTP cancel ctx")
		defer cancel()
	}()

	streamCh := stream(r.Context(), r.Body)
	if streamCh == nil {
		log.Println("streamCh is nil")
		return
	}

	fetchCh := h.svc.Fetch(r.Context(), streamCh)

	go func() {
		defer func() {
			log.Println("ServeHTTP out close")
			close(out)
		}()
		for ch := range fetchCh {
			log.Println("send", ch, "to out chan")
			out <- ch
		}
	}()

	w.WriteHeader(http.StatusOK)
	for ch := range out {
		select {
		case <-ctx.Done():
			log.Println("ctx done")
		default:
			fmt.Fprintf(w, "%s\n", ch)
		}
	}
}

func stream(ctx context.Context, body io.Reader) <-chan string {
	out := make(chan string)

	var req domain.TaskFetchRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		log.Println(err)
		return nil
	}

	go func() {
		defer func() {
			log.Println("stream out close")
			close(out)
		}()
		for _, str := range req.IDs {
			time.Sleep(time.Second)
			log.Println("send", str, "to out chan")
			out <- str
		}
	}()

	return out
}
