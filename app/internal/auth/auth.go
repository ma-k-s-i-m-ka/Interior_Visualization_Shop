package auth

import "golang.org/x/crypto/bcrypt"

/// Структура для авторизации и регистрации пользователей \\\

type AccessToken struct {
	ID      int64  `json:"id" example:"1567"`
	Email   string `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name    string `json:"name" example:"Maksim"`
	Surname string `json:"surname" example:"Petrov"`
}

type RefreshToken struct {
	ID int64 `json:"id" example:"1567"`
}

type AuthByEmail struct {
	Email    string `json:"email" example:"petrovmaksim1992@mail.ru"`
	Password string `json:"password" example:"abcdEFG"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Register struct {
	Email    string `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name     string `json:"name" example:"Maksim"`
	Surname  string `json:"surname" example:"Petrov"`
	Password string `json:"password" example:"sfdsg"`
}
type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *Register) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
