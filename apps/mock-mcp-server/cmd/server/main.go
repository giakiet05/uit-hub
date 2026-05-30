package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverVersion = "0.1.0"

type store struct {
	mu                   sync.RWMutex
	students             map[string]studentProfile
	courses              []course
	schedules            map[string]studentSchedule
	tuition              map[string]tuitionStatus
	notes                map[string]note
	leaveRequests        map[string]leaveRequest
	studyPlans           map[string]studyPlan
	advisorTickets       map[string]advisorTicket
	supportRequests      map[string]supportRequest
	nextNoteID           int
	nextLeaveID          int
	nextStudyPlanID      int
	nextTicketID         int
	nextSupportRequestID int
}

type studentProfile struct {
	StudentID       string  `json:"student_id" jsonschema:"student ID"`
	FullName        string  `json:"full_name" jsonschema:"student full name"`
	Major           string  `json:"major" jsonschema:"student major"`
	GPA             float64 `json:"gpa" jsonschema:"current GPA"`
	CompletedCredit int     `json:"completed_credits" jsonschema:"completed credits"`
	AcademicWarning bool    `json:"academic_warning" jsonschema:"whether the student is under academic warning"`
}

type course struct {
	Code          string   `json:"code" jsonschema:"course code"`
	Name          string   `json:"name" jsonschema:"course name"`
	Credits       int      `json:"credits" jsonschema:"course credit count"`
	Tags          []string `json:"tags" jsonschema:"course tags"`
	RiskLevel     string   `json:"risk_level" jsonschema:"risk level: low, medium, or high"`
	RiskReason    string   `json:"risk_reason" jsonschema:"why this course may be risky"`
	Prerequisites []string `json:"prerequisites" jsonschema:"course prerequisite codes"`
}

type scheduleItem struct {
	CourseCode string `json:"course_code" jsonschema:"course code"`
	Day        string `json:"day" jsonschema:"weekday"`
	StartTime  string `json:"start_time" jsonschema:"start time in HH:MM"`
	EndTime    string `json:"end_time" jsonschema:"end time in HH:MM"`
	Room       string `json:"room" jsonschema:"classroom"`
}

type studentSchedule struct {
	StudentID string         `json:"student_id" jsonschema:"student ID"`
	Term      string         `json:"term" jsonschema:"academic term"`
	Items     []scheduleItem `json:"items" jsonschema:"schedule items"`
}

type note struct {
	NoteID     string   `json:"note_id" jsonschema:"note ID"`
	StudentID  string   `json:"student_id" jsonschema:"student ID"`
	Body       string   `json:"body" jsonschema:"note body"`
	Labels     []string `json:"labels" jsonschema:"note labels"`
	Visibility string   `json:"visibility" jsonschema:"note visibility"`
	UpdatedAt  string   `json:"updated_at" jsonschema:"RFC3339 update timestamp"`
}

type leaveRequest struct {
	RequestID string    `json:"request_id" jsonschema:"leave request ID"`
	StudentID string    `json:"student_id" jsonschema:"student ID"`
	Reason    string    `json:"reason" jsonschema:"leave reason"`
	Period    dateRange `json:"period" jsonschema:"requested leave period"`
	Contacts  []contact `json:"contacts" jsonschema:"emergency or notification contacts"`
	Status    string    `json:"status" jsonschema:"request status"`
	CreatedAt string    `json:"created_at" jsonschema:"RFC3339 creation timestamp"`
}

type tuitionStatus struct {
	StudentID string  `json:"student_id" jsonschema:"student ID"`
	Status    string  `json:"status" jsonschema:"paid, pending, or overdue"`
	Balance   float64 `json:"balance" jsonschema:"remaining tuition balance"`
	Currency  string  `json:"currency" jsonschema:"currency code"`
	DueDate   string  `json:"due_date" jsonschema:"tuition due date in YYYY-MM-DD"`
}

type studyPlan struct {
	PlanID    string   `json:"study_plan_id" jsonschema:"study plan ID"`
	StudentID string   `json:"student_id" jsonschema:"student ID"`
	Title     string   `json:"title" jsonschema:"study plan title"`
	Goals     []string `json:"goals" jsonschema:"study plan goals"`
	Notes     []string `json:"notes" jsonschema:"study plan notes"`
	Status    string   `json:"status" jsonschema:"draft, active, or archived"`
	UpdatedAt string   `json:"updated_at" jsonschema:"RFC3339 update timestamp"`
}

