package service

import (
	"sync"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"
)

type Store struct {
	mu sync.RWMutex

	base Fixture

	students      []model.Student
	courses       []model.Course
	enrollments   []model.Enrollment
	schedules     []model.Schedule
	examSchedules []model.ExamSchedule
	scores        []model.Score
	assignments   []model.Assignment
	materials     []model.Material
	deadlines     []model.Deadline
	rooms         []model.Room
	roomBookings  []model.RoomBooking
	submissions   []model.Submission
	requests      []model.Request
	contacts      []model.Contact
}

func NewStoreFromFixture(f Fixture) *Store {
	s := &Store{base: f}
	s.Reset()
	return s
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.students = append([]model.Student(nil), s.base.Students...)
	s.courses = append([]model.Course(nil), s.base.Courses...)
	s.enrollments = append([]model.Enrollment(nil), s.base.Enrollments...)
	s.schedules = append([]model.Schedule(nil), s.base.Schedules...)
	s.examSchedules = append([]model.ExamSchedule(nil), s.base.ExamSchedules...)
	s.scores = append([]model.Score(nil), s.base.Scores...)
	s.assignments = append([]model.Assignment(nil), s.base.Assignments...)
	s.materials = append([]model.Material(nil), s.base.Materials...)
	s.deadlines = append([]model.Deadline(nil), s.base.Deadlines...)
	s.rooms = append([]model.Room(nil), s.base.Rooms...)
	s.roomBookings = append([]model.RoomBooking(nil), s.base.RoomBookings...)
	s.submissions = append([]model.Submission(nil), s.base.Submissions...)
	s.requests = append([]model.Request(nil), s.base.Requests...)
	s.contacts = append([]model.Contact(nil), s.base.Contacts...)
}

func (s *Store) Students() []model.Student {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Student(nil), s.students...)
}

func (s *Store) FindStudentByID(id string) (*model.Student, bool) {
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

func (s *Store) Courses() []model.Course {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Course(nil), s.courses...)
}

func (s *Store) FindCourseByID(id string) (*model.Course, bool) {
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

func (s *Store) Enrollments() []model.Enrollment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Enrollment(nil), s.enrollments...)
}

func (s *Store) Schedules() []model.Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Schedule(nil), s.schedules...)
}

func (s *Store) ExamSchedules() []model.ExamSchedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.ExamSchedule(nil), s.examSchedules...)
}

func (s *Store) Scores() []model.Score {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Score(nil), s.scores...)
}

func (s *Store) Assignments() []model.Assignment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Assignment(nil), s.assignments...)
}

func (s *Store) FindAssignmentByID(id string) (*model.Assignment, bool) {
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

func (s *Store) Materials() []model.Material {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Material(nil), s.materials...)
}

func (s *Store) Deadlines() []model.Deadline {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Deadline(nil), s.deadlines...)
}

func (s *Store) Rooms() []model.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Room(nil), s.rooms...)
}

func (s *Store) RoomBookings() []model.RoomBooking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.RoomBooking(nil), s.roomBookings...)
}

func (s *Store) AddSubmission(submission model.Submission) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.submissions = append(s.submissions, submission)
}

func (s *Store) SubmissionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.submissions)
}

func (s *Store) AddRequest(request model.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, request)
}

func (s *Store) AddContact(contact model.Contact) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contacts = append(s.contacts, contact)
}

func (s *Store) ContactCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.contacts)
}
