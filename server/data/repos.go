package data

import (
	"server/models"
	"server/usecase"
)

type Repos struct {
	Students      usecase.StudentRepo
	Courses       usecase.CourseRepo
	Enrollments   usecase.EnrollmentRepo
	Schedules     usecase.ScheduleRepo
	ExamSchedules usecase.ExamScheduleRepo
	Scores        usecase.ScoreRepo
	Assignments   usecase.AssignmentRepo
	Materials     usecase.MaterialRepo
	Deadlines     usecase.DeadlineRepo
	Rooms         usecase.RoomRepo
	RoomBookings  usecase.RoomBookingRepo
	Submissions   usecase.SubmissionRepo
	Requests      usecase.RequestRepo
	Contacts      usecase.ContactRepo
}

func NewRepos(store *Store) Repos {
	return Repos{
		Students:      StudentRepository{store: store},
		Courses:       CourseRepository{store: store},
		Enrollments:   EnrollmentRepository{store: store},
		Schedules:     ScheduleRepository{store: store},
		ExamSchedules: ExamScheduleRepository{store: store},
		Scores:        ScoreRepository{store: store},
		Assignments:   AssignmentRepository{store: store},
		Materials:     MaterialRepository{store: store},
		Deadlines:     DeadlineRepository{store: store},
		Rooms:         RoomRepository{store: store},
		RoomBookings:  RoomBookingRepository{store: store},
		Submissions:   SubmissionRepository{store: store},
		Requests:      RequestRepository{store: store},
		Contacts:      ContactRepository{store: store},
	}
}

type StudentRepository struct{ store *Store }

type CourseRepository struct{ store *Store }

type EnrollmentRepository struct{ store *Store }

type ScheduleRepository struct{ store *Store }

type ExamScheduleRepository struct{ store *Store }

type ScoreRepository struct{ store *Store }

type AssignmentRepository struct{ store *Store }

type MaterialRepository struct{ store *Store }

type DeadlineRepository struct{ store *Store }

type RoomRepository struct{ store *Store }

type RoomBookingRepository struct{ store *Store }

type SubmissionRepository struct{ store *Store }

type RequestRepository struct{ store *Store }

type ContactRepository struct{ store *Store }

func (r StudentRepository) List() ([]models.Student, error) {
	return r.store.Students(), nil
}

func (r StudentRepository) FindByID(id string) (*models.Student, bool, error) {
	student, ok := r.store.FindStudentByID(id)
	return student, ok, nil
}

func (r CourseRepository) FindByID(id string) (*models.Course, bool, error) {
	course, ok := r.store.FindCourseByID(id)
	return course, ok, nil
}

func (r EnrollmentRepository) List() ([]models.Enrollment, error) {
	return r.store.Enrollments(), nil
}

func (r ScheduleRepository) List() ([]models.Schedule, error) {
	return r.store.Schedules(), nil
}

func (r ExamScheduleRepository) List() ([]models.ExamSchedule, error) {
	return r.store.ExamSchedules(), nil
}

func (r ScoreRepository) List() ([]models.Score, error) {
	return r.store.Scores(), nil
}

func (r AssignmentRepository) List() ([]models.Assignment, error) {
	return r.store.Assignments(), nil
}

func (r AssignmentRepository) FindByID(id string) (*models.Assignment, bool, error) {
	assignment, ok := r.store.FindAssignmentByID(id)
	return assignment, ok, nil
}

func (r MaterialRepository) List() ([]models.Material, error) {
	return r.store.Materials(), nil
}

func (r DeadlineRepository) List() ([]models.Deadline, error) {
	return r.store.Deadlines(), nil
}

func (r RoomRepository) List() ([]models.Room, error) {
	return r.store.Rooms(), nil
}

func (r RoomBookingRepository) List() ([]models.RoomBooking, error) {
	return r.store.RoomBookings(), nil
}

func (r SubmissionRepository) Create(submission models.Submission) error {
	r.store.AddSubmission(submission)
	return nil
}

func (r SubmissionRepository) Count() (int, error) {
	return r.store.SubmissionCount(), nil
}

func (r RequestRepository) Create(request models.Request) error {
	r.store.AddRequest(request)
	return nil
}

func (r ContactRepository) Create(contact models.Contact) error {
	r.store.AddContact(contact)
	return nil
}

func (r ContactRepository) Count() (int, error) {
	return r.store.ContactCount(), nil
}
