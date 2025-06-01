package student

import (
	"context"
	"database/sql"
	"errors"
	studentv1 "github.com/bwjson/grpc_proto/gen/go/student"
	"github.com/bwjson/grpc_server/internal/db/postgres"
	"github.com/bwjson/grpc_server/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Student interface {
}

type serverAPI struct {
	studentv1.UnimplementedStudentsServer
	studentRepo postgres.StudentRepo
}

func Register(gRPC *grpc.Server, repo postgres.StudentRepo) {
	studentv1.RegisterStudentsServer(gRPC, &serverAPI{studentRepo: repo})
}

func (s *serverAPI) GetAll(ctx context.Context, request *studentv1.GetAllStudentsRequest) (*studentv1.GetAllStudentsResponse, error) {
	students, err := s.studentRepo.GetAll(ctx, request.SortBy, request.SortOrder)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &studentv1.GetAllStudentsResponse{}

	for _, student := range students {
		resp.Students = append(resp.Students, &studentv1.Student{
			Id:    student.Id,
			Name:  student.Name,
			Email: student.Email,
			Gpa:   student.Gpa,
		})
	}

	return resp, nil
}

func (s *serverAPI) GetById(ctx context.Context, request *studentv1.GetByIdStudentRequest) (*studentv1.GetByIdStudentResponse, error) {
	if request.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument")
	}

	student, err := s.studentRepo.GetById(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "student not found")
	}

	resp := &studentv1.GetByIdStudentResponse{
		Student: &studentv1.Student{
			Id:    student.Id,
			Name:  student.Name,
			Email: student.Email,
			Gpa:   student.Gpa,
		},
	}

	return resp, nil
}

func (s *serverAPI) Create(ctx context.Context, request *studentv1.CreateStudentRequest) (*studentv1.CreateStudentResponse, error) {
	student := models.Student{
		Name:  request.Student.Name,
		Gpa:   request.Student.Gpa,
		Email: request.Student.Email,
	}

	id, err := s.studentRepo.Create(ctx, student)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &studentv1.CreateStudentResponse{
		Id: int32(id),
	}

	return resp, nil
}

func (s *serverAPI) Update(ctx context.Context, request *studentv1.UpdateStudentRequest) (*studentv1.UpdateStudentResponse, error) {
	student := models.Student{
		Id:    request.Student.Id,
		Name:  request.Student.Name,
		Gpa:   request.Student.Gpa,
		Email: request.Student.Email,
	}

	id, err := s.studentRepo.Update(ctx, student)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "student not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &studentv1.UpdateStudentResponse{
		Id: int32(id),
	}

	return resp, nil
}

func (s *serverAPI) Delete(ctx context.Context, request *studentv1.DeleteStudentRequest) (*studentv1.DeleteStudentResponse, error) {
	id, err := s.studentRepo.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "student not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &studentv1.DeleteStudentResponse{
		Id: int32(id),
	}

	return resp, nil
}
