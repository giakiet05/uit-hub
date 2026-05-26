package service

<<<<<<< HEAD
import (
	"encoding/json"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"
	"github.com/google/uuid"
)

func (s *Service) getStudent(studentID string) (*model.Student, *apperror.AppError) {
	if studentID == "" {
		return nil, newUnauthorized("authorization required")
	}
	student, ok, err := s.students.FindByID(studentID)
	if err != nil {
		return nil, newInternal("")
	}
	if !ok {
		return nil, newNotFound("student not found")
	}
	return student, nil
}

func (s *Service) ListStudents() (dto.SuccessResponse, *apperror.AppError) {
	students, err := s.students.List()
	if err != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]dto.StudentProfile, 0, len(students))
	for _, student := range students {
		result = append(result, dto.StudentProfile{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
			Phone: student.Phone,
			Major: student.Major,
		})
	}

	return success(result), nil
}

func (s *Service) GetStudentByID(id string) (dto.SuccessResponse, *apperror.AppError) {
	if isBlank(id) {
		return dto.SuccessResponse{}, newBadRequest("id is required")
	}

	student, ok, err := s.students.FindByID(id)
	if err != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if !ok {
		return dto.SuccessResponse{}, newNotFound("student not found")
	}

	profile := dto.StudentProfile{
		ID:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Phone: student.Phone,
		Major: student.Major,
	}

	return success(profile), nil
}

func (s *Service) GetProfile(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	profile := dto.StudentProfile{
		ID:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Phone: student.Phone,
		Major: student.Major,
	}

	return success(profile), nil
}

func (s *Service) GetSchedule(studentID string, query dto.QueryYearSemester) (dto.SuccessResponse, *apperror.AppError) {
	if query.Year < 1900 || query.Semester < 1 {
		return dto.SuccessResponse{}, newBadRequest("year and semester are required")
	}

	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	schedules, repoErr := s.schedules.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]model.Schedule, 0)
	for _, schedule := range schedules {
		if schedule.StudentID == student.ID && schedule.Year == query.Year && schedule.Semester == query.Semester {
			result = append(result, schedule)
		}
	}

	return success(result), nil
}

func (s *Service) GetExamSchedule(studentID string, query dto.QueryYearSemester) (dto.SuccessResponse, *apperror.AppError) {
	if query.Year < 1900 || query.Semester < 1 {
		return dto.SuccessResponse{}, newBadRequest("year and semester are required")
	}

	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	examSchedules, repoErr := s.examSchedules.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]model.ExamSchedule, 0)
	for _, schedule := range examSchedules {
		if schedule.StudentID == student.ID && schedule.Year == query.Year && schedule.Semester == query.Semester {
			result = append(result, schedule)
		}
	}

	return success(result), nil
}

func (s *Service) GetScore(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	scores, repoErr := s.scores.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	result := make([]model.Score, 0)
	for _, score := range scores {
		if score.StudentID == student.ID {
			result = append(result, score)
		}
	}

	return success(result), nil
}

func (s *Service) GetTuitionFee(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	enrollments, repoErr := s.enrollments.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	credits := 0
	for _, enrollment := range enrollments {
		if enrollment.StudentID != student.ID {
			continue
		}
		course, ok, repoErr := s.courses.FindByID(enrollment.CourseID)
		if repoErr != nil {
			return dto.SuccessResponse{}, newInternal("")
		}
		if ok {
			credits += course.Credits
		}
	}

	fee := dto.TuitionFeeData{
		TotalCredits: credits,
		Amount:       credits * tuitionPerCredit,
		Currency:     "VND",
		Status:       "UNPAID",
	}

	return success(fee), nil
}

func (s *Service) GetInsurance(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	insurance := dto.InsuranceData{
		Provider:  "BHXH",
		Status:    "ACTIVE",
		ExpiresAt: "2026-12-31",
	}

	return success(insurance), nil
}

