package data

import (
	"sync"

	"server/models"
)

type Store struct {
	mu sync.RWMutex

	base Fixture

	students      []models.Student
	courses       []models.Course
	enrollments   []models.Enrollment
	schedules     []models.Schedule
	examSchedules []models.ExamSchedule
	scores        []models.Score
	assignments   []models.Assignment
	materials     []models.Material
	deadlines     []models.Deadline
	rooms         []models.Room
	roomBookings  []models.RoomBooking
	submissions   []models.Submission
	requests      []models.Request
	contacts      []models.Contact
}

func NewStoreFromFixture(f Fixture) *Store {
	s := &Store{base: f}
	s.Reset()
	return s
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.students = append([]models.Student(nil), s.base.Students...)
	s.courses = append([]models.Course(nil), s.base.Courses...)
	s.enrollments = append([]models.Enrollment(nil), s.base.Enrollments...)
	s.schedules = append([]models.Schedule(nil), s.base.Schedules...)
	s.examSchedules = append([]models.ExamSchedule(nil), s.base.ExamSchedules...)
	s.scores = append([]models.Score(nil), s.base.Scores...)
	s.assignments = append([]models.Assignment(nil), s.base.Assignments...)
	s.materials = append([]models.Material(nil), s.base.Materials...)
	s.deadlines = append([]models.Deadline(nil), s.base.Deadlines...)
	s.rooms = append([]models.Room(nil), s.base.Rooms...)
	s.roomBookings = append([]models.RoomBooking(nil), s.base.RoomBookings...)
	s.submissions = append([]models.Submission(nil), s.base.Submissions...)
	s.requests = append([]models.Request(nil), s.base.Requests...)
	s.contacts = append([]models.Contact(nil), s.base.Contacts...)
}

func (s *Store) Students() []models.Student {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Student(nil), s.students...)
}

func (s *Store) FindStudentByID(id string) (*models.Student, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.students {
		if s.students[i].ID == id {
			student := s.students[i]
			return &student, true
		}
	}
	return nil, false
}

func (s *Store) Courses() []models.Course {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Course(nil), s.courses...)
}

func (s *Store) FindCourseByID(id string) (*models.Course, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.courses {
		if s.courses[i].ID == id {
			course := s.courses[i]
			return &course, true
		}
	}
	return nil, false
}

func (s *Store) Enrollments() []models.Enrollment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Enrollment(nil), s.enrollments...)
}

func (s *Store) Schedules() []models.Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Schedule(nil), s.schedules...)
}

func (s *Store) ExamSchedules() []models.ExamSchedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.ExamSchedule(nil), s.examSchedules...)
}

func (s *Store) Scores() []models.Score {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Score(nil), s.scores...)
}

func (s *Store) Assignments() []models.Assignment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Assignment(nil), s.assignments...)
}

func (s *Store) FindAssignmentByID(id string) (*models.Assignment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.assignments {
		if s.assignments[i].ID == id {
			assignment := s.assignments[i]
			return &assignment, true
		}
	}
	return nil, false
}

func (s *Store) Materials() []models.Material {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Material(nil), s.materials...)
}

func (s *Store) Deadlines() []models.Deadline {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Deadline(nil), s.deadlines...)
}

func (s *Store) Rooms() []models.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Room(nil), s.rooms...)
}

func (s *Store) RoomBookings() []models.RoomBooking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.RoomBooking(nil), s.roomBookings...)
}

func (s *Store) AddSubmission(submission models.Submission) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.submissions = append(s.submissions, submission)
}

func (s *Store) SubmissionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.submissions)
}

func (s *Store) AddRequest(request models.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, request)
}

func (s *Store) AddContact(contact models.Contact) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contacts = append(s.contacts, contact)
}

func (s *Store) ContactCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.contacts)
}
