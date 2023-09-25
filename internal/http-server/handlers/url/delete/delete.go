package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	storage2 "httpserver/internal/storage"
	"httpserver/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, storage *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log.Info("delete url", slog.String("op", op))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			http.Error(w, "invalid argument", http.StatusBadRequest)
			return
		}

		err := storage.DeleteURL(alias)
		if errors.Is(err, storage2.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Error("failed to delete url from db")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
