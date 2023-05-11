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

func (r *userRepo) Insert(name, email, password string) error {
	return nil
}

func (r *userRepo) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (r *userRepo) Exists(id int) (bool, error) {
	return false, nil
}
