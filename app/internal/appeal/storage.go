package appeal

type Storage interface {
	Create(user *Appeal) (*Appeal, error)
}
