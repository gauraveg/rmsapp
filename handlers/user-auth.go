package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gauraveg/rmsapp/logger"
	"net/http"
	"strings"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
)

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)

	var signUpData models.UserSignUp
	body, ok := r.Context().Value("payload").(string)
	if !ok {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload not present"), "cannot parse payload data")
		return
	}
	err := json.Unmarshal([]byte(body), &signUpData)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	payload := models.SignUpWithRole{
		UserSignUp: signUpData,
		Role:       models.RoleUser,
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload validation failed"), strings.Join(errMsg, "|"))
		return
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	userId, userEr := dbHelper.CreateSignUpHelper(payload.Email, payload.Name, hashedPwd, string(payload.Role), payload.Addresses)
	if userEr != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, userEr, fmt.Sprintf("Failed to create new user with email %s", payload.Email))
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusCreated, userId)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	var payload models.LoginRequest
	//err := utils.ParsePayload(r.Body, &payload)
	body, ok := r.Context().Value("payload").(string)
	if !ok {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload not present"), "cannot parse payload data")
		return
	}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload validation failed"), strings.Join(errMsg, "|"))
		return
	}

	userID, pwdHash, role, userErr := dbHelper.GetUserInfoForLogin(payload)
	if userErr != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, userErr, fmt.Sprintf("Failed to find user with email %v", payload.Email))
		return
	}

	if userID == "" || utils.VerifyPwdHash(payload.Password, pwdHash) {
		sessionID, crtErr := dbHelper.CreateUserSession(userID)
		if crtErr != nil {
			utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, crtErr, fmt.Sprintf("Failed to create user session with userId %v", userID))
			return
		}

		jwtToken, jwtErr := utils.GenerateJwt(userID, role, sessionID)
		if jwtErr != nil {
			utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, jwtErr, fmt.Sprintf("Failed to generate JWT Tokens for the sessionID %v", sessionID))
			return
		}

		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, models.SessionToken{
			Status: "Login success",
			Token:  jwtToken,
		})
	} else {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusOK, errors.New("email/password invalid"), "Login Failed. Email or password invalid")
		return
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	sessionId := userCtx.SessionID
	userId := userCtx.UserID

	err := dbHelper.DeleteUserSession(sessionId)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, err, fmt.Sprintf("Logout Failed. Failed to delete user session with userId %v", userId))
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusCreated, map[string]string{
		"status": "Logout success",
		"userId": userId,
	})
}
