package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/likipiki/RewriteNotes/app/postgres"
	"github.com/likipiki/RewriteNotes/app/routes"
)

func main() {
	connStr := "password='postgres' dbname=notes sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	userController := postgres.UserService{
		DB: db,
	}
	noteController := postgres.NoteService{
		DB: db,
	}

	r := chi.NewRouter()

	userRouter := routes.UserHandlers{
		Controller: userController,
	}.Router()
	noteRouter := routes.NoteHandlers{
		Controller: noteController,
	}.Router()

	r.Mount("/user", userRouter)
	r.Mount("/note", noteRouter)

	fmt.Println("Listening on port :3000")
	http.ListenAndServe(":3000", r)
}
