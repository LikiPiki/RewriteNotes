package app

type User struct {
	Id uint `json:"id"`

	Password string `json:"password"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}
type Users []User

type UserService interface {
	Install()
}
