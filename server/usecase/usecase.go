package usecase

import (
	"encoding/json"
	"strings"
	"time"

	"server/data"
	"server/models"

	"github.com/google/uuid"
)

const (
	successMessage   = "Successfully!"
	tuitionPerCredit = 350000
)

type Service struct {
	now func() time.Time
}

func NewService() *Service {
	return &Service{now: time.Now}
}

func success(data interface{}) SuccessResponse {
	if data == nil {
		data = map[string]interface{}{}
	}

	return SuccessResponse{Message: successMessage, Data: data}
}

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func (s *Service) getStudent(studentID string) (*models.Student, *AppError) {
	if studentID == "" {
		if len(data.Students) == 0 {
			return nil, NewNotFound("student not found")
		}
		return &data.Students[0], nil
	}

	for i := range data.Students {
		if data.Students[i].ID == studentID {
			return &data.Students[i], nil
		}
	}

	return nil, NewNotFound("student not found")
}

func (s *Service) Login(req LoginRequest) (SuccessResponse, *AppError) {
	if isBlank(req.StudentID) || isBlank(req.Password) {
		return SuccessResponse{}, NewBadRequest("student_id and password are required")
	}

	for i := range data.Students {
		student := &data.Students[i]
		if student.ID == req.StudentID {
			if student.Password != req.Password {
				return SuccessResponse{}, NewUnauthorized("invalid credentials")
			}

			token := "mock-" + student.ID
			return success(map[string]string{"token": token}), nil
		}
	}

	return SuccessResponse{}, NewUnauthorized("invalid credentials")
}

func (s *Service) GetProfile(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	profile := StudentProfile{
		ID:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Phone: student.Phone,
		Major: student.Major,
	}

	return success(profile), nil
}

func (s *Service) GetSchedule(studentID string, query QueryYearSemester) (SuccessResponse, *AppError) {
	if query.Year < 1900 || query.Semester < 1 {
		return SuccessResponse{}, NewBadRequest("year and semester are required")
	}

	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	result := make([]models.Schedule, 0)
	for _, schedule := range data.Schedules {
		if schedule.StudentID == student.ID && schedule.Year == query.Year && schedule.Semester == query.Semester {
			result = append(result, schedule)
		}
	}

	return success(result), nil
}

func (s *Service) GetExamSchedule(studentID string, query QueryYearSemester) (SuccessResponse, *AppError) {
	if query.Year < 1900 || query.Semester < 1 {
		return SuccessResponse{}, NewBadRequest("year and semester are required")
	}

	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	result := make([]models.ExamSchedule, 0)
	for _, schedule := range data.ExamSchedules {
		if schedule.StudentID == student.ID && schedule.Year == query.Year && schedule.Semester == query.Semester {
			result = append(result, schedule)
		}
	}

	return success(result), nil
}

func (s *Service) GetScore(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	result := make([]models.Score, 0)
	for _, score := range data.Scores {
		if score.StudentID == student.ID {
			result = append(result, score)
		}
	}

	return success(result), nil
}

func (s *Service) GetTuitionFee(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	credits := 0
	for _, enrollment := range data.Enrollments {
		if enrollment.StudentID != student.ID {
			continue
		}
		course := findCourse(enrollment.CourseID)
		if course != nil {
			credits += course.Credits
		}
	}

	fee := TuitionFeeData{
		TotalCredits: credits,
		Amount:       credits * tuitionPerCredit,
		Currency:     "VND",
		Status:       "UNPAID",
	}

	return success(fee), nil
}

func (s *Service) GetInsurance(studentID string) (SuccessResponse, *AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	insurance := InsuranceData{
		Provider:  "BHXH",
		Status:    "ACTIVE",
		ExpiresAt: "2026-12-31",
	}

	return success(insurance), nil
}

func (s *Service) GetCourses(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	seen := make(map[string]bool)
	result := make([]models.Course, 0)
	for _, enrollment := range data.Enrollments {
		if enrollment.StudentID != student.ID {
			continue
		}
		if seen[enrollment.CourseID] {
			continue
		}
		course := findCourse(enrollment.CourseID)
		if course != nil {
			seen[course.ID] = true
			result = append(result, *course)
		}
	}

	return success(result), nil
}

