package route

import (
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/controller"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(rg *gin.RouterGroup, c *controller.StudentController, authService service.AuthService) {
	student := rg.Group("/student")
	student.Use(middleware.RequireAuth(authService))
	{
		student.GET("/profile", c.GetProfile)
	}

	students := rg.Group("/students")

	students.GET("", c.GetStudents)
	students.GET("/:id", c.GetStudentByID)
}
