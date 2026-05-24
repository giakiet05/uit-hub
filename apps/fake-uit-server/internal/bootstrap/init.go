package bootstrap

import (
	"log"

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

	svc := do.MustInvoke[*service.Service](injector)
	ctrl := do.MustInvoke[*controller.UitController](injector)

	route.RegisterUitRoutes(api, ctrl, svc)
}