type advisorTicket struct {
	TicketID    string   `json:"ticket_id" jsonschema:"advisor ticket ID"`
	StudentID   string   `json:"student_id" jsonschema:"student ID"`
	StudyPlanID string   `json:"study_plan_id,omitempty" jsonschema:"related study plan ID"`
	Summary     string   `json:"summary" jsonschema:"ticket summary"`
	Labels      []string `json:"labels" jsonschema:"ticket labels"`
	Status      string   `json:"status" jsonschema:"open, pending, or closed"`
	UpdatedAt   string   `json:"updated_at" jsonschema:"RFC3339 update timestamp"`
}

type supportRequest struct {
	RequestID string   `json:"support_request_id" jsonschema:"support request ID"`
	StudentID string   `json:"student_id" jsonschema:"student ID"`
	Reason    string   `json:"reason" jsonschema:"support request reason"`
	Channels  []string `json:"channels" jsonschema:"requested support channels"`
	Status    string   `json:"status" jsonschema:"submitted, triaged, or closed"`
	CreatedAt string   `json:"created_at" jsonschema:"RFC3339 creation timestamp"`
}

type dateRange struct {
	StartDate string `json:"start_date" jsonschema:"start date in YYYY-MM-DD"`
	EndDate   string `json:"end_date" jsonschema:"end date in YYYY-MM-DD"`
}

type contact struct {
	Name  string `json:"name" jsonschema:"contact name"`
	Email string `json:"email" jsonschema:"contact email"`
}

type getStudentProfileInput struct {
	StudentID string `json:"student_id" jsonschema:"student ID, for example 22520001"`
}

type getStudentProfileOutput struct {
	Found   bool            `json:"found" jsonschema:"whether a profile was found"`
	Profile *studentProfile `json:"profile,omitempty" jsonschema:"student profile when found"`
}

type searchCoursesInput struct {
	Query string   `json:"query,omitempty" jsonschema:"course code, course name, or tag text"`
	Tags  []string `json:"tags,omitempty" jsonschema:"tags that must be present on the course"`
	Limit int      `json:"limit,omitempty" jsonschema:"maximum number of courses to return"`
}

type searchCoursesOutput struct {
	Count   int      `json:"count" jsonschema:"number of returned courses"`
	Courses []course `json:"courses" jsonschema:"matching courses"`
}

type getStudentScheduleInput struct {
	StudentID string `json:"student_id" jsonschema:"student ID"`
	Term      string `json:"term,omitempty" jsonschema:"academic term; omit to use current term"`
}

type getStudentScheduleOutput struct {
	Found    bool             `json:"found" jsonschema:"whether a schedule was found"`
	Schedule *studentSchedule `json:"schedule,omitempty" jsonschema:"student schedule when found"`
}

type getCourseDetailInput struct {
	CourseCode string `json:"course_code" jsonschema:"course code"`
}

type getCourseDetailOutput struct {
	Found  bool    `json:"found" jsonschema:"whether a course was found"`
	Course *course `json:"course,omitempty" jsonschema:"course detail when found"`
}

type getTuitionStatusInput struct {
	StudentID string `json:"student_id" jsonschema:"student ID"`
}

type getTuitionStatusOutput struct {
	Found  bool           `json:"found" jsonschema:"whether tuition status was found"`
	Status *tuitionStatus `json:"status,omitempty" jsonschema:"tuition status when found"`
}

type upsertStudentNoteInput struct {
	StudentID  string   `json:"student_id" jsonschema:"student ID"`
	NoteID     string   `json:"note_id,omitempty" jsonschema:"existing note ID; omit to create a new note"`
	Body       string   `json:"body" jsonschema:"note body"`
	Labels     []string `json:"labels,omitempty" jsonschema:"labels to attach to the note"`
	Visibility string   `json:"visibility,omitempty" jsonschema:"visibility value: private, advisor, or public"`
}

