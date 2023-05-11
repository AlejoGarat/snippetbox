package repo

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/AlejoGarat/snippetbox/internal/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = r.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (r *userRepo) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (r *userRepo) Exists(id int) (bool, error) {
	return false, nil
}
