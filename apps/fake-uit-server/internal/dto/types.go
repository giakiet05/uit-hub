package dto



type AppError struct {
	Status  int
	Code    string
	Message string
}

func NewBadRequest(message string) *AppError {
	return &AppError{Status: 400, Code: "VALIDATION_ERROR", Message: message}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: message}
}

func NewNotFound(message string) *AppError {
	return &AppError{Status: 404, Code: "NOT_FOUND", Message: message}
}

func NewInternal(message string) *AppError {
	if message == "" {
		message = "Internal error happened!"
	}
	return &AppError{Status: 500, Code: "INTERNAL_ERROR", Message: message}
}

type QueryYearSemester struct {
	Year     int `json:"year" form:"year"`
	Semester int `json:"semester" form:"semester"`
}

type QueryRoomsAvailability struct {
	Date  string `json:"date" form:"date"`
	Start string `json:"start" form:"start"`
	End   string `json:"end" form:"end"`
}


type TranscriptRegisRequest struct {
	Copies          int    `json:"copies"`
	Language        string `json:"language"`
	DeliveryMethod  string `json:"delivery_method"`
	ShippingAddress string `json:"shipping_address"`
	Phone           string `json:"phone"`
	Note            string `json:"note"`
}

type TuitionExtendRequest struct {
	Year             int      `json:"year"`
	Semester         int      `json:"semester"`
	RequestedDueDate string   `json:"requested_due_date"`
	Reason           string   `json:"reason"`
	Phone            string   `json:"phone"`
	Attachments      []string `json:"attachments"`
}

type MonthlyParkingRequest struct {
	VehicleType string `json:"vehicle_type"`
	PlateNumber string `json:"plate_number"`
	Months      int    `json:"months"`
	StartMonth  string `json:"start_month"`
	OwnerName   string `json:"owner_name"`
	Note        string `json:"note"`
}

type GraduateRequest struct {
	Year     int    `json:"year"`
	Semester int    `json:"semester"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Note     string `json:"note"`
}

type TeamMember struct {
	StudentID string `json:"student_id"`
	FullName  string `json:"full_name"`
}

type GraduationThesisRequest struct {
	ThesisTitle  string       `json:"thesis_title"`
	AdvisorName  string       `json:"advisor_name"`
	AdvisorEmail string       `json:"advisor_email"`
	TeamMembers  []TeamMember `json:"team_members"`
	Note         string       `json:"note"`
}

type AssignmentSubmissionRequest struct {
	SubmissionType string   `json:"submission_type"`
	Text           string   `json:"text"`
	URL            string   `json:"url"`
	FileURLs       []string `json:"file_urls"`
	Comment        string   `json:"comment"`
}

type ConfirmLetterRequest struct {
	Language    string `json:"language"`
	Reason      string `json:"reason"`
	OtherReason string `json:"other_reason"`
	RequestType string `json:"request_type"`
	Note        string `json:"note"`
}

type BankLoansRequest struct {
	Benefit      string `json:"benefit"`
	OrphanStatus string `json:"orphan_status"`
	Template     string `json:"template"`
	Note         string `json:"note"`
}

type TrainingPointConfirmRequest struct {
	Language string `json:"language"`
	Note     string `json:"note"`
}

type LanguageCertificateUploadRequest struct {
	DocumentType   string  `json:"document_type"`
	BirthDate      string  `json:"birth_date"`
	IDNumber       string  `json:"id_number"`
	ListeningScore float64 `json:"listening_score"`
	ReadingScore   float64 `json:"reading_score"`
	TotalScore     float64 `json:"total_score"`
	ExamDate       string  `json:"exam_date"`
	ImageFile      string  `json:"image_file"`
}

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type StudentProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Major string `json:"major"`
}

type TuitionFeeData struct {
	TotalCredits int    `json:"total_credits"`
	Amount       int    `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
}

type InsuranceData struct {
	Provider  string `json:"provider"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at"`
}

type TrainingPointsData struct {
	Year     int    `json:"year"`
	Semester int    `json:"semester"`
	Score    int    `json:"score"`
	Rank     string `json:"rank"`
}

type SubmissionResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
}

type RequestStatusData struct {
	RequestID string  `json:"request_id"`
	Status    string  `json:"status"`
	CreatedAt int64   `json:"created_at"`
	PDFURL    *string `json:"pdf_url"`
}

type RoomAvailabilityData struct {
	Date  string     `json:"date"`
	Start string     `json:"start"`
	End   string     `json:"end"`
	Rooms []RoomItem `json:"rooms"`
}

type RoomItem struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