type upsertStudentNoteOutput struct {
	Created bool `json:"created" jsonschema:"whether a new note was created"`
	Note    note `json:"note" jsonschema:"stored note"`
}

type listStudentNotesInput struct {
	StudentID string   `json:"student_id" jsonschema:"student ID"`
	Labels    []string `json:"labels,omitempty" jsonschema:"labels that must be present on the note"`
}

type listStudentNotesOutput struct {
	Count int    `json:"count" jsonschema:"number of returned notes"`
	Notes []note `json:"notes" jsonschema:"matching notes"`
}

type deleteStudentNoteInput struct {
	NoteID string `json:"note_id" jsonschema:"note ID to delete"`
}

type deleteStudentNoteOutput struct {
	Deleted bool   `json:"deleted" jsonschema:"whether a note was deleted"`
	NoteID  string `json:"note_id" jsonschema:"note ID requested for deletion"`
}

type submitLeaveRequestInput struct {
	StudentID string    `json:"student_id" jsonschema:"student ID"`
	Reason    string    `json:"reason" jsonschema:"leave request reason"`
	Period    dateRange `json:"period" jsonschema:"requested leave period"`
	Contacts  []contact `json:"contacts,omitempty" jsonschema:"emergency or notification contacts"`
}

type submitLeaveRequestOutput struct {
	Request leaveRequest `json:"request" jsonschema:"created leave request"`
}

type createStudyPlanInput struct {
	StudentID string   `json:"student_id" jsonschema:"student ID"`
	Title     string   `json:"title" jsonschema:"study plan title"`
	Goals     []string `json:"goals" jsonschema:"study goals"`
	Notes     []string `json:"notes,omitempty" jsonschema:"study plan notes"`
}

type createStudyPlanOutput struct {
	Plan studyPlan `json:"plan" jsonschema:"created study plan"`
}

type patchStudyPlanInput struct {
	PlanID string   `json:"study_plan_id" jsonschema:"study plan ID"`
	Goals  []string `json:"goals,omitempty" jsonschema:"goals to replace when provided"`
	Notes  []string `json:"notes,omitempty" jsonschema:"notes to append"`
	Status string   `json:"status,omitempty" jsonschema:"new plan status"`
}

type patchStudyPlanOutput struct {
	Found bool       `json:"found" jsonschema:"whether the study plan was found"`
	Plan  *studyPlan `json:"plan,omitempty" jsonschema:"updated study plan when found"`
}

type createAdvisorTicketInput struct {
	StudentID   string   `json:"student_id" jsonschema:"student ID"`
	StudyPlanID string   `json:"study_plan_id,omitempty" jsonschema:"related study plan ID"`
	Summary     string   `json:"summary" jsonschema:"ticket summary"`
	Labels      []string `json:"labels,omitempty" jsonschema:"ticket labels"`
}

type createAdvisorTicketOutput struct {
	Ticket advisorTicket `json:"ticket" jsonschema:"created advisor ticket"`
}

type patchAdvisorTicketInput struct {
	TicketID string   `json:"ticket_id" jsonschema:"advisor ticket ID"`
	Summary  string   `json:"summary,omitempty" jsonschema:"replacement ticket summary"`
	Labels   []string `json:"labels,omitempty" jsonschema:"labels to replace when provided"`
	Status   string   `json:"status,omitempty" jsonschema:"new ticket status"`
}

type patchAdvisorTicketOutput struct {
	Found  bool           `json:"found" jsonschema:"whether the advisor ticket was found"`
	Ticket *advisorTicket `json:"ticket,omitempty" jsonschema:"updated advisor ticket when found"`
}

type submitSupportRequestInput struct {
	StudentID string   `json:"student_id" jsonschema:"student ID"`
	Reason    string   `json:"reason" jsonschema:"support request reason"`
	Channels  []string `json:"channels,omitempty" jsonschema:"support channels such as email, phone, or advisor"`
}

type submitSupportRequestOutput struct {
	Request supportRequest `json:"request" jsonschema:"created support request"`
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "uit-hub-mock-mcp-server",
		Version: serverVersion,
	}, nil)

	data := newStore()

	delayStr := os.Getenv("MOCK_UIT_DELAY")
	var delay time.Duration
	if delayStr != "" {
		if d, err := time.ParseDuration(delayStr); err == nil {
			delay = d
		} else {
			log.Printf("invalid MOCK_UIT_DELAY %q: %v", delayStr, err)
		}
	}

	registerTools(server, data, delay)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("mock mcp server failed: %v", err)
	}
}