func (s *Service) GetTrainingPoints(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	total := 0.0
	count := 0
	for _, score := range data.Scores {
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

	result := TrainingPointsData{
		Year:     2025,
		Semester: 1,
		Score:    point,
		Rank:     rank,
	}

	return success(result), nil
}

func (s *Service) GetSurveyForm(studentID string) (SuccessResponse, *AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	return success(map[string]interface{}{"available": false}), nil
}

func (s *Service) TranscriptRegis(studentID string, req TranscriptRegisRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	if req.Copies < 1 {
		return SuccessResponse{}, NewBadRequest("copies must be at least 1")
	}

	if isBlank(req.DeliveryMethod) {
		return SuccessResponse{}, NewBadRequest("delivery_method is required")
	}

	switch req.DeliveryMethod {
	case "PICKUP", "SHIP":
	default:
		return SuccessResponse{}, NewBadRequest("delivery_method is invalid")
	}

	if req.DeliveryMethod == "SHIP" {
		if isBlank(req.ShippingAddress) || isBlank(req.Phone) {
			return SuccessResponse{}, NewBadRequest("shipping_address and phone are required for delivery")
		}
	}

	if req.Language == "" {
		req.Language = "VI"
	}

	if req.Language != "VI" && req.Language != "EN" {
		return SuccessResponse{}, NewBadRequest("language is invalid")
	}

	created := s.createRequest(student.ID, "TRANSCRIPT_REGIS", req)
	return success(created), nil
}

func (s *Service) TuitionExtend(studentID string, req TuitionExtendRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	if req.Year < 1900 || req.Semester < 1 {
		return SuccessResponse{}, NewBadRequest("year and semester are required")
	}
	if isBlank(req.RequestedDueDate) || isBlank(req.Reason) {
		return SuccessResponse{}, NewBadRequest("requested_due_date and reason are required")
	}
	if _, parseErr := time.Parse("2006-01-02", req.RequestedDueDate); parseErr != nil {
		return SuccessResponse{}, NewBadRequest("requested_due_date is invalid")
	}

	created := s.createRequest(student.ID, "TUITION_EXTEND", req)
	return success(created), nil
}

func (s *Service) MonthlyParking(studentID string, req MonthlyParkingRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	if isBlank(req.VehicleType) || isBlank(req.PlateNumber) || req.Months < 1 || isBlank(req.StartMonth) {
		return SuccessResponse{}, NewBadRequest("vehicle_type, plate_number, months, start_month are required")
	}
	if req.Months > 12 {
		return SuccessResponse{}, NewBadRequest("months must be between 1 and 12")
	}
	switch req.VehicleType {
	case "MOTORBIKE", "CAR", "BICYCLE":
	default:
		return SuccessResponse{}, NewBadRequest("vehicle_type is invalid")
	}
	if _, parseErr := time.Parse("2006-01", req.StartMonth); parseErr != nil {
		return SuccessResponse{}, NewBadRequest("start_month is invalid")
	}

	created := s.createRequest(student.ID, "MONTHLY_PARKING", req)
	return success(created), nil
}

func (s *Service) Graduate(studentID string, req GraduateRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	if req.Year < 1900 || req.Semester < 1 {
		return SuccessResponse{}, NewBadRequest("year and semester are required")
	}
	if isBlank(req.Email) || isBlank(req.Phone) {
		return SuccessResponse{}, NewBadRequest("email and phone are required")
	}

	created := s.createRequest(student.ID, "GRADUATE", req)
	return success(created), nil
}

func (s *Service) GraduationThesis(studentID string, req GraduationThesisRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	if isBlank(req.ThesisTitle) || isBlank(req.AdvisorName) {
		return SuccessResponse{}, NewBadRequest("thesis_title and advisor_name are required")
	}
	if len(req.TeamMembers) == 0 {
		return SuccessResponse{}, NewBadRequest("team_members are required")
	}
	for _, member := range req.TeamMembers {
		if isBlank(member.StudentID) {
			return SuccessResponse{}, NewBadRequest("team_members.student_id is required")
		}
	}

	created := s.createRequest(student.ID, "GRADUATION_THESIS", req)
	return success(created), nil
}

func (s *Service) GetDeadlines(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	result := make([]models.Deadline, 0)
	for _, deadline := range data.Deadlines {
		if deadline.StudentID == student.ID {
			result = append(result, deadline)
		}
	}

	return success(result), nil
}

func (s *Service) GetMaterials(studentID, courseID string) (SuccessResponse, *AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	if findCourse(courseID) == nil {
		return SuccessResponse{}, NewNotFound("course not found")
	}

	result := make([]models.Material, 0)
	for _, material := range data.Materials {
		if material.CourseID == courseID {
			result = append(result, material)
		}
	}

	return success(result), nil
}

func (s *Service) GetAssignments(studentID, courseID string) (SuccessResponse, *AppError) {
	_, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	if findCourse(courseID) == nil {
		return SuccessResponse{}, NewNotFound("course not found")
	}

	result := make([]models.Assignment, 0)
	for _, assignment := range data.Assignments {
		if assignment.CourseID == courseID {
			result = append(result, assignment)
		}
	}

	return success(result), nil
}

func (s *Service) SubmitAssignment(studentID, assignmentID string, req AssignmentSubmissionRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}

	assignment := findAssignment(assignmentID)
	if assignment == nil {
		return SuccessResponse{}, NewNotFound("assignment not found")
	}

	submissionType := strings.ToUpper(strings.TrimSpace(req.SubmissionType))
	if submissionType == "" {
		return SuccessResponse{}, NewBadRequest("submission_type is required")
	}

	content := ""
	switch submissionType {
	case "TEXT":
		if isBlank(req.Text) {
			return SuccessResponse{}, NewBadRequest("text is required for TEXT submissions")
		}
		content = req.Text
	case "LINK":
		if isBlank(req.URL) {
			return SuccessResponse{}, NewBadRequest("url is required for LINK submissions")
		}
		content = req.URL
	case "FILE":
		if len(req.FileURLs) == 0 {
			return SuccessResponse{}, NewBadRequest("file_urls are required for FILE submissions")
		}
		content = strings.Join(req.FileURLs, ",")
	default:
		return SuccessResponse{}, NewBadRequest("submission_type is invalid")
	}

	submissionID := uuid.NewString()
	data.Submissions = append(data.Submissions, models.Submission{
		ID:             uint(len(data.Submissions) + 1),
		AssignmentID:   assignment.ID,
		StudentID:      student.ID,
		SubmissionType: submissionType,
		Content:        content,
		CreatedAt:      s.now().Format(time.RFC3339),
	})

	response := SubmissionResponse{SubmissionID: submissionID, Status: "SUBMITTED"}
	return success(response), nil
}

