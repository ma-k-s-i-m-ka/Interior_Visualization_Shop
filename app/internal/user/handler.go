package user

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"Interior_Visualization_Shop/app/internal/handler"
	"Interior_Visualization_Shop/app/internal/response"
	"Interior_Visualization_Shop/app/pkg/logger"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	usersURL       = "/users"
	userByEmailURL = "/users/email"
	userURL        = "/users/profile/:id"
)

/// Структура Handler представляющая собой обработчик объекта userService для пользователей \\\

type Handler struct {
	log         logger.Logger
	userService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(log logger.Logger, userService Service) handler.Hand {
	return &Handler{
		log:         log,
		userService: userService,
	}
}

/// Структура Register регистрирует новые запросы для пользователей \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userByEmailURL, h.GetUserByEmail)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
	router.HandlerFunc(http.MethodDelete, userURL, h.DeleteUser)
	router.HandlerFunc(http.MethodGet, userURL, h.GetUserById)
}

/// Функция GetUserById получает пользователя по его id \\\

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GET USER BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)

	h.log.Printf("Input: %+v\n", &id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей id пациента \\\
	user, err := h.userService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "")
		return
	}
	h.log.Info("GOT USER BY ID")
	response.JSON(w, http.StatusOK, user)
}

/// Функция GetUserByEmail получает пользователя по его email \\\

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GET USER BY EMAIL")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр email из URL \\\

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	email := r.URL.Query().Get("email")
	h.log.Printf("Input: %+v\n", email)
	if email == "" {
		response.BadRequest(w, "empty email", "")
		return
	}

	/// Вызов функции GetByEmail передавая ей email \\\
	user, err := h.userService.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.log.Info("GOT USER BY EMAIL")
	response.JSON(w, http.StatusOK, user)
}

/// Функция CreateUser создает пользователя по полученным данным из input \\\

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: CREATE USER")

	var input CreateUserDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.log.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	user, err := h.userService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedEmail) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}
	h.log.Info("USER CREATED")
	response.JSON(w, http.StatusCreated, user)
}

/// Функция DeleteUser удаляет пользователя по его id \\\

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: DELETE USER")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.log.Printf("Input: %+v\n", id)
	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.userService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.log.Info("USER DELETED")
	response.JSON(w, http.StatusOK, "USER DELETED")
}
