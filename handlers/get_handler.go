package handlers

import (
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
)

func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	role := "user"
	userdata, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userdata, err = dbHelper.GetAddressForUser(userdata)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	if len(userdata) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, userdata)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responsebody": "No record found"})
	}
}

func GetSubAdmins(w http.ResponseWriter, r *http.Request) {
	role := "sub-admin"
	subadmins, err := dbHelper.GetUsersHelper(role)
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

func GetRestaurentsByAdmin(w http.ResponseWriter, r *http.Request) {
	restaurents, err := dbHelper.GetRestaurentHelper()
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	if len(restaurents) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, restaurents)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responsebody": "No record found"})
	}
}

func GetUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	role := "user"
	userctx := middlewares.UserContext(r)
	createdBy := userctx.UserID

	userdata, err := dbHelper.GetUsersSubAdminHelper(role, createdBy)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userdata, err = dbHelper.GetAddressForUser(userdata)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	if len(userdata) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, userdata)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responsebody": "No record found"})
	}
}

func GetRestaurentsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userctx := middlewares.UserContext(r)
	createdBy := userctx.UserID
	restaurents, err := dbHelper.GetRestaurentSubAdminHelper(createdBy)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	if len(restaurents) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, restaurents)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responsebody": "No record found"})
	}
}
