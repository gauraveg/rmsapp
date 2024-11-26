package middlewares

import (
	"errors"
	"fmt"
	"github.com/gauraveg/rmsapp/logger"
	"net/http"

	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
)

func ShouldHaveRole(role models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggers := logger.GetLogContext(r)
			userRole := UserContext(r).Role
			if userRole != role {
				w.WriteHeader(http.StatusForbidden)
				msg := fmt.Sprintf("endpoint forbidden. Cannot access this endpoint as %v", userRole)
				loggers.ErrorWithContext(r.Context(), msg)
				utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("endpoint forbidden"), msg)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
