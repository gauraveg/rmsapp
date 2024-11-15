package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/golang-jwt/jwt/v5"
)

type usercon string

const usercontext usercon = "userContext"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("token")
			if tokenString == "" {
				utils.ResponseWithError(w, http.StatusUnauthorized, nil, "token header missing")
				return
			}
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil || !token.Valid {
				utils.ResponseWithError(w, http.StatusUnauthorized, err, "invalid token")
				return
			}

			claimValues, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				utils.ResponseWithError(w, http.StatusUnauthorized, nil, "invalid token claims")
				return
			}

			sessionId := claimValues["sessionId"].(string)
			userData, err := dbHelper.FetchUserDetails(sessionId)
			if err != nil {
				utils.ResponseWithError(w, http.StatusInternalServerError, err, "internal server error")
				return
			}
			if userData.ArchivedAt != nil {
				utils.ResponseWithError(w, http.StatusUnauthorized, nil, "invalid token")
				return
			}

			user := &models.UserCtx{
				UserID:    claimValues["userId"].(string),
				SessionID: sessionId,
				Role:      claimValues["role"].(string),
				Email:     userData.Email,
			}

			ctx := context.WithValue(r.Context(), usercontext, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}

func UserContext(r *http.Request) *models.UserCtx {
	user, ok := r.Context().Value(usercontext).(*models.UserCtx)
	if !ok {
		return nil
	}
	return user
}
