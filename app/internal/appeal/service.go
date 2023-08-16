package appeal

import (
	"Interior_Visualization_Shop/app/pkg/logger"
	"context"
)

/// Интерфейс Service реализизирующий service и методы для работы с обращениями \\\

type Service interface {
	Create(ctx context.Context, appeal *CreateAppealDTO) (*Appeal, error)
}

/// Структура  service реализизирующая инфтерфейс Service обращений \\\

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

/// Функция Create создает обращение через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateAppealDTO) (*Appeal, error) {
	s.log.Info("SERVICE: CREATE APPEAL")

	/// Создание структуры a на основе полученных данных \\\
	a := Appeal{
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Nickname:    input.Nickname,
		Subject:     input.Subject,
		Message:     input.Message,
		Document:    input.Document,
	}

	/// Вызов функции Create в хранилище записей \\\
	appeal, err := s.storage.Create(&a)
	if err != nil {
		return nil, err
	}
	return appeal, nil
}
