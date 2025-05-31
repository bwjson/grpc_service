package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bwjson/grpc_server/internal/models"
	_ "github.com/lib/pq"
)

type StudentRepo interface {
	GetAll(ctx context.Context) ([]*models.Student, error)
	//GetById(ctx context.Context, id int) (*models.Student, error)
	//Create(ctx context.Context, user models.Student) (int, error)
	//Update(ctx context.Context, user models.Student) (int, error)
	//Delete(ctx context.Context, id int) (int, error)
}

type Database struct {
	db *sql.DB
}

func New(host, port, user, dbname, pass string) (*Database, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, pass)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to postgres: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// REALIZATION OF STUDENTS REPO
func (s *Database) GetAll(ctx context.Context) ([]*models.Student, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Gpa)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
