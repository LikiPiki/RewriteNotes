package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/likipiki/VueGoNotes/server/crypt"
)

// UserService represents a PostgreSQL implementation of UserService.
type UserService struct {
	db *sql.DB
}

type User struct {
	Id uint `json:"id"`

	Password string `json:"password"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}
type Users []User

func NewUserService(db *sql.DB) UserService {
	return UserService{
		db: db,
	}
}

func (s UserService) GetAll() (Users, error) {

	rows, err := s.db.Query(
		"SELECT id, username, password, is_admin FROM users",
	)

	if err != nil {
		return nil, err
	}

	var users Users
	var user User

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
	row := s.db.QueryRow(
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
		user := User{
			Username: "admin",
			Password: pass,
			IsAdmin:  true,
		}
		_, err = s.db.Query(
			"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)",
			&user.Username, &user.Password, &user.IsAdmin,
		)
		if err != nil {
			log.Println("error creating superuser", err)
		}
	}
}

func (s UserService) Get(username string) (User, error) {

	var user User
	err := s.db.QueryRow(
		"SELECT id, username, password, is_admin FROM users WHERE username = $1",
		username,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil

}

func (s UserService) GetById(id interface{}) (User, error) {

	var user User
	err := s.db.QueryRow(
		"SELECT id, username, password, is_admin FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil

}

func (s UserService) Create(user User) error {
	cryptPassword, err := crypt.CryptPassword(user.Password)

	if err != nil {
		return err
	}
	_, err = s.db.Query(
		"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)",
		user.Username, cryptPassword, user.IsAdmin,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) Delete(id string) error {
	_, err := s.db.Query(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
