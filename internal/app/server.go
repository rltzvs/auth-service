package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"auth/internal/config"
	handlers "auth/internal/controller/http"
	"auth/internal/controller/http/middleware"
	"auth/internal/infra/postgres"
	"auth/internal/security"
	"auth/internal/service"
)

type Server struct {
	closer     *Closer
	logger     *zap.Logger
	httpServer *http.Server
	db         *postgres.DB
}

func NewServer(cfg *config.Config, logger *zap.Logger) *Server {
	closer := NewCloser()

	db, err := postgres.NewDbConnection(cfg)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	jwtManager := security.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.TokenDuration)
	passwordHasher := security.NewBcryptHasher()

	authMiddleware := middleware.AuthMiddleware(jwtManager, logger)
	authRepository := postgres.NewAuthRepository(db.Pool)
	authService := service.NewAuthService(authRepository, passwordHasher, jwtManager)
	authHandler := handlers.NewAuthHandler(authService, logger)

	r := chi.NewRouter()

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RecovererMiddleware(logger))
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		// r.Get("/user")
	})
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return &Server{
		closer:     closer,
		logger:     logger,
		httpServer: httpServer,
		db:         db,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.closer.Add(func(ctx context.Context) error {
		s.logger.Info("closing database connection")
		s.db.Close()
		return nil
	})

	s.closer.Add(func(ctx context.Context) error {
		s.logger.Info("shutting down HTTP server")
		return s.httpServer.Shutdown(ctx)
	})

	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("starting HTTP server", zap.String("addr", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("HTTP server failed", zap.Error(err))
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to start http server: %w", err)
		}
		return nil
	default:
		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.closer.Close(ctx)
}
