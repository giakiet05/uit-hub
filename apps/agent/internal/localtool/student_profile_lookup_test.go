package localtool

import (
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestStudentProfileLookupFindsProfile(t *testing.T) {
	lookup := NewStudentProfileLookup()

	result, err := lookup.Execute(context.Background(), tool.Call{
		ID:   "call-student",
		Name: "student_profile_lookup",
		Arguments: map[string]any{
			"student_id": "22520002",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, "Tran Thi Binh") {
		t.Fatalf("result content = %q, want student name", result.Content)
	}
	if !strings.Contains(result.Content, `"academic_warning":true`) {
		t.Fatalf("result content = %q, want academic warning", result.Content)
	}
}
