package repo

type userService struct {
	repo UserRepo
}

type UserRepo interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Exists(int) (bool, error)
}

func NewUserService(repo UserRepo) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) ShowSignup() error {
	return nil
}

func (s *userService) Signup(name string, email string, password string) error {
	return s.repo.Insert(name, email, password)
}

func (s *userService) ShowLogin() error {
	return nil
}

func (s *userService) Login(email string, password string) (int, error) {
	return s.repo.Authenticate(email, password)
}

func (s *userService) Logout() error {
	return nil
}
