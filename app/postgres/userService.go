package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/likipiki/RewriteNotes/app"
	"github.com/likipiki/VueGoNotes/server/crypt"
)

// UserService represents a PostgreSQL implementation of app.UserService.
type UserService struct {
	DB *sql.DB
}

func (service UserService) GetAll() (app.Users, error) {

	rows, err := service.DB.Query(
		"SELECT id, username, password, is_admin FROM users",
	)

	if err != nil {
		return nil, err
	}

	var users app.Users
	var user app.User

	for rows.Next() {
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.IsAdmin,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	defer rows.Close()

	return users, nil
}

func (service UserService) Get(username string) (app.User, error) {

	var user app.User
	err := service.DB.QueryRow(
		"SELECT id, username, password, is_admin FROM users WHERE username = $1",
		username,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
	)

	if err != nil {
		return app.User{}, err
	}

	return user, nil

}

func (service UserService) Create(user app.User) error {
	cryptPassword, err := crypt.CryptPassword(user.Password)

	if err != nil {
		return err
	}
	_, err = service.DB.Query(
		"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)",
		user.Username, cryptPassword, user.IsAdmin,
	)
	if err != nil {
		return err
	}
	return nil
}

func (service UserService) Delete(id string) error {
	_, err := service.DB.Query(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
