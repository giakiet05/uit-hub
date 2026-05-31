package model

type Student struct {
	ID       string ` json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Major    string `json:"major"`
}

type Course struct {
	ID      string ` json:"id"`
	Name    string `json:"name"`
	Credits int    `json:"credits"`
}

type Enrollment struct {
	ID        uint   ` json:"id"`
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
	Year      int    `json:"year"`
	Semester  int    `json:"semester"`
}

type Schedule struct {
	ID        uint   ` json:"id"`
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Room      string `json:"room"`
	Year      int    `json:"year"`
	Semester  int    `json:"semester"`
}

type ExamSchedule struct {
	ID        uint   ` json:"id"`
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
	Date      string `json:"date"`
	Room      string `json:"room"`
	Year      int    `json:"year"`
	Semester  int    `json:"semester"`
}

type Score struct {
	ID        uint    ` json:"id"`
	StudentID string  `json:"student_id"`
	CourseID  string  `json:"course_id"`
	Midterm   float64 `json:"midterm"`
	Final     float64 `json:"final"`
	Total     float64 `json:"total"`
}

type Assignment struct {
	ID       string ` json:"id"`
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
	Deadline string `json:"deadline"`
}

type Submission struct {
	ID             uint   ` json:"id"`
	AssignmentID   string `json:"assignment_id"`
	StudentID      string `json:"student_id"`
	SubmissionType string `json:"submission_type"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
}

type Material struct {
	ID       uint   ` json:"id"`
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type Deadline struct {
	ID        uint   ` json:"id"`
	StudentID string `json:"student_id"`
	Title     string `json:"title"`
	DueDate   string `json:"due_date"`
}

type Request struct {
	ID        string ` json:"id"`
	StudentID string `json:"student_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Payload   string `json:"payload"`
	CreatedAt string `json:"created_at"`
}

type Room struct {
	ID       uint   ` json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type RoomBooking struct {
	ID        uint   ` json:"id"`
	RoomID    uint   `json:"room_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Contact struct {
	ID        uint   ` json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
