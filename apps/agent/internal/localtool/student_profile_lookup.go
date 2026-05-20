package localtool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// StudentProfileLookup returns fake student profiles by student ID.
type StudentProfileLookup struct {
	profiles map[string]studentProfile
}

// studentProfile is the fake backend shape returned by StudentProfileLookup.
type studentProfile struct {
	StudentID       string  `json:"student_id"`
	FullName        string  `json:"full_name"`
	Major           string  `json:"major"`
	GPA             float64 `json:"gpa"`
	CompletedCredit int     `json:"completed_credits"`
	AcademicWarning bool    `json:"academic_warning"`
}

// NewStudentProfileLookup creates the fake student profile lookup tool.
func NewStudentProfileLookup() *StudentProfileLookup {
	return &StudentProfileLookup{
		profiles: map[string]studentProfile{
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
	}
}

// Definition describes the student profile lookup schema.
func (t *StudentProfileLookup) Definition() tool.Definition {
	return tool.Definition{
		Name:        "student_profile_lookup",
		Description: "Look up a student's profile from a fake school backend by student ID.",
		InputSchema: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"student_id": map[string]any{
					"type":        "string",
					"description": "Student ID, for example 22520001.",
				},
			},
			"required": []string{"student_id"},
		},
	}
}

// Execute finds a student profile by ID and returns a JSON result.
func (t *StudentProfileLookup) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	studentID, err := stringArg(call.Arguments, "student_id")
	if err != nil {
		return tool.Result{}, err
	}

	profile, found := t.profiles[strings.TrimSpace(studentID)]
	output, err := json.Marshal(map[string]any{
		"found":   found,
		"profile": profile,
	})
	if err != nil {
		return tool.Result{}, fmt.Errorf("marshal result: %w", err)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: string(output),
	}, nil
}
