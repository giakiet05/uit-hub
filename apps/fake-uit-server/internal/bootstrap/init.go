package bootstrap

import (
<<<<<<< HEAD
	"log"

=======
>>>>>>> origin/dev
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/config"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/controller"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/route"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func Init() (*gin.Engine, error) {
	injector := do.New()

	do.Provide(injector, func(i *do.Injector) (*config.Config, error) {
		return config.NewConfig(), nil
	})

<<<<<<< HEAD
	do.Provide(injector, func(i *do.Injector) (*service.Service, error) {
		// Load fixture data
		fixture, err := service.LoadFixture("default")
		if err != nil {
			log.Fatalf("failed to load fixture: %v", err)
		}
		
		// Initialize the data store
		store := service.NewStoreFromFixture(*fixture)
		
		repos := service.NewRepos(store)
		
		// Create dependencies for the service
		deps := service.Deps{
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
		
		return service.NewService(deps), nil
	})

	do.Provide(injector, func(i *do.Injector) (*controller.UitController, error) {
		svc := do.MustInvoke[*service.Service](i)
		return controller.NewUitController(svc), nil
=======
	do.Provide(injector, func(i *do.Injector) (service.StudentService, error) {
		return service.NewStudentService(), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.AuthService, error) {
		return service.NewAuthService(), nil
	})

	do.Provide(injector, func(i *do.Injector) (*controller.StudentController, error) {
		studentService := do.MustInvoke[service.StudentService](i)
		return controller.NewStudentController(studentService), nil
	})

	do.Provide(injector, func(i *do.Injector) (*controller.AuthController, error) {
		authService := do.MustInvoke[service.AuthService](i)
		return controller.NewAuthController(authService), nil
>>>>>>> origin/dev
	})

	router := gin.Default()
	registerRoutes(router, injector)

	return router, nil
}

func registerRoutes(router *gin.Engine, injector *do.Injector) {
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	cfg := do.MustInvoke[*config.Config](injector)

	api := router.Group("/api/v1")
	api.Use(middleware.Chaos(cfg.Chaos))

<<<<<<< HEAD
	svc := do.MustInvoke[*service.Service](injector)
	ctrl := do.MustInvoke[*controller.UitController](injector)

	route.RegisterUitRoutes(api, ctrl, svc)
=======
	authService := do.MustInvoke[service.AuthService](injector)

	route.RegisterAuthRoutes(api, do.MustInvoke[*controller.AuthController](injector))
	route.RegisterStudentRoutes(api, do.MustInvoke[*controller.StudentController](injector), authService)
>>>>>>> origin/dev
}
