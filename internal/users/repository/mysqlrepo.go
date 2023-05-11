package repo

import (
	"database/sql"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{
		DB: db,
	}
}

func (r *userRepo) ShowSignup() error {
	return nil
}

func (r *userRepo) Signup() error {
	return nil
}

func (r *userRepo) ShowLogin() error {
	return nil
}

func (r *userRepo) Login() error {
	return nil
}

func (r *userRepo) Logout() error {
	return nil
}