func withDelay[T any, U any](
	delay time.Duration,
	handler func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, U, error),
) func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, U, error) {
	if delay <= 0 {
		return handler
	}
	return func(ctx context.Context, req *mcp.CallToolRequest, input T) (*mcp.CallToolResult, U, error) {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			var empty U
			return nil, empty, ctx.Err()
		}
		return handler(ctx, req, input)
	}
}

func safeSchema[In any]() *jsonschema.Schema {
	schema, err := jsonschema.ForType(reflect.TypeFor[In](), nil)
	if err != nil {
		panic(err)
	}
	schema.Extra = map[string]any{"_concurrencySafe": true}
	return schema
}

func writeSchema[In any]() *jsonschema.Schema {
	schema, err := jsonschema.ForType(reflect.TypeFor[In](), nil)
	if err != nil {
		panic(err)
	}
	schema.Extra = map[string]any{"_requireApproval": true}
	return schema
}

func registerTools(server *mcp.Server, data *store, delay time.Duration) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_student_profile",
		Description: "Read a mock student profile by student ID.",
		InputSchema: safeSchema[getStudentProfileInput](),
	}, withDelay(delay, data.getStudentProfile))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_courses",
		Description: "Read mock course catalog entries using query text and tags.",
		InputSchema: safeSchema[searchCoursesInput](),
	}, withDelay(delay, data.searchCourses))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_student_schedule",
		Description: "Read a mock student schedule for the current or requested term.",
		InputSchema: safeSchema[getStudentScheduleInput](),
	}, withDelay(delay, data.getStudentSchedule))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_course_detail",
		Description: "Read detailed mock course information by course code.",
		InputSchema: safeSchema[getCourseDetailInput](),
	}, withDelay(delay, data.getCourseDetail))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_tuition_status",
		Description: "Read mock tuition status for a student.",
		InputSchema: safeSchema[getTuitionStatusInput](),
	}, withDelay(delay, data.getTuitionStatus))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "upsert_student_note",
		Description: "Create or update an advisor note for a student. This is a write operation.",
		InputSchema: writeSchema[upsertStudentNoteInput](),
	}, withDelay(delay, data.upsertStudentNote))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_student_notes",
		Description: "Read advisor notes for a student.",
		InputSchema: safeSchema[listStudentNotesInput](),
	}, withDelay(delay, data.listStudentNotes))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_student_note",
		Description: "Delete an advisor note by note ID. This is a write operation.",
		InputSchema: writeSchema[deleteStudentNoteInput](),
	}, withDelay(delay, data.deleteStudentNote))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "submit_leave_request",
		Description: "Create a mock leave request for a student. This is a write operation with nested object and array input.",
		InputSchema: writeSchema[submitLeaveRequestInput](),
	}, withDelay(delay, data.submitLeaveRequest))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_study_plan",
		Description: "Create a mock study plan for a student. This is a write operation.",
		InputSchema: writeSchema[createStudyPlanInput](),
	}, withDelay(delay, data.createStudyPlan))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "patch_study_plan",
		Description: "Patch a mock study plan by ID. This is a write operation.",
		InputSchema: writeSchema[patchStudyPlanInput](),
	}, withDelay(delay, data.patchStudyPlan))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_advisor_ticket",
		Description: "Create a mock advisor ticket for a student. This is a write operation.",
		InputSchema: writeSchema[createAdvisorTicketInput](),
	}, withDelay(delay, data.createAdvisorTicket))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "patch_advisor_ticket",
		Description: "Patch a mock advisor ticket by ID. This is a write operation.",
		InputSchema: writeSchema[patchAdvisorTicketInput](),
	}, withDelay(delay, data.patchAdvisorTicket))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "submit_support_request",
		Description: "Submit a mock student support request. This is a write operation.",
		InputSchema: writeSchema[submitSupportRequestInput](),
	}, withDelay(delay, data.submitSupportRequest))
}

