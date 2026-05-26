package middleware

import (
	"strings"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
<<<<<<< HEAD
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"
=======
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
>>>>>>> origin/dev
	"github.com/gin-gonic/gin"
)

const StudentIDKey = "student_id"

<<<<<<< HEAD
type AuthProvider interface {
	GetStudentByIDForAuth(id string) (*model.Student, bool)
}

// Since model.Student causes circular dep if not careful, we just use a small interface
type StudentFetcher interface {
	GetStudentByIDForAuth(id string) (any, bool)
}

func RequireAuth(provider any) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			dto.AbortWithError(c, apperror.ErrUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			dto.AbortWithError(c, apperror.ErrUnauthorized)
			return
		}

		token := parts[1]
		if !strings.HasPrefix(token, "mock-") {
			dto.AbortWithError(c, apperror.ErrUnauthorized)
			return
		}

		studentID := strings.TrimPrefix(token, "mock-")
		
		// Optional: actually fetch student to verify
		// We can just trust the token for mock server, but let's set it
		c.Set(StudentIDKey, studentID)
		c.Next()
	}
}
=======
func RequireAuth(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			dto.AbortWithError(ctx, apperror.ErrUnauthorized)
			return
		}

		authUser, ok := authService.ValidateToken(token)
		if !ok {
			dto.AbortWithError(ctx, apperror.ErrUnauthorized)
			return
		}

		ctx.Set(StudentIDKey, authUser.StudentID)
		ctx.Next()
	}
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}
>>>>>>> origin/dev
