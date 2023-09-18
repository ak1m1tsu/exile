package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware.logger"),
		)

		log.Info("Logger middleware initialized")

		fn := func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("host", r.Host),
				slog.String("referer", r.Referer()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				log.Info("request completed",
					slog.Int("status_code", ww.Status()),
					slog.Duration("duration", time.Since(start)),
					slog.Int("bytes_written", ww.BytesWritten()),
					slog.String("content_type", ww.Header().Get("Content-Type")),
				)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