func (s *Service) ConfirmLetter(studentID string, req ConfirmLetterRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	if isBlank(req.Language) || isBlank(req.Reason) || isBlank(req.RequestType) {
		return SuccessResponse{}, NewBadRequest("language, reason, request_type are required")
	}
	switch req.Language {
	case "VI", "EN":
	default:
		return SuccessResponse{}, NewBadRequest("language is invalid")
	}
	switch req.RequestType {
	case "NEW", "REISSUE":
	default:
		return SuccessResponse{}, NewBadRequest("request_type is invalid")
	}
	switch req.Reason {
	case "MILITARY_DEFERMENT", "DORM_EXTEND", "TAX_DEDUCTION_DOCS", "DEFENSE_EDU_REGISTRATION", "OTHER":
	default:
		return SuccessResponse{}, NewBadRequest("reason is invalid")
	}
	if req.Reason == "OTHER" && isBlank(req.OtherReason) {
		return SuccessResponse{}, NewBadRequest("other_reason is required when reason is OTHER")
	}

	created := s.createRequest(student.ID, "CONFIRM_LETTER", req)
	return success(created), nil
}

func (s *Service) BankLoans(studentID string, req BankLoansRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	if isBlank(req.Benefit) || isBlank(req.OrphanStatus) || isBlank(req.Template) {
		return SuccessResponse{}, NewBadRequest("benefit, orphan_status, template are required")
	}
	switch req.Benefit {
	case "NO_DISCOUNT", "TUITION_REDUCTION", "TUITION_EXEMPTION":
	default:
		return SuccessResponse{}, NewBadRequest("benefit is invalid")
	}
	switch req.OrphanStatus {
	case "NOT_ORPHAN", "ORPHAN":
	default:
		return SuccessResponse{}, NewBadRequest("orphan_status is invalid")
	}
	switch req.Template {
	case "LEGACY", "STEM":
	default:
		return SuccessResponse{}, NewBadRequest("template is invalid")
	}

	created := s.createRequest(student.ID, "BANK_LOANS", req)
	return success(created), nil
}

func (s *Service) TrainingPointConfirm(studentID string, req TrainingPointConfirmRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	if isBlank(req.Language) {
		return SuccessResponse{}, NewBadRequest("language is required")
	}
	switch req.Language {
	case "VI", "EN":
	default:
		return SuccessResponse{}, NewBadRequest("language is invalid")
	}

	created := s.createRequest(student.ID, "TRAINING_POINT_CONFIRM", req)
	return success(created), nil
}

