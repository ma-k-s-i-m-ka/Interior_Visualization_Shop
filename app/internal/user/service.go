package user

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"Interior_Visualization_Shop/app/pkg/logger"
	"context"
	"errors"
	"fmt"
)

/// Интерфейс Service реализизирующий service и методы для пользователей \\\

type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service пользователей \\\

type service struct {
	log     logger.Logger
	storage Storage
}

/// Структура NewService возвращает новый экземпляр Service инициализируя переданные в него аргументы \\\

func NewService(storage Storage, log logger.Logger) Service {
	return &service{
		log:     log,
		storage: storage,
	}
}

/// Функция Create создает пользователя через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateUserDTO) (*User, error) {
	s.log.Info("SERVICE: CREATE USER")

	/// Проверка на уникальность email \\\
	checkEmail, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}
	if checkEmail != nil {
		return nil, apperror.ErrRepeatedEmail
	}

	/// Создание структуры u на основе полученных данных \\\
	u := User{
		Email:    input.Email,
		Name:     input.Name,
		Surname:  input.Surname,
		Password: input.Password,
	}

	/// Хэширование пароля \\\
	err = u.HashPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot hash password")
	}

	/// Вызов функции Create в хранилище пользователей \\\
	user, err := s.storage.Create(&u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/// Функция GetByEmail осуществялет поиск пользователей через интерфейс Service принимая входные данные email пользователя \\\

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	s.log.Info("SERVICE: GET USER BY EMAIL")

	/// Вызов функции FindByEmail в хранилище пользователей \\\
	user, err := s.storage.FindByEmail(email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.log.Warn("cannot find user by email:", err)
		return nil, err
	}
	return user, nil
}

/// Функция GetById осуществялет поиск пользователей через интерфейс Service принимая входные данные id пользователя \\\

func (s *service) GetById(ctx context.Context, id int64) (*User, error) {
	s.log.Info("SERVICE: GET USER BY ID")

	/// Вызов функции FindById в хранилище пациентов \\\
	user, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.log.Warn("cannot find user by id:", err)
		return nil, err
	}
	return user, nil
}

/// Функция Delete удаляет пользователя через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.log.Info("SERVICE: DELETE USER")
	/// Вызов функции Delete в хранилище пациентов \\\
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.log.Warn("failed to delete user:", err)
		}
		return err
	}
	return nil
}