func newStore() *store {
	return &store{
		students: map[string]studentProfile{
			"22520001": {
				StudentID:       "22520001",
				FullName:        "Nguyen Van An",
				Major:           "Computer Science",
				GPA:             3.42,
				CompletedCredit: 92,
				AcademicWarning: false,
			},
			"22520002": {
				StudentID:       "22520002",
				FullName:        "Tran Thi Binh",
				Major:           "Information Systems",
				GPA:             1.83,
				CompletedCredit: 48,
				AcademicWarning: true,
			},
		},
		courses: []course{
			{
				Code:          "CS101",
				Name:          "Introduction to Programming",
				Credits:       4,
				Tags:          []string{"programming", "foundation"},
				RiskLevel:     "low",
				RiskReason:    "Foundation course with moderate workload.",
				Prerequisites: []string{},
			},
			{
				Code:          "CS201",
				Name:          "Data Structures and Algorithms",
				Credits:       4,
				Tags:          []string{"programming", "algorithm"},
				RiskLevel:     "high",
				RiskReason:    "High workload and requires consistent practice.",
				Prerequisites: []string{"CS101"},
			},
			{
				Code:          "IS301",
				Name:          "Database Systems",
				Credits:       3,
				Tags:          []string{"database", "information_systems"},
				RiskLevel:     "medium",
				RiskReason:    "Project and SQL practice required.",
				Prerequisites: []string{"CS101"},
			},
			{
				Code:          "SE401",
				Name:          "Software Architecture",
				Credits:       3,
				Tags:          []string{"software_engineering", "architecture"},
				RiskLevel:     "high",
				RiskReason:    "Requires design experience and team coordination.",
				Prerequisites: []string{"CS201"},
			},
		},
		schedules: map[string]studentSchedule{
			"22520001": {
				StudentID: "22520001",
				Term:      "2026-1",
				Items: []scheduleItem{
					{CourseCode: "IS301", Day: "Monday", StartTime: "08:00", EndTime: "10:30", Room: "B1.204"},
					{CourseCode: "SE401", Day: "Wednesday", StartTime: "13:00", EndTime: "15:30", Room: "C.105"},
					{CourseCode: "CS201", Day: "Friday", StartTime: "09:00", EndTime: "11:30", Room: "A2.301"},
				},
			},
			"22520002": {
				StudentID: "22520002",
				Term:      "2026-1",
				Items: []scheduleItem{
					{CourseCode: "CS201", Day: "Tuesday", StartTime: "07:30", EndTime: "10:00", Room: "B2.101"},
					{CourseCode: "SE401", Day: "Thursday", StartTime: "13:00", EndTime: "15:30", Room: "C.205"},
					{CourseCode: "IS301", Day: "Saturday", StartTime: "08:00", EndTime: "10:30", Room: "A1.401"},
				},
			},
		},
		tuition: map[string]tuitionStatus{
			"22520001": {
				StudentID: "22520001",
				Status:    "paid",
				Balance:   0,
				Currency:  "VND",
				DueDate:   "2026-06-15",
			},
			"22520002": {
				StudentID: "22520002",
				Status:    "pending",
				Balance:   4200000,
				Currency:  "VND",
				DueDate:   "2026-06-15",
			},
		},
		notes:                map[string]note{},
		leaveRequests:        map[string]leaveRequest{},
		studyPlans:           map[string]studyPlan{},
		advisorTickets:       map[string]advisorTicket{},
		supportRequests:      map[string]supportRequest{},
		nextNoteID:           1,
		nextLeaveID:          1,
		nextStudyPlanID:      1,
		nextTicketID:         1,
		nextSupportRequestID: 1,
	}
}

func (s *store) getStudentProfile(ctx context.Context, req *mcp.CallToolRequest, input getStudentProfileInput) (*mcp.CallToolResult, getStudentProfileOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, getStudentProfileOutput{}, err
	}

	s.mu.RLock()
	profile, found := s.students[strings.TrimSpace(input.StudentID)]
	s.mu.RUnlock()

	if !found {
		return nil, getStudentProfileOutput{Found: false}, nil
	}
	return nil, getStudentProfileOutput{Found: true, Profile: &profile}, nil
}

