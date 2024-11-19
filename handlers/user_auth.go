package handlers

import (
	"errors"
	"fmt"
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginRequest
	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		zap.L().Error("Payload cannot be parsed. Check the payload",
			zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "cannot parse payload data")
		return
	}

	//Validator to check the payload's required fields
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		utils.LogError("Payload's required validation failed.", err, "payload", fmt.Sprintf("%#v", payload))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed")
		return
	}

	userID, pwdHash, role, userErr := dbHelper.GetUserInfo(payload)
	if userErr != nil {
		utils.LogError("Failed to find user", err, "payload", fmt.Sprintf("%#v", payload))
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, "Failed to find user")
		return
	}

	if userID == "" || utils.VerifyPwdHash(payload.Password, pwdHash) {
		sessionID, crtErr := dbHelper.CreateUserSession(userID)
		if crtErr != nil {
			utils.LogError("Failed to create user session", err, "payload", fmt.Sprintf("%#v", payload))
			utils.ResponseWithError(w, http.StatusInternalServerError, crtErr, "Failed to create user session")
			return
		}

		jwtToken, jwtErr := utils.GenerateJwt(userID, role, sessionID)
		if jwtErr != nil {
			utils.LogError("Failed to generate JWT Tokens", err, "payload", fmt.Sprintf("%#v", payload))
			utils.ResponseWithError(w, http.StatusInternalServerError, jwtErr, "Failed to generate Tokens")
			return
		}

		utils.ResponseWithJson(w, http.StatusOK, models.SessionToken{
			Status: "Login success",
			Token:  jwtToken,
		})
	} else {
		utils.LogError("Login Failed. Email or password invalid", errors.New("password invalid"), "payload", fmt.Sprintf("%#v", payload))
		utils.ResponseWithError(w, http.StatusOK, errors.New("password invalid"), "Login Failed. Email or password invalid")
		return
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	sessionId := userCtx.SessionID
	userId := userCtx.UserID

	err := dbHelper.DeleteUserSession(sessionId)
	if err != nil {
		utils.LogError("Failed to delete user session", err, "userId", fmt.Sprintf("%#v", userId))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to delete user session")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"status": "Logout success",
		"userId": userId,
	})
}
