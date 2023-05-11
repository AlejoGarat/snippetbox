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

func (s *userService) Signup() error {
	return nil
}

func (s *userService) ShowLogin() error {
	return nil
}

func (s *userService) Login() error {
	return nil
}

func (s *userService) Logout() error {
	return nil
}
