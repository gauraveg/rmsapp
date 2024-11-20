package handlers

import (
	"fmt"
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
)

// ----------------------------------------------------------------------------------------------------------
// BY ADMINS
func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	email := userCtx.Email
	role := "user"
	userData, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		utils.LogError("Failed to fetch", err, "Admin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUser(userData)
	if err != nil {
		utils.LogError("Failed to fetch user's address", err, "Admin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	if len(userData) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, userData)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

func GetSubAdminsByAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	email := userCtx.Email
	role := "sub-admin"
	subAdmins, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		utils.LogError("Failed to fetch", err, "Admin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	if len(subAdmins) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, subAdmins)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

func GetRestaurantsByAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	email := userCtx.Email
	restaurants, err := dbHelper.GetRestaurantHelper()
	if err != nil {
		utils.LogError("Failed to fetch restaurants", err, "Admin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	//Get the dishes for each restaurants
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		utils.LogError("Failed to fetch restaurant's dishes", err, "Admin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurant's dishes")
		return
	}

	if len(restaurants) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, restaurants)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

func GetAllDishesByAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	email := userCtx.Email
	dishes, err := dbHelper.GetAllDishHelper()
	if err != nil {
		utils.LogError("Failed to fetch dishes", err, "SubAdmin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	if len(dishes) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, dishes)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

// --------------------------------------------------------------------------------------------------------------------
// BY SUB-ADMINS
func GetUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	role := "user"
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	email := userCtx.Email

	userData, err := dbHelper.GetUsersSubAdminHelper(role, createdBy)
	if err != nil {
		utils.LogError("Failed to fetch user", err, "SubAdmin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUser(userData)
	if err != nil {
		utils.LogError("Failed to fetch user's address", err, "SubAdmin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	if len(userData) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, userData)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

func GetRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	email := userCtx.Email
	restaurants, err := dbHelper.GetRestaurantSubAdminHelper(createdBy)
	if err != nil {
		utils.LogError("Failed to fetch restaurants", err, "SubAdmin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	//Get the dishes for each restaurants
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		utils.LogError("Failed to fetch restaurant's dishes", err, "SubAdmin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurant's dishes")
		return
	}

	if len(restaurants) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, restaurants)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

func GetAllDishesBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	email := userCtx.Email

	dishes, err := dbHelper.GetAllDishSubAdminHelper(createdBy)
	if err != nil {
		utils.LogError("Failed to fetch dishes", err, "SubAdmin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	if len(dishes) > 0 {
		utils.ResponseWithJson(w, http.StatusOK, dishes)
	} else {
		utils.ResponseWithJson(w, http.StatusOK, map[string]string{"responseBody": "No record found"})
	}
}

// ---------------------------------------------------------------------------------------------------------------
// BY USERS
func GetRestaurantsByUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	//createdBy := userCtx.UserID
	email := userCtx.Email

	restaurants, err := dbHelper.GetRestaurantHelper() //Get restaurants
	if err != nil {
		utils.LogError("Failed to fetch restaurants", err, "User", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	//Get the dishes for each restaurants
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		utils.LogError("Failed to fetch restaurant's dishes", err, "Admin", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurant's dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, restaurants)
}

func GetAllDishesByUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	email := userCtx.Email

	dishes, err := dbHelper.GetAllDishesByUserHelper()
	if err != nil {
		utils.LogError("Failed to fetch dishes", err, "User", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, dishes)
}

func GetDishesByRestId(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	email := userCtx.Email
	restaurantId := chi.URLParam(r, "restaurantId")

	dishes, err := dbHelper.GetDishesByRestIdHelper(restaurantId)
	if err != nil {
		utils.LogError("Failed to fetch dishes", err, "User", fmt.Sprintf("%#v", email))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, dishes)
}
