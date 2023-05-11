package repo

type userService struct {
	repo UserRepo
}

type UserRepo interface {
	ShowSignup() error
	Signup() error
	ShowLogin() error
	Login() error
	Logout() error
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
