package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/likipiki/RewriteNotes/app"
	"github.com/likipiki/VueGoNotes/server/crypt"
)

// UserService represents a PostgreSQL implementation of app.UserService.
type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return UserService{
		DB: db,
	}
}

func (s UserService) GetAll() (app.Users, error) {

	rows, err := s.DB.Query(
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

// create admin user if not exists
func (s UserService) Install() {
	row := s.DB.QueryRow(
		"SELECT username FROM users WHERE username = $1",
		"admin",
	)
	var name string
	err := row.Scan(&name)
	if name == "" || err != nil {
		err = row.Scan()
		pass, err := crypt.CryptPassword("admin")
		if err != nil {
			log.Println("error crypting admin password")
		}
		user := app.User{
			Username: "admin",
			Password: pass,
			IsAdmin:  true,
		}
		_, err = s.DB.Query(
			"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)",
			&user.Username, &user.Password, &user.IsAdmin,
		)
		if err != nil {
			log.Println("error creating superuser", err)
		}
	}
}

func (s UserService) Get(username string) (app.User, error) {

	var user app.User
	err := s.DB.QueryRow(
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

func (s UserService) Create(user app.User) error {
	cryptPassword, err := crypt.CryptPassword(user.Password)

	if err != nil {
		return err
	}
	_, err = s.DB.Query(
		"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)",
		user.Username, cryptPassword, user.IsAdmin,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) Delete(id string) error {
	_, err := s.DB.Query(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}