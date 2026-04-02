package data

import "server/models"

var Students = []models.Student{
	{ID: "22520001", Password: "pass123", Name: "Duc", Email: "duc@example.com", Phone: "0901000001", Major: "Software Engineering"},
	{ID: "22520002", Password: "pass123", Name: "An", Email: "an@example.com", Phone: "0901000002", Major: "Information Systems"},
}

var Courses = []models.Course{
	{ID: "SE101", Name: "Intro to Software Engineering", Credits: 3},
	{ID: "IS201", Name: "Database Systems", Credits: 3},
	{ID: "CS301", Name: "Computer Networks", Credits: 4},
}

var Enrollments = []models.Enrollment{
	{ID: 1, StudentID: "22520001", CourseID: "SE101", Year: 2025, Semester: 1},
	{ID: 2, StudentID: "22520001", CourseID: "IS201", Year: 2025, Semester: 1},
	{ID: 3, StudentID: "22520002", CourseID: "CS301", Year: 2025, Semester: 1},
}

var Schedules = []models.Schedule{
	{ID: 1, StudentID: "22520001", CourseID: "SE101", DayOfWeek: 1, StartTime: "08:00", EndTime: "09:50", Room: "A101", Year: 2025, Semester: 1},
	{ID: 2, StudentID: "22520001", CourseID: "IS201", DayOfWeek: 3, StartTime: "13:00", EndTime: "14:50", Room: "B202", Year: 2025, Semester: 1},
	{ID: 3, StudentID: "22520002", CourseID: "CS301", DayOfWeek: 2, StartTime: "09:00", EndTime: "10:50", Room: "C303", Year: 2025, Semester: 1},
}

var ExamSchedules = []models.ExamSchedule{
	{ID: 1, StudentID: "22520001", CourseID: "SE101", Date: "2025-12-20", Room: "A101", Year: 2025, Semester: 1},
	{ID: 2, StudentID: "22520001", CourseID: "IS201", Date: "2025-12-25", Room: "B202", Year: 2025, Semester: 1},
	{ID: 3, StudentID: "22520002", CourseID: "CS301", Date: "2025-12-22", Room: "C303", Year: 2025, Semester: 1},
}

var Scores = []models.Score{
	{ID: 1, StudentID: "22520001", CourseID: "SE101", Midterm: 7.5, Final: 8.2, Total: 8.0},
	{ID: 2, StudentID: "22520001", CourseID: "IS201", Midterm: 6.8, Final: 7.4, Total: 7.2},
	{ID: 3, StudentID: "22520002", CourseID: "CS301", Midterm: 8.0, Final: 8.5, Total: 8.3},
}

var Assignments = []models.Assignment{
	{ID: "A1", CourseID: "SE101", Title: "Sprint plan", Deadline: "2025-10-10"},
	{ID: "A2", CourseID: "IS201", Title: "Schema design", Deadline: "2025-10-15"},
	{ID: "A3", CourseID: "CS301", Title: "Routing lab", Deadline: "2025-10-20"},
}

var Materials = []models.Material{
	{ID: 1, CourseID: "SE101", Title: "Week 1 slides", URL: "https://example.com/se101/week1"},
	{ID: 2, CourseID: "IS201", Title: "ERD notes", URL: "https://example.com/is201/erd"},
	{ID: 3, CourseID: "CS301", Title: "TCP basics", URL: "https://example.com/cs301/tcp"},
}

var Deadlines = []models.Deadline{
	{ID: 1, StudentID: "22520001", Title: "SE101 - Sprint plan", DueDate: "2025-10-10"},
	{ID: 2, StudentID: "22520001", Title: "IS201 - Schema design", DueDate: "2025-10-15"},
	{ID: 3, StudentID: "22520002", Title: "CS301 - Routing lab", DueDate: "2025-10-20"},
}

var Rooms = []models.Room{
	{ID: 1, Name: "A101", Capacity: 80},
	{ID: 2, Name: "B202", Capacity: 60},
	{ID: 3, Name: "C303", Capacity: 40},
}

var RoomBookings = []models.RoomBooking{
	{ID: 1, RoomID: 1, Date: "2025-09-10", StartTime: "08:00", EndTime: "10:00"},
	{ID: 2, RoomID: 2, Date: "2025-09-10", StartTime: "13:00", EndTime: "15:00"},
}

var Submissions = []models.Submission{}

var Requests = []models.Request{}

var Contacts = []models.Contact{}
