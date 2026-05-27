package service

import (
	"net/http"

	"encoding/base64"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
)

func newAppError(status int, code, msg string) *apperror.AppError {
	return &apperror.AppError{Status: status, Code: code, Message: msg}
}

func newBadRequest(msg string) *apperror.AppError {
	if msg == "" {
		msg = "Bad request!"
	}
	return newAppError(http.StatusBadRequest, "BAD_REQUEST", msg)
}
func newUnauthorized(msg string) *apperror.AppError {
	if msg == "" {
		msg = "Unauthorized!"
	}
	return newAppError(http.StatusUnauthorized, "UNAUTHORIZED", msg)
}
func newNotFound(msg string) *apperror.AppError {
	if msg == "" {
		msg = "Not found!"
	}
	return newAppError(http.StatusNotFound, "NOT_FOUND", msg)
}
func newInternal(msg string) *apperror.AppError {
	if msg == "" {
		msg = "Internal error happened!"
	}
	return newAppError(http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

const (
	successMessage   = "Successfully!"
	tuitionPerCredit = 350000
)

type Service struct {
	now func() time.Time

	students      StudentRepo
	courses       CourseRepo
	enrollments   EnrollmentRepo
	schedules     ScheduleRepo
	examSchedules ExamScheduleRepo
	scores        ScoreRepo
	assignments   AssignmentRepo
	materials     MaterialRepo
	deadlines     DeadlineRepo
	rooms         RoomRepo
	roomBookings  RoomBookingRepo
	submissions   SubmissionRepo
	requests      RequestRepo
	contacts      ContactRepo
}

func NewService(deps Deps) *Service {
	return &Service{
		now:           time.Now,
		students:      deps.Students,
		courses:       deps.Courses,
		enrollments:   deps.Enrollments,
		schedules:     deps.Schedules,
		examSchedules: deps.ExamSchedules,
		scores:        deps.Scores,
		assignments:   deps.Assignments,
		materials:     deps.Materials,
		deadlines:     deps.Deadlines,
		rooms:         deps.Rooms,
		roomBookings:  deps.RoomBookings,
		submissions:   deps.Submissions,
		requests:      deps.Requests,
		contacts:      deps.Contacts,
	}
}

func success(data interface{}) dto.SuccessResponse {
	if data == nil {
		data = map[string]interface{}{}
	}

	return dto.SuccessResponse{Message: successMessage, Data: data}
}

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func isValidEmail(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	_, err := mail.ParseAddress(value)
	return err == nil
}

func isValidURL(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func isValidBase64(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if strings.HasPrefix(value, "data:") {
		if idx := strings.Index(value, ","); idx >= 0 {
			value = value[idx+1:]
		}
	}
	_, err := base64.StdEncoding.DecodeString(value)
	if err == nil {
		return true
	}
	_, err = base64.RawStdEncoding.DecodeString(value)
	return err == nil
}
