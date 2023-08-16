package auth

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"Interior_Visualization_Shop/app/internal/user"
	"Interior_Visualization_Shop/app/pkg/config"
	"Interior_Visualization_Shop/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/// Интерфейс Service реализизирующий service и методы для обработки логики аутентификации и регистрации пользователей \\\

type Service interface {
	AuthByEmail(ctx context.Context, user *AuthByEmail) (*user.User, *AuthResponse, error)
	Register(ctx context.Context, user *Register) (*user.User, *RegisterResponse, error)
	CreateAccessToken(cfg *config.Config, user *user.User) (string, error)
	CreateRefreshToken(cfg *config.Config, user *user.User) (string, error)
	ParseToken(token string) (string, error)
}

/// Структура  service реализизирующая инфтерфейс Service пользователей \\\

type service struct {
	log     logger.Logger
	storage user.Storage
	cfg     config.Config
}

/// Структура NewService возвращает новый экземпляр Service инициализируя переданные в него аргументы \\\

func NewService(storage user.Storage, log logger.Logger, cfg config.Config) Service {
	return &service{
		log:     log,
		storage: storage,
		cfg:     cfg,
	}
}

// / Структура tokenClaims хранящая информацию о сессиях пользователей \\\
type tokenClaims struct {
	jwt.MapClaims
	Email string `json:"email"`
}

/// Функция AuthByEmail реализует аутентификацию пользователя по адресу электронной почты через интерфейс Service принимая входные данные input  \\\

func (s *service) AuthByEmail(ctx context.Context, input *AuthByEmail) (*user.User, *AuthResponse, error) {
	s.log.Info("SERVICE: AUTH USER BY EMAIL")

	/// Вызов функции FindByEmail в хранилище пользователей  \\\
	user, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, nil, err
		}
		s.log.Error("cannot find user by email:", err)
		return nil, nil, err
	}
	/// Проверка на соответствие введенного и захэшированного пароля в хранилище \\\
	if !user.CheckPassword(input.Password) {
		s.log.Error("incorrect password:", err)
		return nil, nil, err
	}

	/// Создание токенов доступа \\\
	accessToken, err := s.CreateAccessToken(&s.cfg, user)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := s.CreateRefreshToken(&s.cfg, user)
	if err != nil {
		return nil, nil, err
	}
	return user, &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

/// Функция Register реализует регистрацию пользователя через интерфейс Service принимая входные данные input  \\\

func (s *service) Register(ctx context.Context, input *Register) (*user.User, *RegisterResponse, error) {
	s.log.Info("SERVICE: REGISTER USER")

	/// Проверка на повтаряющийся адрес электронной почты \\\
	/// Вызов функции FindByEmail в хранилище пользователей  \\\
	checkEmail, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, nil, err
		}
	}

	if checkEmail != nil {
		return nil, nil, apperror.ErrRepeatedEmail
	}

	u := user.User{
		Email:    input.Email,
		Name:     input.Name,
		Surname:  input.Surname,
		Password: input.Password,
	}

	/// Хэширование полученного пароля \\\
	err = u.HashPassword()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot hash password")
	}

	/// Вызов функции Create в хранилище пользователей  \\\
	user, err := s.storage.Create(&u)

	/// Создание токенов доступа \\\
	accessToken, err := s.CreateAccessToken(&s.cfg, user)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := s.CreateRefreshToken(&s.cfg, user)
	if err != nil {
		return nil, nil, err
	}
	return user, &RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

/// Функция CreateAccessToken для создания токена доступа AccessToken \\\

func (s *service) CreateAccessToken(cfg *config.Config, user *user.User) (string, error) {
	s.log.Info("SERVICE: CREATE ACCESS TOKEN")
	metadata := AccessToken{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Surname: user.Surname,
	}

	/// Создание нового токена accessToken с указанными утверждениями claims и методом подписи SigningMethodHS256 \\\
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.MapClaims{
			"user": metadata,
			"exp":  time.Duration(cfg.JWT.AccessExpirationMinutes) * time.Minute,
		}, metadata.Email,
	})
	/// Токен подписывается с помощью секретного ключа AccessTokenSecretKey и преобразуется в строку с помощью метода SignedString \\\
	token, err := accessToken.SignedString([]byte(cfg.JWT.AccessTokenSecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

/// Функция CreateRefreshToken для создания токена обновления RefreshToken \\\

func (s *service) CreateRefreshToken(cfg *config.Config, user *user.User) (string, error) {
	s.log.Info("SERVICE: CREATE REFRESH TOKEN")

	/// Добавление id пользователя в мапу claims, и время истечения действия токена RefreshToken \\\
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Duration(cfg.JWT.RefreshExpirationDays) * time.Hour * 24,
	}
	/// Создание нового токена refreshToken с указанными утверждениями claims и методом подписи SigningMethodHS256 \\\
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	/// Токен подписывается с помощью секретного ключа RefreshTokenSecretKey и преобразуется в строку с помощью метода SignedString \\\
	token, err := refreshToken.SignedString([]byte(cfg.JWT.RefreshTokenSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

/// Функция ParseToken для передачи accessToken токена \\\

func (s *service) ParseToken(accessToken string) (string, error) {
	s.log.Info("HANDLER: PARSE TOKEN")

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid metod")
		}
		return []byte(s.cfg.JWT.AccessTokenSecretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("bad recived token")
	}

	claim, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("token claim bad type")
	}
	return claim.Email, nil
}
