package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gauraveg/rmsapp/logger"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gauraveg/rmsapp/models"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func ResponseWithJson(ctx context.Context, loggers *logger.ZapLogger, w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			loggers.ErrorWithContext(ctx, map[string]string{"message": "cannot parse payload"})
			return
		}
	}

	if code == 200 || code == 201 {
		loggers.InfoWithContext(ctx, map[string]interface{}{"message": "response parsed successfully", "response": payload})
	}
}

func ResponseWithError(ctx context.Context, loggers *logger.ZapLogger, w http.ResponseWriter, code int, err error, msg string) {
	loggers.ErrorWithContext(ctx, map[string]string{"message": msg, "error": err.Error()})
	if code > 499 {
		zap.L().Error("Responding with 5XX error", zap.Error(err))
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	ResponseWithJson(ctx, loggers, w, code, errorResponse{
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

func hsin(value float64) float64 {
	return math.Pow(math.Sin(value/2), 2)
}

func CalculateDistBetweenPoints(restPoint []models.Coordinates, userPoint []models.Coordinates) map[string]string {
	AddrDistance := make(map[string]string)
	for _, Point := range userPoint {
		if Point.Address != "" {
			userLatInRad := (Point.Latitude * math.Pi) / 180
			userLongInRad := (Point.Longitude * math.Pi) / 180

			restLatInRad := (restPoint[0].Latitude * math.Pi) / 180
			restLongInRad := (restPoint[0].Longitude * math.Pi) / 180

			earthRadius := float64(6378100)
			h := hsin(restLatInRad-userLatInRad) + math.Cos(userLatInRad)*math.Cos(restLatInRad)*hsin(restLongInRad-userLongInRad)
			result := 2 * earthRadius * math.Asin(math.Sqrt(h)) / 1000
			AddrDistance[Point.Address] = fmt.Sprintf("%.3f km", result)
		}
	}
	return AddrDistance
}

//func GetPayload(ctx context.Context) (map[string]interface{}, error) {
//	body, _ := ctx.Value("payload").(map[string]interface{})
//	jsonBody, _ := json.Marshal(body)
//	err := json.Unmarshal(jsonBody, &payload)
//	return payload, err
//}
