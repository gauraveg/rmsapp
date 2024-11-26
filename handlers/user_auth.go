package handlers

import (
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
	err := utils.ParsePayload(r.Body, &signUpData)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Payload cannot be parsed. Check the payload", "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	payload := models.SignUpWithRole{
		UserSignUp: signUpData,
		Role:       models.RoleUser,
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	userId, userEr := dbHelper.CreateSignUpHelper(payload.Email, payload.Name, hashedPwd, string(payload.Role), payload.Addresses)
	if userEr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to create new user", "email": payload.Email, "error": userEr.Error()})
		utils.ResponseWithError(w, http.StatusInternalServerError, userEr, "Failed to create new user")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, userId)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	var payload models.LoginRequest
	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to create new user", "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	userID, pwdHash, role, userErr := dbHelper.GetUserInfoForLogin(payload)
	if userErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to find user", "email": payload.Email, "error": userErr.Error()})
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, fmt.Sprintf("Failed to find user with email %v", payload.Email))
		return
	}

	if userID == "" || utils.VerifyPwdHash(payload.Password, pwdHash) {
		sessionID, crtErr := dbHelper.CreateUserSession(userID)
		if crtErr != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to create user session", "userId": userID, "error": crtErr.Error()})
			utils.ResponseWithError(w, http.StatusInternalServerError, crtErr, "Failed to create user session")
			return
		}

		jwtToken, jwtErr := utils.GenerateJwt(userID, role, sessionID)
		if jwtErr != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to generate JWT Tokens for the sessionID", "SessionID": sessionID, "error": jwtErr.Error()})
			utils.ResponseWithError(w, http.StatusInternalServerError, jwtErr, fmt.Sprintf("Failed to generate JWT Tokens for the sessionID %v", sessionID))
			return
		}

		utils.ResponseWithJson(w, http.StatusOK, models.SessionToken{
			Status: "Login success",
			Token:  jwtToken,
		})
	} else {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Login Failed. Email or password invalid"})
		utils.ResponseWithError(w, http.StatusOK, errors.New("email/password invalid"), "Login Failed. Email or password invalid")
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
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to delete user session", "userId": userId, "error": err.Error()})
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Logout Failed")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"status": "Logout success",
		"userId": userId,
	})
}
