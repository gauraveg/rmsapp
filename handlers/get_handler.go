package handlers

import (
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
)

func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	role := "user"
	userData, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUser(userData)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	if len(userData) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, userData)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}

func GetSubAdmins(w http.ResponseWriter, r *http.Request) {
	role := "sub-admin"
	subAdmins, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	if len(subAdmins) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, subAdmins)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}

func GetRestaurantsByAdmin(w http.ResponseWriter, r *http.Request) {
	restaurants, err := dbHelper.GetRestaurantHelper()
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	if len(restaurants) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, restaurants)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}

func GetUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	role := "user"
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	userData, err := dbHelper.GetUsersSubAdminHelper(role, createdBy)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUser(userData)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	if len(userData) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, userData)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}

func GetRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	restaurants, err := dbHelper.GetRestaurantSubAdminHelper(createdBy)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	if len(restaurants) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, restaurants)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}

func GetAllDishesByAdmin(w http.ResponseWriter, r *http.Request) {
	dishes, err := dbHelper.GetAllDishHelper()
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	if len(dishes) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, dishes)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}

func GetAllDishesBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	dishes, err := dbHelper.GetAllDishSubAdminHelper(createdBy)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	if len(dishes) > 0 {
		utils.ResponseWithJson(w, http.StatusCreated, dishes)
	} else {
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{"responseBody": "No record found"})
	}
}
