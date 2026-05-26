package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
<<<<<<< HEAD
	"github.com/gin-gonic/gin"
)

func (c *UitController) ListStudents(ctx *gin.Context) {

	res, appErr := c.service.ListStudents()
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetStudentByID(ctx *gin.Context) {

=======
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

type StudentController struct {
	service service.StudentService
}

func NewStudentController(service service.StudentService) *StudentController {
	return &StudentController{service: service}
}

func (c *StudentController) GetStudents(ctx *gin.Context) {
	students := c.service.GetStudents()
	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponses(students))
}

func (c *StudentController) GetProfile(ctx *gin.Context) {
>>>>>>> origin/dev
	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

<<<<<<< HEAD
	res, appErr := c.service.GetStudentByID(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetProfile(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetProfile(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetSchedule(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var query dto.QueryYearSemester
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.GetSchedule(studentID.(string), query)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetExamSchedule(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var query dto.QueryYearSemester
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.GetExamSchedule(studentID.(string), query)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetScore(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetScore(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetTuitionFee(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetTuitionFee(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetInsurance(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetInsurance(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetCourses(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetCourses(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetTrainingPoints(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetTrainingPoints(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetSurveyForm(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetSurveyForm(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) TranscriptRegis(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.TranscriptRegisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.TranscriptRegis(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) TuitionExtend(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.TuitionExtendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.TuitionExtend(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) MonthlyParking(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.MonthlyParkingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.MonthlyParking(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) Graduate(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.GraduateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.Graduate(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GraduationThesis(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.GraduationThesisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.GraduationThesis(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) Contact(ctx *gin.Context) {

	var req dto.ContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.Contact(req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
=======
	student, found := c.service.GetStudentByID(studentID.(string))
	if !found {
		dto.SendError(ctx, apperror.ErrNotFound)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponse(*student))
}

func (c *StudentController) GetStudentByID(ctx *gin.Context) {
	student, found := c.service.GetStudentByID(ctx.Param("id"))
	if !found {
		err := apperror.ErrNotFound
		err.Message = "Student not found"
		dto.SendError(ctx, err)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponse(*student))
>>>>>>> origin/dev
}