func (s *Service) GetCourses(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	enrollments, repoErr := s.enrollments.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	seen := make(map[string]bool)
	result := make([]model.Course, 0)
	for _, enrollment := range enrollments {
		if enrollment.StudentID != student.ID {
			continue
		}
		if seen[enrollment.CourseID] {
			continue
		}
		course, ok, repoErr := s.courses.FindByID(enrollment.CourseID)
		if repoErr != nil {
			return dto.SuccessResponse{}, newInternal("")
		}
		if ok {
			seen[course.ID] = true
			result = append(result, *course)
		}
	}

	return success(result), nil
}

func (s *Service) GetTrainingPoints(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}
	scores, repoErr := s.scores.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	total := 0.0
	count := 0
	for _, score := range scores {
		if score.StudentID == student.ID {
			total += score.Total
			count++
		}
	}

	avg := 0.0
	if count > 0 {
		avg = total / float64(count)
	}

	point := int(avg * 10)
	rank := "D"
	if point >= 90 {
		rank = "A"
	} else if point >= 80 {
		rank = "B"
	} else if point >= 65 {
		rank = "C"
	}

	result := dto.TrainingPointsData{
		Year:     2025,
		Semester: 1,
		Score:    point,
		Rank:     rank,
	}

	return success(result), nil
}

func (s *Service) GetSurveyForm(studentID string) (dto.SuccessResponse, *apperror.AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	return success(map[string]interface{}{"available": false}), nil
}

