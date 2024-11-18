package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			zap.L().Error("Cannot parse payload", zap.Error(err))
			return
		}
	}
}

func ResponseWithError(w http.ResponseWriter, code int, err error, msg string) {
	zap.L().Error("Exception occurred", zap.Error(err))
	if code > 499 {
		zap.L().Error("Responding with 5XX error", zap.Error(err))
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	ResponseWithJson(w, code, errorResponse{
		Error: msg,
	})
}

func ParsePayload(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}

func HashingPwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("cannot hash password", zap.Error(err))
	}

	return string(hash)
}

func VerifyPwdHash(pwd string, userPwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPwdHash), []byte(pwd))
	if err != nil {
		zap.L().Error("cannot verify password", zap.Error(err))
	}
	return err == nil
}

func GenerateJwt(userId, role, sessionId string) (string, error) {
	claims := jwt.MapClaims{
		"userId":    userId,
		"role":      role,
		"sessionId": sessionId,
		"exp":       time.Now().Add(time.Hour * 2).Unix(), //Adding 2 hours for testing. Change this to 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
