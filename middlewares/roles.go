package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
)

func ShouldHaveRole(role models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := LoggerContext(r)
			userRole := UserContext(r).Role
			if userRole != role {
				w.WriteHeader(http.StatusForbidden)
				msg := fmt.Sprintf("Cannot access this endpoint as %v", userRole)
				logger.Error("endpoint forbidden", zap.String("Error", msg))
				utils.ResponseWithError(w, http.StatusBadRequest, errors.New("endpoint forbidden"), msg)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
