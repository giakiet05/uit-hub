package service

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

func (s *Service) GetStudentByIDForAuth(id string) (*model.Student, bool) {
	st, err := s.getStudent(id)
	if err != nil {
		return nil, false
	}
	return st, true
}
