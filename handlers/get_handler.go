package handlers

import (
	"net/http"

	"go.uber.org/zap"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
)

// ----------------------------------------------------------------------------------------------------------
// BY ADMINS
func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)

	role := "user"
	userData, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		logger.Error("Failed to fetch User", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUser(userData)
	if err != nil {
		logger.Error("Failed to fetch User's address", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, userData)
}

func GetSubAdminsByAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	role := "sub-admin"
	subAdmins, err := dbHelper.GetUsersHelper(role)
	if err != nil {
		logger.Error("Failed to fetch sub-admins", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch sub admin user")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, subAdmins)
}

func GetRestaurantsByAdminAndUser(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)

	restaurants, err := dbHelper.GetRestaurantHelper()
	if err != nil {
		logger.Error("Failed to fetch restaurants", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	//Get the dishes for each restaurant
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		logger.Error("Failed to fetch restaurant's dishes", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurant's dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, restaurants)
}

func GetAllDishesByAdminAndUser(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	dishes, err := dbHelper.GetAllDishHelper()
	if err != nil {
		logger.Error("Failed to fetch dishes", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, dishes)
}

// --------------------------------------------------------------------------------------------------------------------
// BY SUB-ADMINS
func GetUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	role := "user"
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	userData, err := dbHelper.GetUsersSubAdminHelper(role, createdBy)
	if err != nil {
		logger.Error("Failed to fetch user", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user")
		return
	}

	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUser(userData)
	if err != nil {
		logger.Error("Failed to fetch user's address", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch user's address")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, userData)
}

func GetRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	restaurants, err := dbHelper.GetRestaurantSubAdminHelper(createdBy)
	if err != nil {
		logger.Error("Failed to fetch restaurants", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	//Get the dishes for each restaurant
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		logger.Error("Failed to fetch restaurant's dishes", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurant's dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, restaurants)
}

func GetAllDishesBySubAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	dishes, err := dbHelper.GetAllDishSubAdminHelper(createdBy)
	if err != nil {
		logger.Error("Failed to fetch dishes", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, dishes)
}

// ---------------------------------------------------------------------------------------------------------------
// BY USERS
func GetRestaurantsByUser(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)

	restaurants, err := dbHelper.GetRestaurantHelper() //Get restaurants
	if err != nil {
		logger.Error("Failed to fetch restaurants", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurants")
		return
	}

	//Get the dishes for each restaurant
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		logger.Error("Failed to fetch restaurant's dishes", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch restaurant's dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, restaurants)
}

func GetDishesByRestId(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	restaurantId := chi.URLParam(r, "restaurantId")

	dishes, err := dbHelper.GetDishesByRestIdHelper(restaurantId)
	if err != nil {
		logger.Error("Failed to fetch dishes", zap.Error(err))
		utils.ResponseWithError(w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, dishes)
}

func DistanceBetweenCoords(w http.ResponseWriter, r *http.Request) {
	// todo
}
