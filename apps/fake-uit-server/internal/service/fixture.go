package service

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

type Fixture struct {
	Students      []model.Student      `json:"students"`
	Courses       []model.Course       `json:"courses"`
	Enrollments   []model.Enrollment   `json:"enrollments"`
	Schedules     []model.Schedule     `json:"schedules"`
	ExamSchedules []model.ExamSchedule `json:"exam_schedules"`
	Scores        []model.Score        `json:"scores"`
	Assignments   []model.Assignment   `json:"assignments"`
	Materials     []model.Material     `json:"materials"`
	Deadlines     []model.Deadline     `json:"deadlines"`
	Rooms         []model.Room         `json:"rooms"`
	RoomBookings  []model.RoomBooking  `json:"room_bookings"`
	Submissions   []model.Submission   `json:"submissions"`
	Requests      []model.Request      `json:"requests"`
	Contacts      []model.Contact      `json:"contacts"`
}
