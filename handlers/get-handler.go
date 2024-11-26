package handlers

import (
	"github.com/gauraveg/rmsapp/logger"
	"net/http"

	"github.com/gauraveg/rmsapp/models"

	"github.com/gauraveg/rmsapp/database"
	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// ----------------------------------------------------------------------------------------------------------
// BY ADMINS
func GetUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	role := string(models.RoleUser)
	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		userData, err := dbHelper.GetUserDataHelper(tx, role)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while fetch user", "error": err.Error()})
			return err
		}
		//Fetch Address for role as user
		userData, err = dbHelper.GetAddressForUserHelper(tx, userData)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch User's address", "error": err.Error()})
			return err
		}

		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, userData)
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": txErr.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}

func GetSubAdminsByAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	role := string(models.RoleSubAdmin)
	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		subAdmins, err := dbHelper.GetUserDataHelper(tx, role)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch sub-admins", "error": err.Error()})
			return err
		}
		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, subAdmins)
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": txErr.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}

func GetRestaurantsByAdminAndUser(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		restaurants, err := dbHelper.GetRestaurantByAdminAndUserHelper(tx)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurants", "error": err.Error()})
			return err
		}
		//Get the dishes for each restaurant
		restaurants, err = dbHelper.GetDishesForRestaurantHelper(tx, restaurants)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurant's dishes", "error": err.Error()})
			return err
		}
		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, restaurants)
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": txErr.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
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

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		userData, err := dbHelper.GetUsersSubAdminHelper(tx, role, createdBy)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch user", "error": err.Error()})
			return err
		}
		//Fetch Address for user role
		userData, err = dbHelper.GetAddressForUserHelper(tx, userData)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch user's address", "error": err.Error()})
			return err
		}
		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, userData)
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": txErr.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}

func GetRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		restaurants, err := dbHelper.GetRestaurantSubAdminHelper(tx, createdBy)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurants", "error": err.Error()})
			return err
		}
		//Get the dishes for each restaurant
		restaurants, err = dbHelper.GetDishesForRestaurantHelper(tx, restaurants)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurant's dishes", "error": err.Error()})
			return err
		}
		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, restaurants)
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": txErr.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
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

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		restPoint, err = dbHelper.GetRestLatitudeAndLongitude(tx, restaurantId)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch restaurant's latitude and longitude", "restaurantId": restaurantId, "error": err.Error()})
			return err
		}

		userPoint, err = dbHelper.GetUserLatitudeAndLongitude(tx, userId)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to fetch User's latitude and longitude", "userId": userId, "error": err.Error()})
			return err
		}
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "userId": userId, "error": txErr.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
	distance := utils.CalculateDistBetweenPoints(restPoint, userPoint)
	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusOK, distance)
}