func (s *Service) LanguageCertificate(studentID string, req LanguageCertificateUploadRequest) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	if isBlank(req.DocumentType) || isBlank(req.BirthDate) || isBlank(req.IDNumber) || isBlank(req.ExamDate) || isBlank(req.ImageFile) {
		return SuccessResponse{}, NewBadRequest("document_type, birth_date, id_number, exam_date, image_file are required")
	}
	switch req.DocumentType {
	case "CONFIRMATION", "DIPLOMA", "CERTIFICATE":
	default:
		return SuccessResponse{}, NewBadRequest("document_type is invalid")
	}
	if _, parseErr := time.Parse("2006-01-02", req.BirthDate); parseErr != nil {
		return SuccessResponse{}, NewBadRequest("birth_date is invalid")
	}
	if _, parseErr := time.Parse("2006-01-02", req.ExamDate); parseErr != nil {
		return SuccessResponse{}, NewBadRequest("exam_date is invalid")
	}
	if req.ListeningScore < 0 || req.ReadingScore < 0 || req.TotalScore < 0 {
		return SuccessResponse{}, NewBadRequest("scores must be greater than or equal to 0")
	}

	created := s.createRequest(student.ID, "LANGUAGE_CERTIFICATE", req)
	return success(created), nil
}

func (s *Service) GetRoomsAvailability(query QueryRoomsAvailability) (SuccessResponse, *AppError) {
	if isBlank(query.Date) || isBlank(query.Start) || isBlank(query.End) {
		return SuccessResponse{}, NewBadRequest("date, start, end are required")
	}
	if _, parseErr := time.Parse("2006-01-02", query.Date); parseErr != nil {
		return SuccessResponse{}, NewBadRequest("date is invalid")
	}

	startMin, err := parseTimeToMinutes(query.Start)
	if err != nil {
		return SuccessResponse{}, NewBadRequest("start is invalid")
	}
	endMin, err := parseTimeToMinutes(query.End)
	if err != nil {
		return SuccessResponse{}, NewBadRequest("end is invalid")
	}
	if startMin >= endMin {
		return SuccessResponse{}, NewBadRequest("start must be before end")
	}

	available := make([]RoomItem, 0)
	for _, room := range data.Rooms {
		if isRoomAvailable(room.ID, query.Date, startMin, endMin) {
			available = append(available, RoomItem{ID: room.ID, Name: room.Name, Capacity: room.Capacity})
		}
	}

	result := RoomAvailabilityData{Date: query.Date, Start: query.Start, End: query.End, Rooms: available}
	return success(result), nil
}

func (s *Service) Contact(req ContactRequest) (SuccessResponse, *AppError) {
	if isBlank(req.Name) || isBlank(req.Email) || isBlank(req.Subject) || isBlank(req.Message) {
		return SuccessResponse{}, NewBadRequest("name, email, subject, message are required")
	}

	contactID := uuid.NewString()
	data.Contacts = append(data.Contacts, models.Contact{
		ID:        uint(len(data.Contacts) + 1),
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Subject:   req.Subject,
		Message:   req.Message,
		CreatedAt: s.now().Format(time.RFC3339),
	})

	return success(map[string]string{"contact_id": contactID}), nil
}

func (s *Service) createRequest(studentID, requestType string, payload interface{}) RequestStatusData {
	requestID := uuid.NewString()
	payloadBytes, _ := json.Marshal(payload)
	data.Requests = append(data.Requests, models.Request{
		ID:        requestID,
		StudentID: studentID,
		Type:      requestType,
		Status:    "PENDING",
		Payload:   string(payloadBytes),
		CreatedAt: s.now().Format(time.RFC3339),
	})

	return RequestStatusData{
		RequestID: requestID,
		Status:    "PENDING",
		CreatedAt: s.now().Unix(),
		PDFURL:    nil,
	}
}

func findCourse(courseID string) *models.Course {
	for i := range data.Courses {
		if data.Courses[i].ID == courseID {
			return &data.Courses[i]
		}
	}

	return nil
}

func findAssignment(assignmentID string) *models.Assignment {
	for i := range data.Assignments {
		if data.Assignments[i].ID == assignmentID {
			return &data.Assignments[i]
		}
	}

	return nil
}

func parseTimeToMinutes(value string) (int, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}

	return parsed.Hour()*60 + parsed.Minute(), nil
}

func isRoomAvailable(roomID uint, date string, startMin, endMin int) bool {
	for _, booking := range data.RoomBookings {
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
}
