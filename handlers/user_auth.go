package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"go.uber.org/zap"
)

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)

	var signUpData models.UserSignUp
	err := utils.ParsePayload(r.Body, &signUpData)
	if err != nil {
		logger.Error("Payload cannot be parsed. Check the payload", zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	payload := models.SignUpWithRole{
		UserSignUp: signUpData,
		Role:       "user",
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(payload, logger)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	userId, userEr := dbHelper.CreateSignUpHelper(payload.Email, payload.Name, hashedPwd, string(payload.Role), payload.Addresses)
	if userEr != nil {
		logger.Error("Failed to create new user", zap.String("email", payload.Email), zap.Error(userEr))
		utils.ResponseWithError(w, http.StatusInternalServerError, userEr, "Failed to create new user")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, userId)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	var payload models.LoginRequest
	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		logger.Error("Payload cannot be parsed. Check the payload", zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(payload, logger)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	userID, pwdHash, role, userErr := dbHelper.GetUserInfoForLogin(payload)
	if userErr != nil {
		logger.Error("Failed to find user", zap.String("Email", payload.Email))
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, fmt.Sprintf("Failed to find user with email %v", payload.Email))
		return
	}

	if userID == "" || utils.VerifyPwdHash(payload.Password, pwdHash) {
		sessionID, crtErr := dbHelper.CreateUserSession(userID)
		if crtErr != nil {
			logger.Error("Failed to create user session", zap.String("UserId", userID))
			utils.ResponseWithError(w, http.StatusInternalServerError, crtErr, "Failed to create user session")
			return
		}

		jwtToken, jwtErr := utils.GenerateJwt(userID, role, sessionID)
		if jwtErr != nil {
			logger.Error("Failed to generate JWT Tokens for the sessionID", zap.String("SessionID", sessionID))
			utils.ResponseWithError(w, http.StatusInternalServerError, jwtErr, fmt.Sprintf("Failed to generate JWT Tokens for the sessionID %v", sessionID))
			return
		}

		utils.ResponseWithJson(w, http.StatusOK, models.SessionToken{
			Status: "Login success",
			Token:  jwtToken,
		})
	} else {
		logger.Error("Login Failed. Email or password invalid")
		utils.ResponseWithError(w, http.StatusOK, errors.New("email/password invalid"), "Login Failed. Email or password invalid")
		return
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	userCtx := middlewares.UserContext(r)
	sessionId := userCtx.SessionID
	userId := userCtx.UserID

	err := dbHelper.DeleteUserSession(sessionId)
	if err != nil {
		logger.Error("Failed to delete user session", zap.String("userId", fmt.Sprintf("%#v", userId)))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Logout Failed")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"status": "Logout success",
		"userId": userId,
	})
}
