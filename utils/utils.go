package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			log.Printf("Cannot parse payload : %v", err)
			return
		}
	}
}

func ResponseWithError(w http.ResponseWriter, code int, err error, msg string) {
	log.Printf("Error: %v", err)
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
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
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJwt(userId, role, sessionId string) (string, error) {
	claims := jwt.MapClaims{
		"userId":    userId,
		"role":      role,
		"sessionId": sessionId,
		"exp":       time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}