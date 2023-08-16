package response

import (
	"Interior_Visualization_Shop/app/internal/apperror"
	"net/http"
)

func Error(w http.ResponseWriter, code int, message, developerMessage string) {
	appError := apperror.NewAppError(code, message, developerMessage)
	JSON(w, code, appError)
}
func BadRequest(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusBadRequest, message, developerMessage)
}

func ErrorAuth(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusUnauthorized, message, developerMessage)
}

func NotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, apperror.ErrNotFound)
}

func InternalError(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusInternalServerError, message, developerMessage)
}
