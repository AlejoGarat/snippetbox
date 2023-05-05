package repo

import (
	"database/sql"
)

type HomeRepo struct {
	DB *sql.DB
}

func (hr *HomeRepo) Get() error {
	return nil
}
