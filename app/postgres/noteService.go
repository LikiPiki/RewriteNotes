package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/likipiki/RewriteNotes/app"
)

// NoteServiced represents a PostgreSQL implementation of app.NoteServiced.
type NoteService struct {
	DB *sql.DB
}

func (s NoteService) Get(id string) (app.Note, error) {
	var note app.Note
	row := s.DB.QueryRow(
		"SELECT id, user_id, title, content, created_at FROM notes WHERE id = $1",
		&id,
	)

	err := row.Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)
	if err != nil {
		return app.Note{}, err
	}

	return note, nil
}

func (s NoteService) GetAll(id string) (app.Notes, error) {
	rows, err := s.DB.Query(
		"SELECT id, user_id, title, content, created_at FROM notes WHERE user_id = $1",
		&id,
	)

	if err != nil {
		return nil, err
	}
	var note app.Note
	notes := make(app.Notes, 0)

	for rows.Next() {
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (s NoteService) Create(newNote app.Note) (bool, error) {
	_, err := s.DB.Query(
		"INSERT INTO notes(user_id, title, content, created_at) VALUES ($1, $2, $3, $4)",
		&newNote.UserID,
		&newNote.Title,
		&newNote.Content,
		&newNote.CreatedAt,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s NoteService) Delete(id string) error {
	_, err := s.DB.Query(
		"DELETE FROM notes WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
