package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/likipiki/RewriteNotes/app/postgres"
)

func TestLoginHandler(t *testing.T) {
	connStr := "password='postgres' dbname=notes sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()

	if err != nil {
		t.Error(err)
	}
	userController := postgres.NewUserService(db)
	hh := NewUserHandlers(userController)

	data, err := json.Marshal(map[string]string{
		"username": "admin",
		"password": "admin",
	})
	if err != nil {
		t.Error("marshal user error", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:3000/user/login", bytes.NewBuffer(data))
	if err != nil {
		t.Error("request error", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hh.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	type Result struct {
		Status bool   `json:"status"`
		Token  string `json:"token"`
	}

	var result Result
	err = json.NewDecoder(rr.Body).Decode(&result)
	if err != nil {
		t.Error(err)
	}

	if !(result.Status == true) && (result.Token != "") {
		t.Error(
			"Unexpected status, or invalid token",
			fmt.Sprintf("%+v", result),
		)
	}
}
