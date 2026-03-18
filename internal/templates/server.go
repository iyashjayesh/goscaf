package templates

import "github.com/iyashjayesh/goscaf/internal/config"

// ServerGo returns the server.go template for the given framework.
func ServerGo(framework config.Framework) string {
	switch framework {
	case config.FrameworkGin:
		return serverGin
	case config.FrameworkFiber:
		return serverFiber
	case config.FrameworkChi:
		return serverChi
	case config.FrameworkEcho:
		return serverEcho
	case config.FrameworkGorilla:
		return serverGorilla
	default:
		return serverNone
	}
}

const serverGin = `package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/handler"
)

// Server holds the HTTP server and dependencies.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	router *gin.Engine
}

// New creates a new Server.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(slogMiddleware(logger))

	s := &Server{cfg: cfg, logger: logger, router: router}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	h := handler.New(s.cfg)

	s.router.GET("/health", h.Health)

	v1 := s.router.Group("/api/v1")
	_ = v1 // add your routes here
}

// Start starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.App.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("server listening", "port", s.cfg.App.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func slogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip", c.ClientIP(),
		)
	}
}
`

const serverFiber = `package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/handler"
)

// Server holds the Fiber app and dependencies.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	app    *fiber.App
}

// New creates a new Server.
func New(cfg *config.Config, log *slog.Logger) *Server {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  15,
		WriteTimeout: 15,
		IdleTimeout:  60,
	})
	app.Use(recover.New())
	app.Use(logger.New())

	s := &Server{cfg: cfg, logger: log, app: app}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	h := handler.New(s.cfg)

	s.app.Get("/health", h.Health)

	v1 := s.app.Group("/api/v1")
	_ = v1 // add your routes here
}

// Start starts the Fiber server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.cfg.App.Port)
	s.logger.Info("server listening", "port", s.cfg.App.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := s.app.Listen(addr); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server...")
		return s.app.ShutdownWithTimeout(10e9)
	case err := <-errCh:
		return err
	}
}
`

const serverChi = `package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/handler"
)

// Server holds the HTTP server and dependencies.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	router chi.Router
}

// New creates a new Server.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(slogMiddleware(logger))

	s := &Server{cfg: cfg, logger: logger, router: r}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	h := handler.New(s.cfg)

	s.router.Get("/health", h.Health)
	s.router.Route("/api/v1", func(r chi.Router) {
		// add your routes here
		_ = json.NewEncoder
	})
}

// Start starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.App.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("server listening", "port", s.cfg.App.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func slogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"latency", time.Since(start).String(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
`

const serverEcho = `package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/handler"
)

// Server holds the Echo server and dependencies.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	echo   *echo.Echo
}

// New creates a new Server.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())

	s := &Server{cfg: cfg, logger: logger, echo: e}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	h := handler.New(s.cfg)

	s.echo.GET("/health", h.Health)

	v1 := s.echo.Group("/api/v1")
	_ = v1 // add your routes here
}

// Start starts the Echo server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.App.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("server listening", "port", s.cfg.App.Port)
		if err := s.echo.StartServer(httpSrv); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.echo.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
`

const serverGorilla = `package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/handler"
)

// Server holds the HTTP server and dependencies.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	router *mux.Router
}

// New creates a new Server.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	r := mux.NewRouter()

	s := &Server{cfg: cfg, logger: logger, router: r}
	r.Use(s.loggingMiddleware)
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	h := handler.New(s.cfg)

	s.router.HandleFunc("/health", h.Health).Methods(http.MethodGet)

	api := s.router.PathPrefix("/api/v1").Subrouter()
	_ = api // add your routes here
}

// Start starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.App.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("server listening", "port", s.cfg.App.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"latency", time.Since(start).String(),
		)
	})
}
`

const serverNone = `package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/handler"
)

// Server holds the HTTP server and dependencies.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	mux    *http.ServeMux
}

// New creates a new Server.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	s := &Server{cfg: cfg, logger: logger, mux: mux}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	h := handler.New(s.cfg)

	s.mux.HandleFunc("/health", h.Health)
	// add /api/v1 routes here
	_ = json.NewEncoder
}

// Start starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.App.Port),
		Handler:      s.mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("server listening", "port", s.cfg.App.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
`
