package appeal

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"Interior_Visualization_Shop/app/internal/auth"
	"Interior_Visualization_Shop/app/internal/handler"
	"Interior_Visualization_Shop/app/internal/mail"
	"Interior_Visualization_Shop/app/internal/response"
	"Interior_Visualization_Shop/app/pkg/config"
	"Interior_Visualization_Shop/app/pkg/logger"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	appealURL = "/protected/appeal"
)

/// Структура Handler представляющая собой обработчик объекта appealService для обращений \\\

type Handler struct {
	log           logger.Logger
	appealService Service
	cfg           config.Config
	authService   auth.Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(log logger.Logger, appealService Service, cfg config.Config, authService auth.Service) handler.Hand {
	return &Handler{
		log:           log,
		appealService: appealService,
		cfg:           cfg,
		authService:   authService,
	}
}

/// Структура Register регистрирует новые запросы для обращений \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, appealURL, h.AuthMiddleware(h.CreateAppeal))
}

/// Структура AuthMiddleware проверят авторизирован ли пользователь в системе \\\

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("HANDLER: CHECK AUTH")

		// Извлекаем JWT-токен из заголовка запроса
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.ErrorAuth(w, "empty auth header", "")
			return
		}

		// Извлекаем строку токена из заголовка
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.ErrorAuth(w, "invalid auth header", "")
			return
		}

		// Передаем токен
		token, err := h.authService.ParseToken(tokenString)
		if err != nil {
			response.ErrorAuth(w, err.Error(), "")
			return
		}
		r.Header.Set("email", token)
		next(w, r)
	}
}

/// Вызов функции CreateAppeal для обработки запроса на создание обращения \\\

func (h *Handler) CreateAppeal(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: CREATE APPEAL")

	var input CreateAppealDTO

	/// Чтение данных типа form-data входящего запроса r. Email \\\
	input.Email = strings.TrimSpace(r.FormValue("email"))
	if input.Email == "" {
		response.BadRequest(w, "empty email", "")
		return
	}

	/// Чтение данных типа form-data входящего запроса r. Phonenumber \\\
	input.PhoneNumber = strings.TrimSpace(r.FormValue("phonenumber"))
	if input.PhoneNumber == "" {
		response.BadRequest(w, "empty phonenumber", "")
		return
	}
	/// Чтение данных типа form-data входящего запроса r. Nickname \\\
	input.Nickname = strings.TrimSpace(r.FormValue("nickname"))
	if input.Nickname == "" {
		response.BadRequest(w, "empty nickname", "")
		return
	}
	/// Чтение данных типа form-data входящего запроса r. Subject \\\
	subject := strings.TrimSpace(r.FormValue("subject"))
	if subject == "" {
		subject = "Feedback form"
	}
	input.Subject = &subject
	/// Чтение данных типа form-data входящего запроса r. Message \\\
	input.Message = strings.TrimSpace(r.FormValue("message"))
	if input.Message == "" {
		response.BadRequest(w, "empty message", "")
		return
	}

	/// Принимает файл входящего запроса r. Document \\\
	var docPath string
	file, header, err := r.FormFile("document")
	if err != nil {
		docPath = "without a file"
	} else {
		defer file.Close()
		/// Создаем путь и новый файл для сохранения документов \\\
		docPath = "./appealdocuments/" + input.Email + header.Filename
		out, err := os.Create(docPath)
		if err != nil {
			response.InternalError(w, fmt.Sprintf("error saving document: %v", err), "")
			return
		}
		defer out.Close()
		/// Копируем загруженный файл в новый созданный файл на севрере по созданному путю \\\
		_, err = io.Copy(out, file)
		if err != nil {
			response.InternalError(w, fmt.Sprintf("error coping document: %v", err), "")
			return
		}
	}
	/// Считываем данные поты отправителя из env файла \\\
	from := h.cfg.MAIL.MailAddress
	password := h.cfg.MAIL.MailPassword

	/// Формируем и отправляем письмо пользователю\\\
	h.log.Info("HANDLER: SENDING MESSAGE")
	err = mail.SendAppealEmail(from, password, input.Email, input.Nickname, *input.Subject)
	if err != nil {
		h.log.Errorf("failed to send message: %v", err)
	}

	a := CreateAppealDTO{
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Nickname:    input.Nickname,
		Subject:     input.Subject,
		Message:     input.Message,
		Document:    &docPath,
	}
	h.log.Printf("Input: %+v\n", &a)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру a \\\
	appeal, err := h.appealService.Create(r.Context(), &a)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedEmail) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}

	h.log.Info("APPEAL CREATED")
	response.JSON(w, http.StatusCreated, appeal)
}
