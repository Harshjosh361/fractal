package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Book represents a book entity with a title and an ISBN number.
type Book struct {
	Title string `json:"title"`
	ISBN  int    `json:"isbn"`
}

// Add handles adding a book to the database.
func Add(ctx context.Context, db *sql.DB, req *http.Request) (interface{}, error) {
	var book Book
	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	query := "INSERT INTO books (title, isbn) VALUES (?, ?)"
	result, err := db.ExecContext(ctx, query, book.Title, book.ISBN)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return lastInsertID, nil
}

func TestAdd(t *testing.T) {
	greenTick := "\033[32m✔\033[0m"
	redCross := "\033[31m✘\033[0m"

	logTestStatus := func(description string, err error) {
		if err == nil {
			t.Logf("%s %s", greenTick, description)
		} else {
			t.Logf("%s %s: %v", redCross, description, err)
		}
	}

	type gofrResponse struct {
		result interface{}
		err    error
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s Unexpected error while creating mock: %v", redCross, err)
	}
	defer db.Close()

	ctx := context.Background()

	tests := []struct {
		name             string
		requestBody      string
		mockExpect       func()
		expectedResponse gofrResponse
	}{
		{
			name:        "Error while Binding",
			requestBody: `title":"Book Title","isbn":12345}`, // Invalid JSON
			mockExpect:  func() {},                           // No DB interaction expected
			expectedResponse: gofrResponse{
				nil,
				errors.New("invalid request body"),
			},
		},
		{
			name:        "Successful Insertion",
			requestBody: `{"title":"Book Title","isbn":12345}`,
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Book Title", 12345).
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			expectedResponse: gofrResponse{
				int64(12),
				nil,
			},
		},
		{
			name:        "Error on Insertion",
			requestBody: `{"title":"Book Title","isbn":12345}`,
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Book Title", 12345).
					WillReturnError(sql.ErrConnDone)
			},
			expectedResponse: gofrResponse{
				nil,
				sql.ErrConnDone,
			},
		},
		{
			name:        "Error while fetching LastInsertId",
			requestBody: `{"title":"Book Title","isbn":12345}`,
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Book Title", 12345).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("mocked result error")))
			},
			expectedResponse: gofrResponse{
				nil,
				errors.New("mocked result error"),
			},
		},
	}

	t.Logf("=== RUN   sql_test.go")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			req := httptest.NewRequest(
				http.MethodPost,
				"/book",
				bytes.NewBuffer([]byte(tt.requestBody)),
			)
			req.Header.Set("Content-Type", "application/json")

			val, err := Add(ctx, db, req)

			response := gofrResponse{val, err}

			if assert.Equal(t, tt.expectedResponse, response, "Test failed: %s", tt.name) {
				logTestStatus(tt.name+" - Validate Response", nil)
			} else {
				logTestStatus(tt.name+" - Validate Response", assert.AnError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				logTestStatus("Mock expectations met", err)
			} else {
				logTestStatus("Mock expectations met", nil)
			}
		})
	}
}
