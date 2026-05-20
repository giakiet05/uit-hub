package usecase

type Deps struct {
	Students      StudentRepo
	Courses       CourseRepo
	Enrollments   EnrollmentRepo
	Schedules     ScheduleRepo
	ExamSchedules ExamScheduleRepo
	Scores        ScoreRepo
	Assignments   AssignmentRepo
	Materials     MaterialRepo
	Deadlines     DeadlineRepo
	Rooms         RoomRepo
	RoomBookings  RoomBookingRepo
	Submissions   SubmissionRepo
	Requests      RequestRepo
	Contacts      ContactRepo
}
