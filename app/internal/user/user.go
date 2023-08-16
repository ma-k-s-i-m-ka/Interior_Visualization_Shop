package user

import (
	"golang.org/x/crypto/bcrypt"
)

/// Структура для создания пользователей \\\

type User struct {
	ID       int64  `json:"id" example:"1567"`
	Email    string `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name     string `json:"name" example:"Maksim"`
	Surname  string `json:"surname" example:"Petrov"`
	Password string `json:"password"`
}

type CreateUserDTO struct {
	Email    string `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name     string `json:"name" example:"Maksim"`
	Surname  string `json:"surname" example:"Petrov"`
	Password string `json:"password" example:"sfdsg"`
}

/// Хэширование паролей \\\

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

/// Проверка введенного пароля на соответсвие паролю пользователя в БД \\\

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
