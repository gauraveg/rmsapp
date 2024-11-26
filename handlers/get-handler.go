package handlers

import (
	"github.com/gauraveg/rmsapp/logger"
	"net/http"

	"github.com/gauraveg/rmsapp/models"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
)

// ----------------------------------------------------------------------------------------------------------
// BY ADMINS
func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	role := string(models.RoleUser)

	userData, err := dbHelper.GetUserDataHelper(role)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while fetch user", "error": err.Error()})
		return
	}
	//Fetch Address for role as user
	userData, err = dbHelper.GetAddressForUserHelper(userData)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch User's address", "error": err.Error()})
		return
	}
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, userData)
}

func GetSubAdminsByAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	role := string(models.RoleSubAdmin)

	subAdmins, err := dbHelper.GetUserDataHelper(role)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch sub-admins", "error": err.Error()})
		return
	}
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, subAdmins)
}

func GetRestaurantsByAdminAndUser(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)

	restaurants, err := dbHelper.GetRestaurantByAdminAndUserHelper()
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurants", "error": err.Error()})
		return
	}
	//Get the dishes for each restaurant
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurant's dishes", "error": err.Error()})
		return
	}
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, restaurants)
}

func GetAllDishesByAdminAndUser(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	dishes, err := dbHelper.GetAllDishHelper()
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch dishes", "error": err.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, dishes)
}

// --------------------------------------------------------------------------------------------------------------------
// BY SUB-ADMINS
func GetUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	role := string(models.RoleUser)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	userData, err := dbHelper.GetUsersSubAdminHelper(role, createdBy)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch user", "error": err.Error()})
		return
	}
	//Fetch Address for user role
	userData, err = dbHelper.GetAddressForUserHelper(userData)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch user's address", "error": err.Error()})
		return
	}
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, userData)
}

func GetRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	restaurants, err := dbHelper.GetRestaurantSubAdminHelper(createdBy)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurants", "error": err.Error()})
		return
	}
	//Get the dishes for each restaurant
	restaurants, err = dbHelper.GetDishesForRestaurantHelper(restaurants)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurant's dishes", "error": err.Error()})
		return
	}
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, restaurants)
}

func GetAllDishesBySubAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	dishes, err := dbHelper.GetAllDishSubAdminHelper(createdBy)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch dishes", "error": err.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, dishes)
}

// ---------------------------------------------------------------------------------------------------------------
// BY USERS

func GetDishesByRestId(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	restaurantId := chi.URLParam(r, "restaurantId")

	dishes, err := dbHelper.GetDishesByRestIdHelper(restaurantId)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch dishes", "error": err.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, err, "Failed to fetch dishes")
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, dishes)
}

func DistanceBetweenCoords(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	restaurantId := chi.URLParam(r, "restaurantId")
	userCtx := middlewares.UserContext(r)
	userId := userCtx.UserID
	var err error
	var restPoint []models.Coordinates
	var userPoint []models.Coordinates

	restPoint, err = dbHelper.GetRestLatitudeAndLongitude(restaurantId)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurant's latitude and longitude", "restaurantId": restaurantId, "error": err.Error()})
		return
	}
	userPoint, err = dbHelper.GetUserLatitudeAndLongitude(userId)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch User's latitude and longitude", "userId": userId, "error": err.Error()})
		return
	}

	distance := utils.CalculateDistBetweenPoints(restPoint, userPoint)
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, distance)
}
