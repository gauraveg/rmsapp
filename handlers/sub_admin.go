package handlers

import (
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
)

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	var payload models.SubAdminRequest
	userctx := middlewares.UserContext(r)
	createdBy := userctx.UserID
	role := "sub-admin"

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

	userId, userEr := dbHelper.CreateSubAdminHelper(payload.Email, payload.Name, hashedPwd, createdBy, role)
	if userEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, userEr, "Failed to create sub admin user")
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

func GetSubAdmins(w http.ResponseWriter, r *http.Request) {
	subadmins, err := dbHelper.GetSubAdminsHelper()
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	if len(subadmins) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, subadmins)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responsebody": "No record found"})
	}
}
