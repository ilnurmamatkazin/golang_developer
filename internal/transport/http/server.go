package http

import (
	"context"
	"errors"
	"net/http"

	"test/internal/storage"
	"test/pkg/config"
	"test/pkg/log"
)

type HttpServer struct {
	cfg    *config.Config
	stor   *storage.Storage
	server *http.Server
}

// Инициализируем http сервер
func Init(ctx context.Context, stor *storage.Storage) *HttpServer {

	log.Info.Println("HTTP сервер создан.")

	return &HttpServer{
		cfg:  ctx.Value(config.ContextKeyConfig).(*config.Config),
		stor: stor,
	}

}

// Запускаем http сервер
func (s *HttpServer) Run() error {

	s.server = &http.Server{
		Addr:              s.cfg.HTTP.Host,
		Handler:           s.NewRouter(),
		ReadTimeout:       s.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: s.cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.HTTP.WriteTimeout,
	}

	log.Info.Println("HTTP сервер запущен.")

	if err := s.server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Info.Println("HTTP сервер остановлен")
		default:
			return err
		}
	}

	return nil

}

// Останавливаем http сервер
func (s *HttpServer) CloseGracefully(ctx context.Context) {

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error.Printf("HTTP shutdown error: %v", err)

		if err = s.server.Close(); err != nil {
			log.Error.Printf("HTTP close error: %v", err)
		}
	}

}
