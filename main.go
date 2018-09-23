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

	userController := postgres.NewUserService(db)
	noteController := postgres.NewNoteService(db)

	// list Controllers to install default values
	postgres.Install(
		userController,
	)

	r := chi.NewRouter()

	defaultHandler := routes.NewDefaultHandlers()
	userHandlers := routes.NewUserHandlers(userController)
	noteHandlers := routes.NewNoteHandlers(noteController)

	r.Mount("/user", userHandlers.Router)
	r.Mount("/note", noteHandlers.Router)
	r.Mount("/", defaultHandler.Router)

	fmt.Println("Listening on port :3000")
	http.ListenAndServe(":3000", r)
}