func (s *Service) TranscriptRegis(studentID string, req dto.TranscriptRegisRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	if req.Copies < 1 {
		return dto.SuccessResponse{}, newBadRequest("copies must be at least 1")
	}

	if isBlank(req.DeliveryMethod) {
		return dto.SuccessResponse{}, newBadRequest("delivery_method is required")
	}

	switch req.DeliveryMethod {
	case "PICKUP", "SHIP":
	default:
		return dto.SuccessResponse{}, newBadRequest("delivery_method is invalid")
	}

	if req.DeliveryMethod == "SHIP" {
		if isBlank(req.ShippingAddress) || isBlank(req.Phone) {
			return dto.SuccessResponse{}, newBadRequest("shipping_address and phone are required for delivery")
		}
		if len(strings.TrimSpace(req.Phone)) < 6 {
			return dto.SuccessResponse{}, newBadRequest("phone is invalid")
		}
	}

	if req.Language == "" {
		req.Language = "VI"
	}

	if req.Language != "VI" && req.Language != "EN" {
		return dto.SuccessResponse{}, newBadRequest("language is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "TRANSCRIPT_REGIS", req); appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(nil), nil
}

func (s *Service) TuitionExtend(studentID string, req dto.TuitionExtendRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	if req.Year < 1900 || req.Semester < 1 {
		return dto.SuccessResponse{}, newBadRequest("year and semester are required")
	}
	if isBlank(req.RequestedDueDate) || isBlank(req.Reason) {
		return dto.SuccessResponse{}, newBadRequest("requested_due_date and reason are required")
	}
	if _, parseErr := time.Parse("2006-01-02", req.RequestedDueDate); parseErr != nil {
		return dto.SuccessResponse{}, newBadRequest("requested_due_date is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "TUITION_EXTEND", req); appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(nil), nil
}

func (s *Service) MonthlyParking(studentID string, req dto.MonthlyParkingRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	if isBlank(req.VehicleType) || isBlank(req.PlateNumber) || req.Months < 1 || isBlank(req.StartMonth) {
		return dto.SuccessResponse{}, newBadRequest("vehicle_type, plate_number, months, start_month are required")
	}
	if req.Months > 12 {
		return dto.SuccessResponse{}, newBadRequest("months must be between 1 and 12")
	}
	if len(strings.TrimSpace(req.PlateNumber)) < 3 {
		return dto.SuccessResponse{}, newBadRequest("plate_number is invalid")
	}
	switch req.VehicleType {
	case "MOTORBIKE", "CAR", "BICYCLE":
	default:
		return dto.SuccessResponse{}, newBadRequest("vehicle_type is invalid")
	}
	if _, parseErr := time.Parse("2006-01", req.StartMonth); parseErr != nil {
		return dto.SuccessResponse{}, newBadRequest("start_month is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "MONTHLY_PARKING", req); appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(nil), nil
}

func (s *Service) Graduate(studentID string, req dto.GraduateRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	if req.Year < 1900 || req.Semester < 1 {
		return dto.SuccessResponse{}, newBadRequest("year and semester are required")
	}
	if isBlank(req.Email) || isBlank(req.Phone) {
		return dto.SuccessResponse{}, newBadRequest("email and phone are required")
	}
	if !isValidEmail(req.Email) {
		return dto.SuccessResponse{}, newBadRequest("email is invalid")
	}
	if len(strings.TrimSpace(req.Phone)) < 6 {
		return dto.SuccessResponse{}, newBadRequest("phone is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "GRADUATE", req); appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(nil), nil
}

func (s *Service) GraduationThesis(studentID string, req dto.GraduationThesisRequest) (dto.SuccessResponse, *apperror.AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return dto.SuccessResponse{}, err
	}

	if isBlank(req.ThesisTitle) || isBlank(req.AdvisorName) {
		return dto.SuccessResponse{}, newBadRequest("thesis_title and advisor_name are required")
	}
	if len(req.TeamMembers) == 0 {
		return dto.SuccessResponse{}, newBadRequest("team_members are required")
	}
	for _, member := range req.TeamMembers {
		if isBlank(member.StudentID) {
			return dto.SuccessResponse{}, newBadRequest("team_members.student_id is required")
		}
	}
	if !isBlank(req.AdvisorEmail) && !isValidEmail(req.AdvisorEmail) {
		return dto.SuccessResponse{}, newBadRequest("advisor_email is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "GRADUATION_THESIS", req); appErr != nil {
		return dto.SuccessResponse{}, appErr
	}
	return success(nil), nil
}

func (s *Service) Contact(req dto.ContactRequest) (dto.SuccessResponse, *apperror.AppError) {
	if isBlank(req.Name) || isBlank(req.Email) || isBlank(req.Subject) || isBlank(req.Message) {
		return dto.SuccessResponse{}, newBadRequest("name, email, subject, message are required")
	}
	if !isValidEmail(req.Email) {
		return dto.SuccessResponse{}, newBadRequest("email is invalid")
	}

	_ = uuid.NewString()
	count, repoErr := s.contacts.Count()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	if err := s.contacts.Create(model.Contact{
		ID:        uint(count + 1),
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Subject:   req.Subject,
		Message:   req.Message,
		CreatedAt: s.now().Format(time.RFC3339),
	}); err != nil {
		return dto.SuccessResponse{}, newInternal("")
	}

	return success(nil), nil
}

func (s *Service) createRequest(studentID, requestType string, payload interface{}) (dto.RequestStatusData, *apperror.AppError) {
	requestID := uuid.NewString()
	payloadBytes, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return dto.RequestStatusData{}, newInternal("")
	}
	if err := s.requests.Create(model.Request{
		ID:        requestID,
		StudentID: studentID,
		Type:      requestType,
		Status:    "PENDING",
		Payload:   string(payloadBytes),
		CreatedAt: s.now().Format(time.RFC3339),
	}); err != nil {
		return dto.RequestStatusData{}, newInternal("")
	}

	return dto.RequestStatusData{
		RequestID: requestID,
		Status:    "PENDING",
		CreatedAt: s.now().Unix(),
		PDFURL:    nil,
	}, nil
}

func parseTimeToMinutes(value string) (int, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}

	return parsed.Hour()*60 + parsed.Minute(), nil
}

func isRoomAvailable(roomID uint, date string, startMin, endMin int, bookings []model.RoomBooking) bool {
	for _, booking := range bookings {
		if booking.RoomID != roomID || booking.Date != date {
			continue
		}

		bookingStart, err := parseTimeToMinutes(booking.StartTime)
		if err != nil {
			continue
		}
		bookingEnd, err := parseTimeToMinutes(booking.EndTime)
		if err != nil {
			continue
		}

		if startMin < bookingEnd && endMin > bookingStart {
			return false
		}
	}

	return true
=======
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
>>>>>>> origin/dev
}
