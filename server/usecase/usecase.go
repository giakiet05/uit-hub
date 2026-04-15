package usecase

import (
	"encoding/base64"
	"encoding/json"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"server/models"

	"github.com/google/uuid"
)

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

func success(data interface{}) SuccessResponse {
	if data == nil {
		data = map[string]interface{}{}
	}

	return SuccessResponse{Message: successMessage, Data: data}
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

func (s *Service) getStudent(studentID string) (*models.Student, *AppError) {
	if studentID == "" {
		return nil, NewUnauthorized("authorization required")
	}
	student, ok, err := s.students.FindByID(studentID)
	if err != nil {
		return nil, NewInternal("")
	}
	if !ok {
		return nil, NewNotFound("student not found")
	}
	return student, nil
}

func (s *Service) Login(req LoginRequest) (SuccessResponse, *AppError) {
	if isBlank(req.StudentID) || isBlank(req.Password) {
		return SuccessResponse{}, NewBadRequest("student_id and password are required")
	}
	student, ok, err := s.students.FindByID(req.StudentID)
	if err != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if !ok || student.Password != req.Password {
		return SuccessResponse{}, NewUnauthorized("invalid credentials")
	}

	token := "mock-" + student.ID
	return success(map[string]string{"token": token}), nil
}

func (s *Service) ListStudents() (SuccessResponse, *AppError) {
	students, err := s.students.List()
	if err != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]StudentProfile, 0, len(students))
	for _, student := range students {
		result = append(result, StudentProfile{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
			Phone: student.Phone,
			Major: student.Major,
		})
	}

	return success(result), nil
}

