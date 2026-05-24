package service

import (
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
)

func (s *Service) ConfirmLetter(studentID string, req dto.ConfirmLetterRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	if isBlank(req.Language) || isBlank(req.Reason) || isBlank(req.RequestType) {
		return dto.SuccessResponse{}, newBadRequest("language, reason, request_type are required")
	}
	switch req.Language {
	case "VI", "EN":
	default:
		return dto.SuccessResponse{}, newBadRequest("language is invalid")
	}
	switch req.RequestType {
	case "NEW", "REISSUE":
	default:
		return dto.SuccessResponse{}, newBadRequest("request_type is invalid")
	}
	switch req.Reason {
	case "MILITARY_DEFERMENT", "DORM_EXTEND", "TAX_DEDUCTION_DOCS", "DEFENSE_EDU_REGISTRATION", "OTHER":
	default:
		return dto.SuccessResponse{}, newBadRequest("reason is invalid")
	}
	if req.Reason == "OTHER" && isBlank(req.OtherReason) {
		return dto.SuccessResponse{}, newBadRequest("other_reason is required when reason is OTHER")
	}
	if req.Reason == "OTHER" {
		if !strings.HasPrefix(strings.TrimSpace(req.OtherReason), "Bổ sung hồ sơ") {
			return dto.SuccessResponse{}, newBadRequest("other_reason is invalid")
		}
	} else {
		if !isBlank(req.OtherReason) {
			return dto.SuccessResponse{}, newBadRequest("other_reason is not allowed")
		}
	}

	created, appErr := s.createRequest(student.ID, "CONFIRM_LETTER", req)
	if appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(created), nil
}

func (s *Service) BankLoans(studentID string, req dto.BankLoansRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	if isBlank(req.Benefit) || isBlank(req.OrphanStatus) || isBlank(req.Template) {
		return dto.SuccessResponse{}, newBadRequest("benefit, orphan_status, template are required")
	}
	switch req.Benefit {
	case "NO_DISCOUNT", "TUITION_REDUCTION", "TUITION_EXEMPTION":
	default:
		return dto.SuccessResponse{}, newBadRequest("benefit is invalid")
	}
	switch req.OrphanStatus {
	case "NOT_ORPHAN", "ORPHAN":
	default:
		return dto.SuccessResponse{}, newBadRequest("orphan_status is invalid")
	}
	switch req.Template {
	case "LEGACY", "STEM":
	default:
		return dto.SuccessResponse{}, newBadRequest("template is invalid")
	}

	created, appErr := s.createRequest(student.ID, "BANK_LOANS", req)
	if appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(created), nil
}

func (s *Service) TrainingPointConfirm(studentID string, req dto.TrainingPointConfirmRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	if isBlank(req.Language) {
		return dto.SuccessResponse{}, newBadRequest("language is required")
	}
	switch req.Language {
	case "VI", "EN":
	default:
		return dto.SuccessResponse{}, newBadRequest("language is invalid")
	}

	created, appErr := s.createRequest(student.ID, "TRAINING_POINT_CONFIRM", req)
	if appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(created), nil
}

func (s *Service) LanguageCertificate(studentID string, req dto.LanguageCertificateUploadRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	if isBlank(req.DocumentType) || isBlank(req.BirthDate) || isBlank(req.IDNumber) || isBlank(req.ExamDate) || isBlank(req.ImageFile) {
		return dto.SuccessResponse{}, newBadRequest("document_type, birth_date, id_number, exam_date, image_file are required")
	}
	if len(strings.TrimSpace(req.IDNumber)) < 6 {
		return dto.SuccessResponse{}, newBadRequest("id_number is invalid")
	}
	if !isValidBase64(req.ImageFile) {
		return dto.SuccessResponse{}, newBadRequest("image_file is invalid")
	}
	switch req.DocumentType {
	case "CONFIRMATION", "DIPLOMA", "CERTIFICATE":
	default:
		return dto.SuccessResponse{}, newBadRequest("document_type is invalid")
	}
	if _, parseErr := time.Parse("2006-01-02", req.BirthDate); parseErr != nil {
		return dto.SuccessResponse{}, newBadRequest("birth_date is invalid")
	}
	if _, parseErr := time.Parse("2006-01-02", req.ExamDate); parseErr != nil {
		return dto.SuccessResponse{}, newBadRequest("exam_date is invalid")
	}
	if req.ListeningScore < 0 || req.ReadingScore < 0 || req.TotalScore < 0 {
		return dto.SuccessResponse{}, newBadRequest("scores must be greater than or equal to 0")
	}

	created, appErr := s.createRequest(student.ID, "LANGUAGE_CERTIFICATE", req)
	if appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(created), nil
}
