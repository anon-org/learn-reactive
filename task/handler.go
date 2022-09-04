package task

import (
	"context"
	"encoding/json"
	"example/domain"
	"io"
	"log"
	"net/http"
	"time"
)

type handler struct {
	svc *service
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	defer func() {
		log.Println(time.Now().Sub(now))
	}()
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
	resp := &domain.TaskFetchResponse{
		Names: make([]string, 0),
	}

	for ch := range out {
		select {
		case <-ctx.Done():
			log.Println("ctx done")
		default:
			resp.Names = append(resp.Names, ch)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
			log.Println("send", str, "to out chan")
			out <- str
		}
	}()

	return out
}