func (s *Service) GetStudentByID(id string) (SuccessResponse, *AppError) {
	if isBlank(id) {
		return SuccessResponse{}, NewBadRequest("id is required")
	}

	student, ok, err := s.students.FindByID(id)
	if err != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if !ok {
		return SuccessResponse{}, NewNotFound("student not found")
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
	schedules, repoErr := s.schedules.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]models.Schedule, 0)
	for _, schedule := range schedules {
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
	examSchedules, repoErr := s.examSchedules.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]models.ExamSchedule, 0)
	for _, schedule := range examSchedules {
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
	scores, repoErr := s.scores.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]models.Score, 0)
	for _, score := range scores {
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
	enrollments, repoErr := s.enrollments.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	credits := 0
	for _, enrollment := range enrollments {
		if enrollment.StudentID != student.ID {
			continue
		}
		course, ok, repoErr := s.courses.FindByID(enrollment.CourseID)
		if repoErr != nil {
			return SuccessResponse{}, NewInternal("")
		}
		if ok {
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
	enrollments, repoErr := s.enrollments.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	seen := make(map[string]bool)
	result := make([]models.Course, 0)
	for _, enrollment := range enrollments {
		if enrollment.StudentID != student.ID {
			continue
		}
		if seen[enrollment.CourseID] {
			continue
		}
		course, ok, repoErr := s.courses.FindByID(enrollment.CourseID)
		if repoErr != nil {
			return SuccessResponse{}, NewInternal("")
		}
		if ok {
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
	scores, repoErr := s.scores.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
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
		if len(strings.TrimSpace(req.Phone)) < 6 {
			return SuccessResponse{}, NewBadRequest("phone is invalid")
		}
	}

	if req.Language == "" {
		req.Language = "VI"
	}

	if req.Language != "VI" && req.Language != "EN" {
		return SuccessResponse{}, NewBadRequest("language is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "TRANSCRIPT_REGIS", req); appErr != nil {
		return SuccessResponse{}, appErr
	}
	return success(nil), nil
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

	if _, appErr := s.createRequest(student.ID, "TUITION_EXTEND", req); appErr != nil {
		return SuccessResponse{}, appErr
	}
	return success(nil), nil
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
	if len(strings.TrimSpace(req.PlateNumber)) < 3 {
		return SuccessResponse{}, NewBadRequest("plate_number is invalid")
	}
	switch req.VehicleType {
	case "MOTORBIKE", "CAR", "BICYCLE":
	default:
		return SuccessResponse{}, NewBadRequest("vehicle_type is invalid")
	}
	if _, parseErr := time.Parse("2006-01", req.StartMonth); parseErr != nil {
		return SuccessResponse{}, NewBadRequest("start_month is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "MONTHLY_PARKING", req); appErr != nil {
		return SuccessResponse{}, appErr
	}
	return success(nil), nil
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
	if !isValidEmail(req.Email) {
		return SuccessResponse{}, NewBadRequest("email is invalid")
	}
	if len(strings.TrimSpace(req.Phone)) < 6 {
		return SuccessResponse{}, NewBadRequest("phone is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "GRADUATE", req); appErr != nil {
		return SuccessResponse{}, appErr
	}
	return success(nil), nil
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
	if !isBlank(req.AdvisorEmail) && !isValidEmail(req.AdvisorEmail) {
		return SuccessResponse{}, NewBadRequest("advisor_email is invalid")
	}

	if _, appErr := s.createRequest(student.ID, "GRADUATION_THESIS", req); appErr != nil {
		return SuccessResponse{}, appErr
	}
	return success(nil), nil
}

func (s *Service) GetDeadlines(studentID string) (SuccessResponse, *AppError) {
	student, err := s.getStudent(studentID)
	if err != nil {
		return SuccessResponse{}, err
	}
	deadlines, repoErr := s.deadlines.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]models.Deadline, 0)
	for _, deadline := range deadlines {
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
	_, ok, repoErr := s.courses.FindByID(courseID)
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if !ok {
		return SuccessResponse{}, NewNotFound("course not found")
	}
	materials, repoErr := s.materials.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]models.Material, 0)
	for _, material := range materials {
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
	_, ok, repoErr := s.courses.FindByID(courseID)
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if !ok {
		return SuccessResponse{}, NewNotFound("course not found")
	}
	assignments, repoErr := s.assignments.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}

	result := make([]models.Assignment, 0)
	for _, assignment := range assignments {
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
	assignment, ok, repoErr := s.assignments.FindByID(assignmentID)
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if !ok {
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
		if !isValidURL(req.URL) {
			return SuccessResponse{}, NewBadRequest("url is invalid")
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
	_ = uuid.NewString()
	count, repoErr := s.submissions.Count()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if err := s.submissions.Create(models.Submission{
		ID:             uint(count + 1),
		AssignmentID:   assignment.ID,
		StudentID:      student.ID,
		SubmissionType: submissionType,
		Content:        content,
		CreatedAt:      s.now().Format(time.RFC3339),
	}); err != nil {
		return SuccessResponse{}, NewInternal("")
	}

	return success(nil), nil
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
	if req.Reason == "OTHER" {
		if !strings.HasPrefix(strings.TrimSpace(req.OtherReason), "Bổ sung hồ sơ") {
			return SuccessResponse{}, NewBadRequest("other_reason is invalid")
		}
	} else {
		if !isBlank(req.OtherReason) {
			return SuccessResponse{}, NewBadRequest("other_reason is not allowed")
		}
	}

	created, appErr := s.createRequest(student.ID, "CONFIRM_LETTER", req)
	if appErr != nil {
		return SuccessResponse{}, appErr
	}
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

	created, appErr := s.createRequest(student.ID, "BANK_LOANS", req)
	if appErr != nil {
		return SuccessResponse{}, appErr
	}
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

	created, appErr := s.createRequest(student.ID, "TRAINING_POINT_CONFIRM", req)
	if appErr != nil {
		return SuccessResponse{}, appErr
	}
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
	if len(strings.TrimSpace(req.IDNumber)) < 6 {
		return SuccessResponse{}, NewBadRequest("id_number is invalid")
	}
	if !isValidBase64(req.ImageFile) {
		return SuccessResponse{}, NewBadRequest("image_file is invalid")
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

	created, appErr := s.createRequest(student.ID, "LANGUAGE_CERTIFICATE", req)
	if appErr != nil {
		return SuccessResponse{}, appErr
	}
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
	rooms, repoErr := s.rooms.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	bookings, repoErr := s.roomBookings.List()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	for _, room := range rooms {
		if isRoomAvailable(room.ID, query.Date, startMin, endMin, bookings) {
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
	if !isValidEmail(req.Email) {
		return SuccessResponse{}, NewBadRequest("email is invalid")
	}

	_ = uuid.NewString()
	count, repoErr := s.contacts.Count()
	if repoErr != nil {
		return SuccessResponse{}, NewInternal("")
	}
	if err := s.contacts.Create(models.Contact{
		ID:        uint(count + 1),
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Subject:   req.Subject,
		Message:   req.Message,
		CreatedAt: s.now().Format(time.RFC3339),
	}); err != nil {
		return SuccessResponse{}, NewInternal("")
	}

	return success(nil), nil
}

func (s *Service) createRequest(studentID, requestType string, payload interface{}) (RequestStatusData, *AppError) {
	requestID := uuid.NewString()
	payloadBytes, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return RequestStatusData{}, NewInternal("")
	}
	if err := s.requests.Create(models.Request{
		ID:        requestID,
		StudentID: studentID,
		Type:      requestType,
		Status:    "PENDING",
		Payload:   string(payloadBytes),
		CreatedAt: s.now().Format(time.RFC3339),
	}); err != nil {
		return RequestStatusData{}, NewInternal("")
	}

	return RequestStatusData{
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

func isRoomAvailable(roomID uint, date string, startMin, endMin int, bookings []models.RoomBooking) bool {
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
}
