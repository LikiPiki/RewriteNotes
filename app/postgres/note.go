package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// NoteServiced represents a PostgreSQL implementation of NoteService.
type NoteService struct {
	db *sql.DB
}

func NewNoteService(db *sql.DB) NoteService {
	return NoteService{
		db: db,
	}
}

type Note struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`

	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Notes []Note

func (s NoteService) Get(id string) (Note, error) {
	var note Note
	row := s.db.QueryRow(
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
		return Note{}, err
	}

	return note, nil
}

func (s NoteService) GetAll(id string) (Notes, error) {
	rows, err := s.db.Query(
		"SELECT id, user_id, title, content, created_at FROM notes WHERE user_id = $1",
		&id,
	)

	if err != nil {
		return nil, err
	}
	var note Note
	notes := make(Notes, 0)

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

func (s NoteService) Create(newNote Note) (bool, error) {
	_, err := s.db.Query(
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
	_, err := s.db.Query(
		"DELETE FROM notes WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
