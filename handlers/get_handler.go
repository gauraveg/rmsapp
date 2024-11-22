package handlers

import (
	"net/http"

	"go.uber.org/zap"

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
	logger := middlewares.LoggerContext(r)
	role := "user"
	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		userData, err := dbHelper.GetUserDataHelper(tx, role)
		if err != nil {
			logger.Error("Failed to fetch User", zap.Error(err))
			return err
		}
		//Fetch Address for role as user
		userData, err = dbHelper.GetAddressForUserHelper(tx, userData)
		if err != nil {
			logger.Error("Failed to fetch User's address", zap.Error(err))
			return err
		}

		utils.ResponseWithJson(w, http.StatusOK, userData)
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}

func GetSubAdminsByAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	role := "sub-admin"
	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		subAdmins, err := dbHelper.GetUserDataHelper(tx, role)
		if err != nil {
			logger.Error("Failed to fetch sub-admins", zap.Error(err))
			return err
		}
		utils.ResponseWithJson(w, http.StatusOK, subAdmins)
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}

func GetRestaurantsByAdminAndUser(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		restaurants, err := dbHelper.GetRestaurantByAdminAndUserHelper(tx)
		if err != nil {
			logger.Error("Failed to fetch restaurants", zap.Error(err))
			return err
		}
		//Get the dishes for each restaurant
		restaurants, err = dbHelper.GetDishesForRestaurantHelper(tx, restaurants)
		if err != nil {
			logger.Error("Failed to fetch restaurant's dishes", zap.Error(err))
			return err
		}
		utils.ResponseWithJson(w, http.StatusOK, restaurants)
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
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

	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		userData, err := dbHelper.GetUsersSubAdminHelper(tx, role, createdBy)
		if err != nil {
			logger.Error("Failed to fetch user", zap.Error(err))
			return err
		}
		//Fetch Address for user role
		userData, err = dbHelper.GetAddressForUserHelper(tx, userData)
		if err != nil {
			logger.Error("Failed to fetch user's address", zap.Error(err))
			return err
		}
		utils.ResponseWithJson(w, http.StatusOK, userData)
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}

func GetRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.LoggerContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		restaurants, err := dbHelper.GetRestaurantSubAdminHelper(tx, createdBy)
		if err != nil {
			logger.Error("Failed to fetch restaurants", zap.Error(err))
			return err
		}
		//Get the dishes for each restaurant
		restaurants, err = dbHelper.GetDishesForRestaurantHelper(tx, restaurants)
		if err != nil {
			logger.Error("Failed to fetch restaurant's dishes", zap.Error(err))
			return err
		}
		utils.ResponseWithJson(w, http.StatusOK, restaurants)
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
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

	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		restaurants, err := dbHelper.GetRestaurantByAdminAndUserHelper(tx) //Get restaurants
		if err != nil {
			logger.Error("Failed to fetch restaurants", zap.Error(err))
			return err
		}
		//Get the dishes for each restaurant
		restaurants, err = dbHelper.GetDishesForRestaurantHelper(tx, restaurants)
		if err != nil {
			logger.Error("Failed to fetch restaurant's dishes", zap.Error(err))
			return err
		}

		utils.ResponseWithJson(w, http.StatusOK, restaurants)
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
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
