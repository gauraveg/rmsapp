package middlewares

import (
	"context"
	"errors"
	"go.uber.org/zap"
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
			logger := LoggerContext(r)
			tokenString := r.Header.Get("token")
			if tokenString == "" {
				logger.Error("No token provided")
				utils.ResponseWithError(w, http.StatusUnauthorized, nil, "token header missing")
				return
			}
			logger.Info("Parsing JWT token")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil || !token.Valid {
				logger.Error("invalid token", zap.String("Toekn", tokenString))
				utils.ResponseWithError(w, http.StatusUnauthorized, err, "invalid token")
				return
			}

			claimValues, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				logger.Error("invalid token claims")
				utils.ResponseWithError(w, http.StatusUnauthorized, nil, "invalid token claims")
				return
			}

			sessionId := claimValues["sessionId"].(string)
			userData, err := dbHelper.FetchUserDataBySessionId(sessionId)
			if err != nil {
				logger.Error("Failed to fetch User data using the sessionID", zap.String("sessionId", sessionId))
				utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch User data using the sessionID")
				return
			}
			if userData.ArchivedAt != nil {
				logger.Error("Session is already expired", zap.String("sessionId", sessionId))
				utils.ResponseWithError(w, http.StatusUnauthorized, nil, "Session is already expired")
				return
			}

			user := &models.UserCtx{
				UserID:    claimValues["userId"].(string),
				SessionID: sessionId,
				Role:      claimValues["role"].(string),
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
