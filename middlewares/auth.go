package middlewares

import (
	"context"
	"errors"
	"github.com/gauraveg/rmsapp/logger"
	"net/http"
	"os"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/golang-jwt/jwt/v5"
)

type userCon string

const userContext userCon = "userContext"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			loggers := logger.GetLogContext(r)
			tokenString := r.Header.Get("token")
			if tokenString == "" {
				loggers.ErrorWithContext(r.Context(), "No token provided")
				utils.ResponseWithError(r.Context(), loggers, w, http.StatusUnauthorized, nil, "token header missing")
				return
			}
			loggers.InfoWithContext(r.Context(), "Parsing JWT token")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil || !token.Valid {
				loggers.ErrorWithContext(r.Context(), map[string]string{"message": "invalid token", "Token": tokenString})
				utils.ResponseWithError(r.Context(), loggers, w, http.StatusUnauthorized, err, "invalid token")
				return
			}

			claimValues, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				loggers.ErrorWithContext(r.Context(), map[string]string{"message": "invalid token claims", "Token": tokenString})
				utils.ResponseWithError(r.Context(), loggers, w, http.StatusUnauthorized, nil, "invalid token claims")
				return
			}

			sessionId := claimValues["sessionId"].(string)
			userData, err := dbHelper.FetchUserDataBySessionId(sessionId)
			if err != nil {
				loggers.ErrorWithContext(r.Context(), map[string]string{"message": "failed to fetch User data using the sessionId", "sessionId": sessionId})
				utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, err, "Failed to fetch User data using the sessionID")
				return
			}
			if userData.ArchivedAt != nil {
				loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Session is already expired", "sessionId": sessionId})
				utils.ResponseWithError(r.Context(), loggers, w, http.StatusUnauthorized, nil, "Session is already expired")
				return
			}

			user := &models.UserCtx{
				UserID:    claimValues["userId"].(string),
				SessionID: sessionId,
				Role:      models.Role(claimValues["role"].(string)),
				Email:     userData.Email,
			}

			ctx := context.WithValue(r.Context(), userContext, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}

func UserContext(r *http.Request) *models.UserCtx {
	user, ok := r.Context().Value(userContext).(*models.UserCtx)
	if !ok {
		return nil
	}
	return user
}
