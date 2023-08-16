package server

import (
	"Interior_Visualization_Shop/app/internal/appeal"
	"Interior_Visualization_Shop/app/internal/auth"
	"Interior_Visualization_Shop/app/internal/user"
	"Interior_Visualization_Shop/app/pkg/config"
	"Interior_Visualization_Shop/app/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type Server struct {
	srv     *http.Server
	log     *logger.Logger
	cfg     *config.Config
	handler *httprouter.Router
}

func NewServer(cfg *config.Config, handler *httprouter.Router, log *logger.Logger) *Server {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:63342"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(handler)

	return &Server{
		srv: &http.Server{
			Handler:      handlerWithCORS,
			WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
			Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		},
		log:     log,
		cfg:     cfg,
		handler: handler,
	}
}

func (s *Server) Run(dbConn *pgx.Conn) error {

	reqTimeout := s.cfg.PostgreSQL.RequestTimeout

	userStorage := user.NewStorage(dbConn, reqTimeout)
	userService := user.NewService(userStorage, *s.log)
	userHandler := user.NewHandler(*s.log, userService)
	userHandler.Register(s.handler)
	s.log.Info("initialized user routes")

	authStorage := user.NewStorage(dbConn, reqTimeout)
	authService := auth.NewService(authStorage, *s.log, *s.cfg)
	authHandler := auth.NewHandler(*s.log, authService, *s.cfg)
	authHandler.Register(s.handler)
	s.log.Info("initialized auth routes")

	appealStorage := appeal.NewStorage(dbConn, reqTimeout)
	appealService := appeal.NewService(appealStorage, *s.log)
	appealHandler := appeal.NewHandler(*s.log, appealService, *s.cfg, authService)
	appealHandler.Register(s.handler)
	s.log.Info("initialized appeal routes")

	fs := http.FileServer(http.Dir("public"))
	s.handler.Handler(http.MethodGet, "/", fs)

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
