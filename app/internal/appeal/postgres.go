package appeal

import (
	"Interior_Visualization_Shop/app/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

var _ Storage = &AppealStorage{}

/// Структура AppealStorage содержащая поля для работы с БД \\\

type AppealStorage struct {
	log            logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр AppealStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &AppealStorage{
		log:            logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности AppealStorage создает записи обращений в БД \\\

func (d *AppealStorage) Create(appeal *Appeal) (*Appeal, error) {
	d.log.Info("POSTGRES: CREATE APPEAL")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO appeal (email, phone_number, nickname, subject, message, document)
			 VALUES($1,$2,$3,$4,$5,$6) 
			 RETURNING id`,
		appeal.Email, appeal.PhoneNumber, appeal.Nickname, appeal.Subject, appeal.Message, appeal.Document)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&appeal.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create appeal query: %v", err)
		return nil, err
	}
	return appeal, nil
}
