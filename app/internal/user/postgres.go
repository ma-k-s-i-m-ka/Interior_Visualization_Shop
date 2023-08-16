package user

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"Interior_Visualization_Shop/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

var _ Storage = &UserStorage{}

/// Структура UserStorage содержащая поля для работы с БД \\\

type UserStorage struct {
	log            logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр UserStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &UserStorage{
		log:            logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности UserStorage создает записи пользователя в БД \\\

func (d *UserStorage) Create(user *User) (*User, error) {
	d.log.Info("POSTGRES: CREATE USER")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO users (email, name, surname, password)
			 VALUES($1,$2,$3,$4) 
			 RETURNING id`,
		user.Email, user.Name, user.Surname, user.Password)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&user.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create user query: %v", err)
		return nil, err
	}
	return user, nil
}

/// Функция FindByEmail для сущности UserStorage получает записи пациентов из БД по адресу электронной почты \\\

func (d *UserStorage) FindByEmail(email string) (*User, error) {
	d.log.Info("POSTGRES: GET USER BY EMAIL")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM users
			 WHERE email = $1`, email)
	user := &User{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&user.ID, &user.Email, &user.Name, &user.Surname, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		err = fmt.Errorf("failed to execute find user by email query: %v", err)
		return nil, err
	}
	return user, nil
}

/// Функция FindById для сущности UserStorage получает записи пользователя из БД по id \\\

func (d *UserStorage) FindById(id int64) (*User, error) {
	d.log.Info("POSTGRES: GET USER BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM users
			 WHERE id = $1`, id)
	user := &User{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&user.ID, &user.Email, &user.Name, &user.Surname, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find user by id query: %v", err)
		return nil, err
	}
	return user, nil
}

/// Функция Delete для сущности UserStorage удаляет записи о пользователях из БД \\\

func (d *UserStorage) Delete(id int64) error {
	d.log.Info("POSTGRES: DELETE USER")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