func (s *store) searchCourses(ctx context.Context, req *mcp.CallToolRequest, input searchCoursesInput) (*mcp.CallToolResult, searchCoursesOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, searchCoursesOutput{}, err
	}

	query := strings.ToLower(strings.TrimSpace(input.Query))
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]course, 0, len(s.courses))
	for _, item := range s.courses {
		if query != "" && !courseMatchesQuery(item, query) {
			continue
		}
		if !courseHasTags(item, input.Tags) {
			continue
		}
		results = append(results, item)
		if len(results) >= limit {
			break
		}
	}

	return nil, searchCoursesOutput{
		Count:   len(results),
		Courses: results,
	}, nil
}

func (s *store) getStudentSchedule(ctx context.Context, req *mcp.CallToolRequest, input getStudentScheduleInput) (*mcp.CallToolResult, getStudentScheduleOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, getStudentScheduleOutput{}, err
	}

	studentID := strings.TrimSpace(input.StudentID)
	s.mu.RLock()
	schedule, found := s.schedules[studentID]
	s.mu.RUnlock()
	if !found {
		return nil, getStudentScheduleOutput{Found: false}, nil
	}

	term := strings.TrimSpace(input.Term)
	if term != "" && schedule.Term != term {
		return nil, getStudentScheduleOutput{Found: false}, nil
	}
	return nil, getStudentScheduleOutput{Found: true, Schedule: &schedule}, nil
}

func (s *store) getCourseDetail(ctx context.Context, req *mcp.CallToolRequest, input getCourseDetailInput) (*mcp.CallToolResult, getCourseDetailOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, getCourseDetailOutput{}, err
	}

	courseCode := strings.TrimSpace(input.CourseCode)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.courses {
		if strings.EqualFold(item.Code, courseCode) {
			course := item
			return nil, getCourseDetailOutput{Found: true, Course: &course}, nil
		}
	}
	return nil, getCourseDetailOutput{Found: false}, nil
}

func (s *store) getTuitionStatus(ctx context.Context, req *mcp.CallToolRequest, input getTuitionStatusInput) (*mcp.CallToolResult, getTuitionStatusOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, getTuitionStatusOutput{}, err
	}

	studentID := strings.TrimSpace(input.StudentID)
	s.mu.RLock()
	status, found := s.tuition[studentID]
	s.mu.RUnlock()
	if !found {
		return nil, getTuitionStatusOutput{Found: false}, nil
	}
	return nil, getTuitionStatusOutput{Found: true, Status: &status}, nil
}

func (s *store) upsertStudentNote(ctx context.Context, req *mcp.CallToolRequest, input upsertStudentNoteInput) (*mcp.CallToolResult, upsertStudentNoteOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, upsertStudentNoteOutput{}, err
	}
	if strings.TrimSpace(input.StudentID) == "" {
		return nil, upsertStudentNoteOutput{}, fmt.Errorf("student_id is required")
	}
	if strings.TrimSpace(input.Body) == "" {
		return nil, upsertStudentNoteOutput{}, fmt.Errorf("body is required")
	}

	visibility := strings.TrimSpace(input.Visibility)
	if visibility == "" {
		visibility = "private"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	created := false
	noteID := strings.TrimSpace(input.NoteID)
	if noteID == "" {
		noteID = fmt.Sprintf("note-%03d", s.nextNoteID)
		s.nextNoteID++
		created = true
	} else if _, exists := s.notes[noteID]; !exists {
		created = true
	}

	stored := note{
		NoteID:     noteID,
		StudentID:  strings.TrimSpace(input.StudentID),
		Body:       input.Body,
		Labels:     append([]string{}, input.Labels...),
		Visibility: visibility,
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	s.notes[noteID] = stored

	return nil, upsertStudentNoteOutput{
		Created: created,
		Note:    stored,
	}, nil
}

func (s *store) listStudentNotes(ctx context.Context, req *mcp.CallToolRequest, input listStudentNotesInput) (*mcp.CallToolResult, listStudentNotesOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, listStudentNotesOutput{}, err
	}

	studentID := strings.TrimSpace(input.StudentID)

	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]note, 0)
	for _, item := range s.notes {
		if item.StudentID != studentID {
			continue
		}
		if !noteHasLabels(item, input.Labels) {
			continue
		}
		results = append(results, item)
	}

	return nil, listStudentNotesOutput{
		Count: len(results),
		Notes: results,
	}, nil
}

