package routes

import (
	"log"
	"os"

	"server/data"
	"server/handlers"
	"server/usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	scenario := os.Getenv("SEED")
	fixture, err := data.LoadFixture(scenario)
	if err != nil {
		log.Fatalf("load fixture: %v", err)
	}
	store := data.NewStoreFromFixture(*fixture)
	repos := data.NewRepos(store)

	deps := usecase.Deps{
		Students:      repos.Students,
		Courses:       repos.Courses,
		Enrollments:   repos.Enrollments,
		Schedules:     repos.Schedules,
		ExamSchedules: repos.ExamSchedules,
		Scores:        repos.Scores,
		Assignments:   repos.Assignments,
		Materials:     repos.Materials,
		Deadlines:     repos.Deadlines,
		Rooms:         repos.Rooms,
		RoomBookings:  repos.RoomBookings,
		Submissions:   repos.Submissions,
		Requests:      repos.Requests,
		Contacts:      repos.Contacts,
	}

	uc := usecase.NewService(deps)
	h := handlers.NewHandler(uc, store.Reset)

	api := app.Group("/api")
	api.Get("/students", h.GetStudents)
	api.Get("/students/:id", h.GetStudentByID)

	app.Post("/login", h.Login)

	admin := app.Group("/admin")
	admin.Post("/reset", h.ResetStore)

	student := app.Group("/student")
	student.Get("/profile", h.GetProfile)
	student.Get("/schedule", h.GetSchedule)
	student.Get("/schedule/exam", h.GetExamSchedule)
	student.Get("/score", h.GetScore)
	student.Get("/lookup/tuitionfee", h.GetTuitionFee)
	student.Get("/insurance", h.GetInsurance)
	student.Get("/courses", h.GetCourses)
	student.Get("/training-points", h.GetTrainingPoints)
	student.Get("/survey-form", h.GetSurveyForm)
	student.Post("/transcript-regis", h.TranscriptRegis)
	student.Post("/tuition-extend", h.TuitionExtend)
	student.Post("/monthly-parking", h.MonthlyParking)
	student.Post("/graduate", h.Graduate)
	student.Post("/graduation-thesis", h.GraduationThesis)
	student.Get("/deadlines", h.GetDeadlines)
	student.Get("/courses/:courseId/materials", h.GetMaterials)
	student.Get("/courses/:courseId/assignments", h.GetAssignments)
	student.Post("/assignments/:assignmentId/submissions", h.SubmitAssignment)
	student.Post("/confirm-letter", h.ConfirmLetter)
	student.Post("/bank-loans", h.BankLoans)
	student.Post("/training-point-confirm", h.TrainingPointConfirm)
	student.Post("/language-certificate", h.LanguageCertificate)

	app.Get("/rooms/availability", h.GetRoomsAvailability)
	app.Post("/contact", h.Contact)
}
