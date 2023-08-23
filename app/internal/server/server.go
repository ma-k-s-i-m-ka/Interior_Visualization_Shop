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
	"github.com/pkg/browser"
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

	/// список разрешенных источников (доменов), заголовков и HTTP-методов, которые разрешены для выполнения на сервере \\\
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

/// Функция инициализирующая хранище storage, сервисы services и обработчики handler \\\
/// Запускает сервер и начинает обрабатывать входящие HTTP запросы \\\

func (s *Server) Run(dbConn *pgx.Conn) error {

	reqTimeout := s.cfg.PostgreSQL.RequestTimeout

	/// Инициализация хранилища userStorage, создание объекта сервиса userService, создание обработчика userHandler для пользователей \\\
	/// Тот же принцип работы для остальных route \\\

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

	/// создание файлового сервера для статических файлов, которые находятся в директории "public" \\\
	fs := http.FileServer(http.Dir("public"))
	s.handler.Handler(http.MethodGet, "/", fs)

	/// обработка запросов для пути "/style.css", обслуживаеться файловым сервером fs \\\
	s.handler.Handler(http.MethodGet, "/style.css", fs)
	s.handler.Handler(http.MethodGet, "/index.html", fs)
	s.handler.Handler(http.MethodGet, "/appeal.html", fs)
	s.handler.Handler(http.MethodGet, "/contacts.html", fs)
	s.handler.Handler(http.MethodGet, "/portfolio.html", fs)
	s.handler.Handler(http.MethodGet, "/service.html", fs)
	s.handler.Handler(http.MethodGet, "/sign-in.html", fs)
	s.handler.Handler(http.MethodGet, "/sign-up.html", fs)
	s.handler.Handler(http.MethodGet, "/1.jpg", fs)
	s.handler.Handler(http.MethodGet, "/2.jpg", fs)
	s.handler.Handler(http.MethodGet, "/3.jpg", fs)
	s.handler.Handler(http.MethodGet, "/4.jpg", fs)
	s.handler.Handler(http.MethodGet, "/5.jpg", fs)
	s.handler.Handler(http.MethodGet, "/6.jpg", fs)
	s.handler.Handler(http.MethodGet, "/7.jpg", fs)
	s.handler.Handler(http.MethodGet, "/8.jpg", fs)
	s.handler.Handler(http.MethodGet, "/9.jpg", fs)
	s.handler.Handler(http.MethodGet, "/10.jpg", fs)
	s.handler.Handler(http.MethodGet, "/tel.png", fs)
	s.handler.Handler(http.MethodGet, "/vk.jpg", fs)
	s.handler.Handler(http.MethodGet, "/wp.jpg", fs)
	s.handler.Handler(http.MethodGet, "/gm.jpg", fs)
	s.handler.Handler(http.MethodGet, "/clients-bg.jpg", fs)

	/// открытие веб-страницы в браузере \\\
	err := browser.OpenURL("http://" + s.srv.Addr + "/")
	if err != nil {
		return err
	}

	return s.srv.ListenAndServe()
}

/// Метоод Shutdown структуры Server. Функция для завершения работы сервера \\\

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
