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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UserData
	userctx := middlewares.UserContext(r)
	createdBy := userctx.UserID
	role := "user"

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	exist, err := dbHelper.IsUserExists(payload.Email)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error whiile finding user")
		return
	}
	if exist {
		utils.ResponseWithError(w, http.StatusConflict, nil, "User Already Exists")
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	userId, userEr := dbHelper.CreateUserHelper(payload.Email, payload.Name, hashedPwd, createdBy, role, payload.Address)
	if userEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, userEr, "Failed to create new user")
		return
	}

	var user models.User
	user, userErr := dbHelper.GetUserById(userId, role)
	if userErr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, "Failed to create and fetch sub admin user")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, user)
}

func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	userdata, err := dbHelper.GetUsersByAdminHelper()
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	if len(userdata) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, userdata)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responsebody": "No record found"})
	}
}
