package service

import (
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetDeadlines(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	deadlines, repoErr := s.deadlines.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]model.Deadline, 0)
	for _, deadline := range deadlines {
		if deadline.StudentID == student.ID {
			result = append(result, deadline)
		}
	}

	return success(result), nil
}

func (s *Service) GetMaterials(studentID, courseID string) (dto.SuccessResponse, *apperror.AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	_, ok, repoErr := s.courses.FindByID(courseID)
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if !ok {
		return dto.SuccessResponse{}, newNotFound("course not found")
	}
	materials, repoErr := s.materials.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]model.Material, 0)
	for _, material := range materials {
		if material.CourseID == courseID {
			result = append(result, material)
		}
	}

	return success(result), nil
}

func (s *Service) GetAssignments(studentID, courseID string) (dto.SuccessResponse, *apperror.AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	_, ok, repoErr := s.courses.FindByID(courseID)
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if !ok {
		return dto.SuccessResponse{}, newNotFound("course not found")
	}
	assignments, repoErr := s.assignments.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]model.Assignment, 0)
	for _, assignment := range assignments {
		if assignment.CourseID == courseID {
			result = append(result, assignment)
		}
	}

	return success(result), nil
}

func (s *Service) SubmitAssignment(studentID, assignmentID string, req dto.AssignmentSubmissionRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	assignment, ok, repoErr := s.assignments.FindByID(assignmentID)
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if !ok {
		return dto.SuccessResponse{}, newNotFound("assignment not found")
	}

	submissionType := strings.ToUpper(strings.TrimSpace(req.SubmissionType))
	if submissionType == "" {
		return dto.SuccessResponse{}, newBadRequest("submission_type is required")
	}

	content := ""
	switch submissionType {
	case "TEXT":
		if isBlank(req.Text) {
			return dto.SuccessResponse{}, newBadRequest("text is required for TEXT submissions")
		}
		content = req.Text
	case "LINK":
		if isBlank(req.URL) {
			return dto.SuccessResponse{}, newBadRequest("url is required for LINK submissions")
		}
		if !isValidURL(req.URL) {
			return dto.SuccessResponse{}, newBadRequest("url is invalid")
		}
		content = req.URL
	case "FILE":
		if len(req.FileURLs) == 0 {
			return dto.SuccessResponse{}, newBadRequest("file_urls are required for FILE submissions")
		}
		content = strings.Join(req.FileURLs, ",")
	default:
		return dto.SuccessResponse{}, newBadRequest("submission_type is invalid")
	}
	_ = uuid.NewString()
	count, repoErr := s.submissions.Count()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if err := s.submissions.Create(model.Submission{
		ID:             uint(count + 1),
		AssignmentID:   assignment.ID,
		StudentID:      student.ID,
		SubmissionType: submissionType,
		Content:        content,
		CreatedAt:      s.now().Format(time.RFC3339),
	}); err != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	return success(nil), nil
}
