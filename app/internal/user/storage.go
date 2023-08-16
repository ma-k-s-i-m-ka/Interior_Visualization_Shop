package user

type Storage interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindById(id int64) (*User, error)
	Delete(id int64) error
}
