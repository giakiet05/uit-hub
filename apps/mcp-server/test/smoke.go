// Quick smoke test: calls all tools via the HTTP client directly (not via MCP stdio).
// Run: go run test/smoke_test.go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/client"
)

func main() {
	c := client.New()
	fmt.Println("Smoke test → fake-uit-server at", c.BaseURL)

	passed, failed := 0, 0

	run := func(name string, fn func() (map[string]any, error)) map[string]any {
		resp, err := fn()
		if err != nil {
			fmt.Printf("💥 %s → system error: %v\n", name, err)
			failed++
			return nil
		}
		if client.IsBusinessError(resp) {
			fmt.Printf("❌ %s → business error: %v\n", name, resp["message"])
			failed++
			return resp
		}
		fmt.Printf("✅ %s\n", name)
		passed++
		return resp
	}

	// ── 1. Auth ─────────────────────────────────────────────────
	loginResp := run("auth_login", func() (map[string]any, error) {
		return c.Post("/login", "", map[string]any{
			"student_id": "22520001",
			"password":   "pass123",
		})
	})

	token := ""
	if loginResp != nil {
		if data, ok := loginResp["data"].(map[string]any); ok {
			if t, ok := data["token"].(string); ok {
				token = t
			}
		}
	}
	if token == "" {
		fmt.Println("Cannot proceed without token")
		os.Exit(1)
	}
	fmt.Println("  token:", token)

	// ── 2. Student Queries ──────────────────────────────────────
	run("get_student_profile", func() (map[string]any, error) {
		return c.Get("/student/profile", token, nil)
	})
	run("get_student_schedule", func() (map[string]any, error) {
		return c.Get("/student/schedule", token, map[string]string{"year": "2025", "semester": "1"})
	})
	run("get_exam_schedule", func() (map[string]any, error) {
		return c.Get("/student/schedule/exam", token, map[string]string{"year": "2025", "semester": "1"})
	})
	run("get_student_scores", func() (map[string]any, error) {
		return c.Get("/student/score", token, nil)
	})
	run("get_tuition_fee", func() (map[string]any, error) {
		return c.Get("/student/lookup/tuitionfee", token, nil)
	})
	run("get_insurance", func() (map[string]any, error) {
		return c.Get("/student/insurance", token, nil)
	})
	run("get_enrolled_courses", func() (map[string]any, error) {
		return c.Get("/student/courses", token, nil)
	})
	run("get_training_points", func() (map[string]any, error) {
		return c.Get("/student/training-points", token, nil)
	})
	run("get_survey_forms", func() (map[string]any, error) {
		return c.Get("/student/survey-form", token, nil)
	})

	// ── 3. Student Mutations ────────────────────────────────────
	run("create_contact", func() (map[string]any, error) {
		return c.Post("/contact", "", map[string]any{
			"name": "Test", "email": "test@uit.edu.vn",
			"subject": "Test", "message": "Smoke test",
		})
	})
	run("create_transcript_request", func() (map[string]any, error) {
		return c.Post("/student/transcript-regis", token, map[string]any{
			"copies": 1, "language": "VI", "delivery_method": "PICKUP",
		})
	})
	run("create_tuition_extension", func() (map[string]any, error) {
		return c.Post("/student/tuition-extend", token, map[string]any{
			"year": 2024, "semester": 1, "requested_due_date": "2024-12-31", "reason": "Test",
		})
	})
	run("create_monthly_parking", func() (map[string]any, error) {
		return c.Post("/student/monthly-parking", token, map[string]any{
			"vehicle_type": "MOTORBIKE", "plate_number": "59X1-12345",
			"months": 3, "start_month": "2024-05",
		})
	})
	run("create_graduate_request", func() (map[string]any, error) {
		return c.Post("/student/graduate", token, map[string]any{
			"year": 2024, "semester": 1,
			"email": "22520001@gm.uit.edu.vn", "phone": "0901234001",
		})
	})
	run("create_graduation_thesis_request", func() (map[string]any, error) {
		return c.Post("/student/graduation-thesis", token, map[string]any{
			"thesis_title": "AI", "advisor_name": "Dr. A",
			"team_members": []map[string]any{{"student_id": "22520001"}},
		})
	})

	// ── 4. Rooms ────────────────────────────────────────────────
	run("get_rooms_availability", func() (map[string]any, error) {
		return c.Get("/rooms/availability", token, map[string]string{
			"date": "2025-04-15", "start": "08:00", "end": "10:00",
		})
	})

	// ── 5. Course ───────────────────────────────────────────────
	run("get_deadlines", func() (map[string]any, error) {
		return c.Get("/student/deadlines", token, nil)
	})
	run("get_course_materials", func() (map[string]any, error) {
		return c.Get("/student/courses/IT010/materials", token, nil)
	})
	run("get_course_assignments", func() (map[string]any, error) {
		return c.Get("/student/courses/IT010/assignments", token, nil)
	})
	run("submit_assignment", func() (map[string]any, error) {
		return c.Post("/student/assignments/ASG001/submissions", token, map[string]any{
			"submission_type": "LINK", "url": "https://github.com/test",
		})
	})

	// ── 6. CTSV ─────────────────────────────────────────────────
	run("create_confirm_letter", func() (map[string]any, error) {
		return c.Post("/student/confirm-letter", token, map[string]any{
			"language": "VI", "reason": "DORM_EXTEND", "request_type": "NEW",
		})
	})
	run("create_bank_loan_request", func() (map[string]any, error) {
		return c.Post("/student/bank-loans", token, map[string]any{
			"benefit": "NO_DISCOUNT", "orphan_status": "NOT_ORPHAN", "template": "LEGACY",
		})
	})
	run("create_training_point_confirm", func() (map[string]any, error) {
		return c.Post("/student/training-point-confirm", token, map[string]any{
			"language": "VI",
		})
	})
	run("upload_language_certificate", func() (map[string]any, error) {
		return c.Post("/student/language-certificate", token, map[string]any{
			"document_type": "CERTIFICATE", "birth_date": "2004-01-01",
			"id_number": "079204000001", "listening_score": 450,
			"reading_score": 400, "total_score": 850,
			"exam_date": "2024-01-01", "image_file": "aGVsbG8=",
		})
	})

	// ── 7. Public ───────────────────────────────────────────────
	run("list_students", func() (map[string]any, error) {
		return c.Get("/students", "", nil)
	})

	fmt.Println()
	fmt.Println("═══════════════ REPORT ═══════════════")
	fmt.Printf("Total: %d | Passed: %d | Failed: %d\n", passed+failed, passed, failed)

	_ = json.Marshal // suppress unused import if needed
}
