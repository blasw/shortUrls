package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// New is a middleware logger
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	//setting up path to the logger
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			//passing control to the next handler
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func test() {
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
	}

	params := url.Values{
		"info_hash":  []string{"123"},
		"peer_id":    []string{"456"},
		"port":       []string{"789"},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
	}
}