func (s *store) deleteStudentNote(ctx context.Context, req *mcp.CallToolRequest, input deleteStudentNoteInput) (*mcp.CallToolResult, deleteStudentNoteOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, deleteStudentNoteOutput{}, err
	}

	noteID := strings.TrimSpace(input.NoteID)
	if noteID == "" {
		return nil, deleteStudentNoteOutput{}, fmt.Errorf("note_id is required")
	}

	s.mu.Lock()
	_, exists := s.notes[noteID]
	if exists {
		delete(s.notes, noteID)
	}
	s.mu.Unlock()

	return nil, deleteStudentNoteOutput{
		Deleted: exists,
		NoteID:  noteID,
	}, nil
}

func (s *store) submitLeaveRequest(ctx context.Context, req *mcp.CallToolRequest, input submitLeaveRequestInput) (*mcp.CallToolResult, submitLeaveRequestOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, submitLeaveRequestOutput{}, err
	}
	if strings.TrimSpace(input.StudentID) == "" {
		return nil, submitLeaveRequestOutput{}, fmt.Errorf("student_id is required")
	}
	if strings.TrimSpace(input.Reason) == "" {
		return nil, submitLeaveRequestOutput{}, fmt.Errorf("reason is required")
	}
	if strings.TrimSpace(input.Period.StartDate) == "" || strings.TrimSpace(input.Period.EndDate) == "" {
		return nil, submitLeaveRequestOutput{}, fmt.Errorf("period start_date and end_date are required")
	}

	s.mu.Lock()
	requestID := fmt.Sprintf("leave-%03d", s.nextLeaveID)
	s.nextLeaveID++
	request := leaveRequest{
		RequestID: requestID,
		StudentID: strings.TrimSpace(input.StudentID),
		Reason:    input.Reason,
		Period:    input.Period,
		Contacts:  append([]contact{}, input.Contacts...),
		Status:    "submitted",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.leaveRequests[requestID] = request
	s.mu.Unlock()

	return nil, submitLeaveRequestOutput{Request: request}, nil
}

func (s *store) createStudyPlan(ctx context.Context, req *mcp.CallToolRequest, input createStudyPlanInput) (*mcp.CallToolResult, createStudyPlanOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, createStudyPlanOutput{}, err
	}
	if strings.TrimSpace(input.StudentID) == "" {
		return nil, createStudyPlanOutput{}, fmt.Errorf("student_id is required")
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, createStudyPlanOutput{}, fmt.Errorf("title is required")
	}

	s.mu.Lock()
	planID := fmt.Sprintf("plan-%03d", s.nextStudyPlanID)
	s.nextStudyPlanID++
	plan := studyPlan{
		PlanID:    planID,
		StudentID: strings.TrimSpace(input.StudentID),
		Title:     input.Title,
		Goals:     append([]string{}, input.Goals...),
		Notes:     append([]string{}, input.Notes...),
		Status:    "draft",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.studyPlans[planID] = plan
	s.mu.Unlock()

	return nil, createStudyPlanOutput{Plan: plan}, nil
}

func (s *store) patchStudyPlan(ctx context.Context, req *mcp.CallToolRequest, input patchStudyPlanInput) (*mcp.CallToolResult, patchStudyPlanOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, patchStudyPlanOutput{}, err
	}

	planID := strings.TrimSpace(input.PlanID)
	if planID == "" {
		return nil, patchStudyPlanOutput{}, fmt.Errorf("study_plan_id is required")
	}

	s.mu.Lock()
	plan, found := s.studyPlans[planID]
	if found {
		if len(input.Goals) > 0 {
			plan.Goals = append([]string{}, input.Goals...)
		}
		if len(input.Notes) > 0 {
			plan.Notes = append(plan.Notes, input.Notes...)
		}
		if strings.TrimSpace(input.Status) != "" {
			plan.Status = strings.TrimSpace(input.Status)
		}
		plan.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		s.studyPlans[planID] = plan
	}
	s.mu.Unlock()

	if !found {
		return nil, patchStudyPlanOutput{Found: false}, nil
	}
	return nil, patchStudyPlanOutput{Found: true, Plan: &plan}, nil
}

