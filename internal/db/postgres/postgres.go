package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwjson/grpc_server/internal/models"
	_ "github.com/lib/pq"
	"strings"
)

type StudentRepo interface {
	GetAll(ctx context.Context, sortBy, sortOrder *string) ([]*models.Student, error)
	GetById(ctx context.Context, id int32) (*models.Student, error)
	Create(ctx context.Context, user models.Student) (int64, error)
	Update(ctx context.Context, user models.Student) (int64, error)
	Delete(ctx context.Context, id int32) (int32, error)
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

func (s *Database) GetAll(ctx context.Context, sortBy, sortOrder *string) ([]*models.Student, error) {
	var (
		sortByParam    = "id"
		sortOrderParam = "asc"
	)

	validSortBy := map[string]bool{
		"id":    true,
		"name":  true,
		"email": true,
		"gpa":   true,
	}

	validSortOrder := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if sortBy != nil {
		if validSortBy[strings.ToLower(*sortBy)] {
			sortByParam = *sortBy
		} else {
			return nil, errors.New("Invalid sort by parameter")
		}
	}

	if sortOrder != nil {
		if validSortOrder[strings.ToLower(*sortOrder)] {
			sortOrderParam = *sortOrder
		} else {
			return nil, errors.New("Invalid sort order parameter")
		}
	}

	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM students ORDER BY %s %s", sortByParam, sortOrderParam))
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

func (s *Database) GetById(ctx context.Context, id int32) (*models.Student, error) {
	var student models.Student

	row := s.db.QueryRowContext(ctx, "SELECT * FROM students WHERE id = $1", id)

	err := row.Scan(&student.Id, &student.Name, &student.Email, &student.Gpa)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (s *Database) Create(ctx context.Context, student models.Student) (int64, error) {
	var id int64

	err := s.db.QueryRowContext(ctx,
		"INSERT INTO students (name, email, gpa) VALUES ($1, $2, $3) RETURNING id",
		student.Name, student.Email, student.Gpa).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Database) Update(ctx context.Context, student models.Student) (int64, error) {
	var id int64

	err := s.db.QueryRowContext(ctx,
		"UPDATE students SET name = $1, email = $2, gpa = $3 WHERE id = $4 RETURNING id",
		student.Name, student.Email, student.Gpa, student.Id).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Database) Delete(ctx context.Context, id int32) (int32, error) {
	err := s.db.QueryRowContext(ctx, "DELETE FROM students WHERE id = $1 RETURNING id", id).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
