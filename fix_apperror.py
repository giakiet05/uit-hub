import re

with open('apps/fake-uit-server/internal/service/usecase.go', 'r') as f:
    content = f.read()

# Fix AppError pointers
content = content.replace('*apperror.AppError', '*apperror.AppError') # Keep pointer for now

# Add New helper functions to usecase.go to avoid changing all 50 occurrences
helpers = """
import (
	"net/http"
"""
content = content.replace('import (', helpers, 1)

helpers_func = """

func newAppError(status int, code, msg string) *apperror.AppError {
	return &apperror.AppError{Status: status, Code: code, Message: msg}
}

func newBadRequest(msg string) *apperror.AppError {
	if msg == "" { msg = "Bad request!" }
	return newAppError(http.StatusBadRequest, "BAD_REQUEST", msg)
}
func newUnauthorized(msg string) *apperror.AppError {
	if msg == "" { msg = "Unauthorized!" }
	return newAppError(http.StatusUnauthorized, "UNAUTHORIZED", msg)
}
func newNotFound(msg string) *apperror.AppError {
	if msg == "" { msg = "Not found!" }
	return newAppError(http.StatusNotFound, "NOT_FOUND", msg)
}
func newInternal(msg string) *apperror.AppError {
	if msg == "" { msg = "Internal error happened!" }
	return newAppError(http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

"""
content = content.replace('const (', helpers_func + '\nconst (', 1)

content = content.replace('apperror.NewBadRequest', 'newBadRequest')
content = content.replace('apperror.NewUnauthorized', 'newUnauthorized')
content = content.replace('apperror.NewInternal', 'newInternal')
content = content.replace('apperror.NewNotFound', 'newNotFound')

with open('apps/fake-uit-server/internal/service/usecase.go', 'w') as f:
    f.write(content)

