package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverVersion = "0.1.0"

type store struct {
	mu            sync.RWMutex
	students      map[string]studentProfile
	courses       []course
	notes         map[string]note
	leaveRequests map[string]leaveRequest
	nextNoteID    int
	nextLeaveID   int
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
	Code    string   `json:"code" jsonschema:"course code"`
	Name    string   `json:"name" jsonschema:"course name"`
	Credits int      `json:"credits" jsonschema:"course credit count"`
	Tags    []string `json:"tags" jsonschema:"course tags"`
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

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "uit-hub-mock-mcp-server",
		Version: serverVersion,
	}, nil)

	data := newStore()
	registerTools(server, data)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("mock mcp server failed: %v", err)
	}
}

func registerTools(server *mcp.Server, data *store) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_student_profile",
		Description: "Read a mock student profile by student ID.",
	}, data.getStudentProfile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_courses",
		Description: "Read mock course catalog entries using query text and tags.",
	}, data.searchCourses)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "upsert_student_note",
		Description: "Create or update an advisor note for a student. This is a write operation.",
	}, data.upsertStudentNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_student_notes",
		Description: "Read advisor notes for a student.",
	}, data.listStudentNotes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_student_note",
		Description: "Delete an advisor note by note ID. This is a write operation.",
	}, data.deleteStudentNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "submit_leave_request",
		Description: "Create a mock leave request for a student. This is a write operation with nested object and array input.",
	}, data.submitLeaveRequest)
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
			{Code: "CS101", Name: "Introduction to Programming", Credits: 4, Tags: []string{"programming", "foundation"}},
			{Code: "CS201", Name: "Data Structures and Algorithms", Credits: 4, Tags: []string{"programming", "algorithm"}},
			{Code: "IS301", Name: "Database Systems", Credits: 3, Tags: []string{"database", "information_systems"}},
			{Code: "SE401", Name: "Software Architecture", Credits: 3, Tags: []string{"software_engineering", "architecture"}},
		},
		notes:         map[string]note{},
		leaveRequests: map[string]leaveRequest{},
		nextNoteID:    1,
		nextLeaveID:   1,
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
