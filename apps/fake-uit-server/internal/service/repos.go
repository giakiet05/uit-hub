package service

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

type StudentRepo interface {
	List() ([]model.Student, error)
	FindByID(id string) (*model.Student, bool, error)
}

type CourseRepo interface {
	FindByID(id string) (*model.Course, bool, error)
}

type EnrollmentRepo interface {
	List() ([]model.Enrollment, error)
}

type ScheduleRepo interface {
	List() ([]model.Schedule, error)
}

type ExamScheduleRepo interface {
	List() ([]model.ExamSchedule, error)
}

type ScoreRepo interface {
	List() ([]model.Score, error)
}

type AssignmentRepo interface {
	List() ([]model.Assignment, error)
	FindByID(id string) (*model.Assignment, bool, error)
}

type MaterialRepo interface {
	List() ([]model.Material, error)
}

type DeadlineRepo interface {
	List() ([]model.Deadline, error)
}

type RoomRepo interface {
	List() ([]model.Room, error)
}

type RoomBookingRepo interface {
	List() ([]model.RoomBooking, error)
}

type SubmissionRepo interface {
	Create(submission model.Submission) error
	Count() (int, error)
}

type RequestRepo interface {
	Create(request model.Request) error
}

type ContactRepo interface {
	Create(contact model.Contact) error
	Count() (int, error)
}
