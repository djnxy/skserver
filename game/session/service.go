package session

type Service interface {
	Login(id string) error
	ListUsers() ([]string, error)
	Logout(id string) error
}

type service struct {
	users Repository
}

func (s *service) Login(id string) error {
	return s.users.Store(id)
}
func (s *service) ListUsers() ([]string, error) {
	return s.users.FindAll(), nil
}
func (s *service) Logout(id string) error {
	return s.users.Remove(id)
}

func NewService(users Repository) Service {
	return &service{users: users}
}

type Repository interface {
	Store(id string) error
	Remove(id string) error
	Find(id string) (string, error)
	FindAll() []string
}
