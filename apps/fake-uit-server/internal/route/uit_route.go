package route

import (
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/controller"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterUitRoutes(rg *gin.RouterGroup, c *controller.UitController, s *service.Service) {
	// Dummy struct to satisfy auth middleware requirement of AuthService for GetStudentByID
	// The middleware uses authService.GetStudentByID(id). In the new architecture, it's just s.GetStudentByID(id)
	
	// auth_login
	rg.POST("/login", c.Login)
	
	// student protected
	student := rg.Group("/student")
	student.Use(middleware.RequireAuth(s)) 
	{
		student.GET("/profile", c.GetProfile)
		student.GET("/schedule", c.GetSchedule)
		student.GET("/schedule/exam", c.GetExamSchedule)
		student.GET("/score", c.GetScore)
		student.GET("/lookup/tuitionfee", c.GetTuitionFee)
		student.GET("/insurance", c.GetInsurance)
		student.GET("/courses", c.GetCourses)
		student.GET("/training-points", c.GetTrainingPoints)
		student.GET("/survey-form", c.GetSurveyForm)
		student.POST("/transcript-regis", c.TranscriptRegis)
		student.POST("/tuition-extend", c.TuitionExtend)
		student.POST("/monthly-parking", c.MonthlyParking)
		student.POST("/graduate", c.Graduate)
		student.POST("/graduation-thesis", c.GraduationThesis)
		student.GET("/deadlines", c.GetDeadlines)
		student.GET("/courses/:courseId/materials", c.GetMaterials)
		student.GET("/courses/:courseId/assignments", c.GetAssignments)
		student.POST("/assignments/:assignmentId/submissions", c.SubmitAssignment)
		student.POST("/confirm-letter", c.ConfirmLetter)
		student.POST("/bank-loans", c.BankLoans)
		student.POST("/training-point-confirm", c.TrainingPointConfirm)
		student.POST("/language-certificate", c.LanguageCertificate)
	}
	
	// contact
	rg.POST("/contact", c.Contact)

	// rooms
	rooms := rg.Group("/rooms")
	rooms.Use(middleware.RequireAuth(s)) 
	{
		rooms.GET("/availability", c.GetRoomsAvailability)
	}
	
	// public students list (from old server)
	students := rg.Group("/students")
	students.GET("", c.ListStudents)
	students.GET("/:id", c.GetStudentByID)
}
