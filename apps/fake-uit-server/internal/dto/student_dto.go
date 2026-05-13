package dto

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

type StudentResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Major string `json:"major"`
}

func NewStudentResponse(student model.Student) StudentResponse {
	return StudentResponse{
		ID:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Phone: student.Phone,
		Major: student.Major,
	}
}

func NewStudentResponses(students []model.Student) []StudentResponse {
	responses := make([]StudentResponse, 0, len(students))
	for _, student := range students {
		responses = append(responses, NewStudentResponse(student))
	}

	return responses
}
