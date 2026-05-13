package service

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

type StudentService interface {
	GetStudents() []model.Student
	GetStudentByID(id string) (*model.Student, bool)
}

type studentService struct {
	students []model.Student
}

func NewStudentService() StudentService {
	return &studentService{
		students: []model.Student{
			{
				ID:    "1",
				Name:  "Duc",
				Email: "duc@example.com",
				Phone: "0900000001",
				Major: "Software Engineering",
			},
			{
				ID:    "2",
				Name:  "An",
				Email: "an@example.com",
				Phone: "0900000002",
				Major: "Information Systems",
			},
		},
	}
}

func (s *studentService) GetStudents() []model.Student {
	return s.students
}

func (s *studentService) GetStudentByID(id string) (*model.Student, bool) {
	for _, student := range s.students {
		if student.ID == id {
			return &student, true
		}
	}

	return nil, false
}
