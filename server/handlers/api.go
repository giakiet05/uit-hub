package handlers

import (
	"os"
	"strings"

	"server/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	uc    *usecase.Service
	reset func()
}

func NewHandler(uc *usecase.Service, reset func()) *Handler {
	return &Handler{uc: uc, reset: reset}
}

func (h *Handler) respond(c *fiber.Ctx, resp usecase.SuccessResponse, err *usecase.AppError) error {
	if err != nil {
		return c.Status(err.Status).JSON(usecase.ErrorResponse{
			Message:   err.Message,
			ErrorCode: err.Code,
		})
	}

	return c.JSON(resp)
}

func (h *Handler) parseBody(c *fiber.Ctx, out interface{}) *usecase.AppError {
	if err := c.BodyParser(out); err != nil {
		return usecase.NewBadRequest("invalid json body")
	}
	return nil
}

func (h *Handler) requireAuth(c *fiber.Ctx) (string, *usecase.AppError) {
	studentID := getStudentID(c)
	if studentID == "" {
		return "", usecase.NewUnauthorized("authorization required")
	}
	return studentID, nil
}

func getStudentID(c *fiber.Ctx) string {
	authorization := strings.TrimSpace(c.Get("Authorization"))
	if strings.HasPrefix(authorization, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		if strings.HasPrefix(token, "mock-") {
			return strings.TrimPrefix(token, "mock-")
		}
		return token
	}

	return ""
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req usecase.LoginRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.Login(req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) ResetStore(c *fiber.Ctx) error {
	if h.reset == nil {
		return h.respond(c, usecase.SuccessResponse{}, usecase.NewInternal("reset is not configured"))
	}

	adminToken := strings.TrimSpace(os.Getenv("ADMIN_TOKEN"))
	if adminToken != "" {
		provided := strings.TrimSpace(c.Get("X-Admin-Token"))
		if provided == "" {
			return h.respond(c, usecase.SuccessResponse{}, usecase.NewUnauthorized("admin token required"))
		}
		if provided != adminToken {
			return h.respond(c, usecase.SuccessResponse{}, usecase.NewUnauthorized("invalid admin token"))
		}
	}

	h.reset()
	return h.respond(c, usecase.SuccessResponse{Message: "Successfully!", Data: map[string]interface{}{}}, nil)
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetProfile(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetSchedule(c *fiber.Ctx) error {
	query := usecase.QueryYearSemester{Year: c.QueryInt("year"), Semester: c.QueryInt("semester")}
	resp, appErr := h.uc.GetSchedule(getStudentID(c), query)
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetExamSchedule(c *fiber.Ctx) error {
	query := usecase.QueryYearSemester{Year: c.QueryInt("year"), Semester: c.QueryInt("semester")}
	resp, appErr := h.uc.GetExamSchedule(getStudentID(c), query)
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetScore(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetScore(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetTuitionFee(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetTuitionFee(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetInsurance(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetInsurance(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetCourses(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetCourses(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetTrainingPoints(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetTrainingPoints(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetSurveyForm(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetSurveyForm(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) TranscriptRegis(c *fiber.Ctx) error {
	var req usecase.TranscriptRegisRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.TranscriptRegis(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) TuitionExtend(c *fiber.Ctx) error {
	var req usecase.TuitionExtendRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.TuitionExtend(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) MonthlyParking(c *fiber.Ctx) error {
	var req usecase.MonthlyParkingRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.MonthlyParking(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) Graduate(c *fiber.Ctx) error {
	var req usecase.GraduateRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.Graduate(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) GraduationThesis(c *fiber.Ctx) error {
	var req usecase.GraduationThesisRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.GraduationThesis(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetDeadlines(c *fiber.Ctx) error {
	resp, appErr := h.uc.GetDeadlines(getStudentID(c))
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetMaterials(c *fiber.Ctx) error {
	courseID := c.Params("courseId")
	resp, appErr := h.uc.GetMaterials(getStudentID(c), courseID)
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetAssignments(c *fiber.Ctx) error {
	courseID := c.Params("courseId")
	resp, appErr := h.uc.GetAssignments(getStudentID(c), courseID)
	return h.respond(c, resp, appErr)
}

func (h *Handler) SubmitAssignment(c *fiber.Ctx) error {
	assignmentID := c.Params("assignmentId")
	var req usecase.AssignmentSubmissionRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.SubmitAssignment(getStudentID(c), assignmentID, req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) ConfirmLetter(c *fiber.Ctx) error {
	var req usecase.ConfirmLetterRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.ConfirmLetter(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) BankLoans(c *fiber.Ctx) error {
	var req usecase.BankLoansRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.BankLoans(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) TrainingPointConfirm(c *fiber.Ctx) error {
	var req usecase.TrainingPointConfirmRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.TrainingPointConfirm(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) LanguageCertificate(c *fiber.Ctx) error {
	var req usecase.LanguageCertificateUploadRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.LanguageCertificate(getStudentID(c), req)
	return h.respond(c, resp, appErr)
}

func (h *Handler) GetRoomsAvailability(c *fiber.Ctx) error {
	if _, err := h.requireAuth(c); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}
	query := usecase.QueryRoomsAvailability{Date: c.Query("date"), Start: c.Query("start"), End: c.Query("end")}
	resp, appErr := h.uc.GetRoomsAvailability(query)
	return h.respond(c, resp, appErr)
}

func (h *Handler) Contact(c *fiber.Ctx) error {
	if _, err := h.requireAuth(c); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}
	var req usecase.ContactRequest
	if err := h.parseBody(c, &req); err != nil {
		return h.respond(c, usecase.SuccessResponse{}, err)
	}

	resp, appErr := h.uc.Contact(req)
	return h.respond(c, resp, appErr)
}
