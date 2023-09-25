package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "httpserver/internal/lib/api/response"
	storage2 "httpserver/internal/storage"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, storage URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("alias is empty")
			http.Error(w, "alias is empty", http.StatusBadRequest)
			return
		}

		log.Info("redirecting", slog.String("alias", alias))

		url, err := storage.GetURL(alias)
		if errors.Is(err, storage2.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("url not found"))
		}

		if err != nil {
			log.Error("failed to get url from storage")
			render.JSON(w, r, resp.Error("internal error"))
		}

		http.Redirect(w, r, url, http.StatusMovedPermanently)

	}
}
