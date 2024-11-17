package handlers

import (
	"errors"
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-playground/validator/v10"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginRequest
	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	//Validator to check the payload's required fields
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed.")
		return
	}

	userID, pwdHash, role, userErr := dbHelper.GetUserInfo(payload)
	if userErr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, "Failed to find user")
		return
	}

	if userID == "" || utils.VerifyPwdHash(payload.Password, pwdHash) {
		sessionID, crtErr := dbHelper.CreateUserSession(userID)
		if crtErr != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, crtErr, "Failed to create user session")
			return
		}

		jwtToken, jwtErr := utils.GenerateJwt(userID, role, sessionID)
		if jwtErr != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, jwtErr, "Failed to generate Tokens")
			return
		}

		utils.ResponseWithJson(w, http.StatusOK, models.SessionToken{
			Status: "Login success",
			Token:  jwtToken,
		})
	} else {
		utils.ResponseWithError(w, http.StatusOK, errors.New("password invalid"), "Login Failed. Check email or password")
		return
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	sessionId := userCtx.SessionID
	userId := userCtx.UserID

	err := dbHelper.DeleteUserSession(sessionId)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to delete user session")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"status": "Logout success",
		"userId": userId,
	})
}
