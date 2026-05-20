package data

import "server/models"

type Fixture struct {
	Students      []models.Student      `json:"students"`
	Courses       []models.Course       `json:"courses"`
	Enrollments   []models.Enrollment   `json:"enrollments"`
	Schedules     []models.Schedule     `json:"schedules"`
	ExamSchedules []models.ExamSchedule `json:"exam_schedules"`
	Scores        []models.Score        `json:"scores"`
	Assignments   []models.Assignment   `json:"assignments"`
	Materials     []models.Material     `json:"materials"`
	Deadlines     []models.Deadline     `json:"deadlines"`
	Rooms         []models.Room         `json:"rooms"`
	RoomBookings  []models.RoomBooking  `json:"room_bookings"`
	Submissions   []models.Submission   `json:"submissions"`
	Requests      []models.Request      `json:"requests"`
	Contacts      []models.Contact      `json:"contacts"`
}
