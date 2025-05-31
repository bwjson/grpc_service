package student

import (
	"context"
	studentv1 "github.com/bwjson/grpc_proto/gen/go/student"
	"github.com/bwjson/grpc_server/internal/db/postgres"
	"google.golang.org/grpc"
)

type Student interface {
}

type serverAPI struct {
	studentv1.UnimplementedStudentsServer
	studentRepo postgres.StudentRepo
}

// REALIZATION OF STUDENTS REPO
func Register(gRPC *grpc.Server, repo postgres.StudentRepo) {
	studentv1.RegisterStudentsServer(gRPC, &serverAPI{studentRepo: repo})
}

func (s *serverAPI) GetAll(ctx context.Context, request *studentv1.GetAllStudentsRequest) (*studentv1.GetAllStudentsResponse, error) {
	students, err := s.studentRepo.GetAll(ctx)
	if err != nil {
		return nil, err
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
	panic("implement me")
}

func (s *serverAPI) Create(ctx context.Context, request *studentv1.CreateStudentRequest) (*studentv1.CreateStudentResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Update(ctx context.Context, request *studentv1.UpdateStudentRequest) (*studentv1.UpdateStudentResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Delete(ctx context.Context, request *studentv1.DeleteStudentRequest) (*studentv1.DeleteStudentResponse, error) {
	panic("implement me")
}
