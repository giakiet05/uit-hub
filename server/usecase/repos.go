package usecase

import "server/models"

type StudentRepo interface {
	List() ([]models.Student, error)
	FindByID(id string) (*models.Student, bool, error)
}

type CourseRepo interface {
	FindByID(id string) (*models.Course, bool, error)
}

type EnrollmentRepo interface {
	List() ([]models.Enrollment, error)
}

type ScheduleRepo interface {
	List() ([]models.Schedule, error)
}

type ExamScheduleRepo interface {
	List() ([]models.ExamSchedule, error)
}

type ScoreRepo interface {
	List() ([]models.Score, error)
}

type AssignmentRepo interface {
	List() ([]models.Assignment, error)
	FindByID(id string) (*models.Assignment, bool, error)
}

type MaterialRepo interface {
	List() ([]models.Material, error)
}

type DeadlineRepo interface {
	List() ([]models.Deadline, error)
}

type RoomRepo interface {
	List() ([]models.Room, error)
}

type RoomBookingRepo interface {
	List() ([]models.RoomBooking, error)
}

type SubmissionRepo interface {
	Create(submission models.Submission) error
	Count() (int, error)
}

type RequestRepo interface {
	Create(request models.Request) error
}

type ContactRepo interface {
	Create(contact models.Contact) error
	Count() (int, error)
}
