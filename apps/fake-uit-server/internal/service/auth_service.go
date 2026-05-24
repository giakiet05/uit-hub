package service

import (
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
)

func (s *Service) Login(req dto.LoginRequest) (dto.SuccessResponse, *apperror.AppError) {
	if isBlank(req.StudentID) || isBlank(req.Password) {
		return dto.SuccessResponse{}, newBadRequest("student_id and password are required")
	}
	student, ok, err := s.students.FindByID(req.StudentID)
	if err != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if !ok || student.Password != req.Password {
		return dto.SuccessResponse{}, newUnauthorized("invalid credentials")
	}

	token := "mock-" + student.ID
	return success(map[string]string{"token": token}), nil
}