func (s *store) createAdvisorTicket(ctx context.Context, req *mcp.CallToolRequest, input createAdvisorTicketInput) (*mcp.CallToolResult, createAdvisorTicketOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, createAdvisorTicketOutput{}, err
	}
	if strings.TrimSpace(input.StudentID) == "" {
		return nil, createAdvisorTicketOutput{}, fmt.Errorf("student_id is required")
	}
	if strings.TrimSpace(input.Summary) == "" {
		return nil, createAdvisorTicketOutput{}, fmt.Errorf("summary is required")
	}

	s.mu.Lock()
	ticketID := fmt.Sprintf("ticket-%03d", s.nextTicketID)
	s.nextTicketID++
	ticket := advisorTicket{
		TicketID:    ticketID,
		StudentID:   strings.TrimSpace(input.StudentID),
		StudyPlanID: strings.TrimSpace(input.StudyPlanID),
		Summary:     input.Summary,
		Labels:      append([]string{}, input.Labels...),
		Status:      "open",
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	s.advisorTickets[ticketID] = ticket
	s.mu.Unlock()

	return nil, createAdvisorTicketOutput{Ticket: ticket}, nil
}

func (s *store) patchAdvisorTicket(ctx context.Context, req *mcp.CallToolRequest, input patchAdvisorTicketInput) (*mcp.CallToolResult, patchAdvisorTicketOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, patchAdvisorTicketOutput{}, err
	}

	ticketID := strings.TrimSpace(input.TicketID)
	if ticketID == "" {
		return nil, patchAdvisorTicketOutput{}, fmt.Errorf("ticket_id is required")
	}

	s.mu.Lock()
	ticket, found := s.advisorTickets[ticketID]
	if found {
		if strings.TrimSpace(input.Summary) != "" {
			ticket.Summary = input.Summary
		}
		if len(input.Labels) > 0 {
			ticket.Labels = append([]string{}, input.Labels...)
		}
		if strings.TrimSpace(input.Status) != "" {
			ticket.Status = strings.TrimSpace(input.Status)
		}
		ticket.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		s.advisorTickets[ticketID] = ticket
	}
	s.mu.Unlock()

	if !found {
		return nil, patchAdvisorTicketOutput{Found: false}, nil
	}
	return nil, patchAdvisorTicketOutput{Found: true, Ticket: &ticket}, nil
}

func (s *store) submitSupportRequest(ctx context.Context, req *mcp.CallToolRequest, input submitSupportRequestInput) (*mcp.CallToolResult, submitSupportRequestOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, submitSupportRequestOutput{}, err
	}
	if strings.TrimSpace(input.StudentID) == "" {
		return nil, submitSupportRequestOutput{}, fmt.Errorf("student_id is required")
	}
	if strings.TrimSpace(input.Reason) == "" {
		return nil, submitSupportRequestOutput{}, fmt.Errorf("reason is required")
	}

	s.mu.Lock()
	requestID := fmt.Sprintf("support-%03d", s.nextSupportRequestID)
	s.nextSupportRequestID++
	request := supportRequest{
		RequestID: requestID,
		StudentID: strings.TrimSpace(input.StudentID),
		Reason:    input.Reason,
		Channels:  append([]string{}, input.Channels...),
		Status:    "submitted",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.supportRequests[requestID] = request
	s.mu.Unlock()

	return nil, submitSupportRequestOutput{Request: request}, nil
}

func courseMatchesQuery(item course, query string) bool {
	if strings.Contains(strings.ToLower(item.Code), query) || strings.Contains(strings.ToLower(item.Name), query) {
		return true
	}
	for _, tag := range item.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func courseHasTags(item course, required []string) bool {
	for _, requiredTag := range required {
		found := false
		for _, tag := range item.Tags {
			if strings.EqualFold(strings.TrimSpace(requiredTag), tag) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func noteHasLabels(item note, required []string) bool {
	for _, requiredLabel := range required {
		found := false
		for _, label := range item.Labels {
			if strings.EqualFold(strings.TrimSpace(requiredLabel), label) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
