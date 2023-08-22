package auth

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"Interior_Visualization_Shop/app/internal/handler"
	"Interior_Visualization_Shop/app/internal/mail"
	"Interior_Visualization_Shop/app/internal/response"
	"Interior_Visualization_Shop/app/pkg/config"
	"Interior_Visualization_Shop/app/pkg/logger"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	userAuthByEmailURL   = "/sign_in/mail"
	userRegisterURL      = "/sign_up"
	userRegisterCheckURL = "/sign_up/checkmail"
)

/// Структура Handler представляющая собой обработчик объекта authService для пользователей \\\

type Handler struct {
	log         logger.Logger
	authService Service
	cfg         config.Config
	CodeChan    chan string
	CodeMutex   sync.Mutex
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(log logger.Logger, authService Service, cfg config.Config) handler.Hand {
	return &Handler{
		log:         log,
		authService: authService,
		cfg:         cfg,
	}
}

/// Структура Register регистрирует новые запросы для авторизации \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, userAuthByEmailURL, h.GetUserByEmail)
	router.HandlerFunc(http.MethodPost, userRegisterURL, h.RegisterUser)
	router.HandlerFunc(http.MethodPost, userRegisterCheckURL, h.CheckMailCode)
}

/// Функция GetUserByEmail получает пользователя по его адресу электронной почты и паролю \\\

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: AUTH BY EMAIL")

	var input AuthByEmail
	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.log.Printf("Input: %+v\n", &input)
	/// Вызов функции AuthByEmail передавая ей полученные значения и ссылку на структуру input \\\
	user, jwt, err := h.authService.AuthByEmail(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}

	h.log.Info("AUTH BY EMAIL IS COMPLETED")
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
		"jwt":  jwt,
	})
}

/// Функция RegisterUser регистрирует пользователя \\\

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: REGISTER USER")

	var input Register
	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	/// Формируем 4-х значный код подтверждения \\\
	h.log.Printf("Input: %+v\n", &input)
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000
	codeStr := strconv.Itoa(code)

	/// Считываем данные поты отправителя из env файла \\\
	from := h.cfg.MAIL.MailAddress
	password := h.cfg.MAIL.MailPassword

	/// Формируем и отправляем письмо пользователю \\\
	err := mail.SendEmail(from, password, input.Email, input.Name, input.Surname, codeStr)
	if err != nil {
		h.log.Errorf("failed to send messeg: %v", err)
	}

	/// Ждем пока пользователь введет код \\\
	h.log.Info("HANDLER: WAITING FOR THE CODE")
	h.CodeChan = make(chan string)
	go h.CheckMailCode(w, r)

	/// Получаем введенный пользователем код из CheckMailCode \\\
	checkcode := <-h.CodeChan
	h.log.Info("HANDLER: CODE RECEIVED")
	h.log.Printf("Input: %+v\n", checkcode)

	/// Сравниваем код введенный пользователем с тем кодом который был отправлен пользователю на почту \\\
	if checkcode != codeStr {
		response.BadRequest(w, "the entered code is not correct", apperror.ErrInvalidMailCode.Error())
		return
	}

	/// Вызов функции Register передавая ей полученные значения и ссылку на структуру input \\\
	user, jwt, err := h.authService.Register(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedEmail) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}

	h.log.Info("REGISTER USER IS COMPLETED")
	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user": user,
		"jwt":  jwt,
	})
}

/// Функция CheckMailCode получает код подтверждения для RegisterUser \\\

func (h *Handler) CheckMailCode(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GETTING THE REGISTRATION CODE")

	var Check struct {
		Code string `json:"code"`
	}

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &Check); err != nil {
		//response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	/// передача введенного пользователем кода через канал CodeChan в RegisterUser \\\
	h.CodeChan <- Check.Code

	response.JSON(w, http.StatusOK, Check.Code)
}
