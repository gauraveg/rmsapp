package handlers

import (
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginRequest
	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	userID, role, userErr := dbHelper.GetUserInfo(payload)
	if userErr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, "Failed to find user")
		return
	}

	if userID == "" || role == "" {
		utils.ResponseWithError(w, http.StatusOK, nil, "user not found")
		return
	}

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

	utils.ResponseWithJson(w, http.StatusCreated, models.SessionToken{
		Status: "Login success",
		Token:  jwtToken,
	})
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	userctx := middlewares.UserContext(r)
	sessionId := userctx.SessionID
	userId := userctx.UserID

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
