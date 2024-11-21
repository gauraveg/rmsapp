package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gauraveg/rmsapp/database"
	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UserData
	logger := middlewares.LoggerContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		logger.Error("Payload cannot be parsed. Check the payload", zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(payload, logger)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Check if the user with the email exists in system or not
	exist, err := dbHelper.IsUserExists(payload.Email)
	if err != nil {
		logger.Error("Error while finding user", zap.String("email", payload.Email), zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding user")
		return
	}
	if exist {
		logger.Error("User Already Exists", zap.String("email", payload.Email), zap.Error(err))
		utils.ResponseWithError(w, http.StatusConflict, nil, "User Already Exists")
		return
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		_, userEr := dbHelper.CreateUserHelper(tx, payload.Email, payload.Name, hashedPwd, createdBy, payload.Role, payload.Addresses)
		if userEr != nil {
			return userEr
		}
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction", zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"responseBody": fmt.Sprintf("User created with email %v", payload.Email),
	})
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var payload models.RestaurantsRequest
	logger := middlewares.LoggerContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		logger.Error("Payload cannot be parsed. Check the payload", zap.Error(err))

		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(payload, logger)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	exist, err := dbHelper.IsRestaurantExists(payload.Name, payload.Address)
	if err != nil {
		logger.Error("Error while finding restaurant", zap.String("restaurant's name", fmt.Sprintf("%#v", payload.Name)), zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding restaurant")
		return
	}
	if exist {
		logger.Error("Restaurant Already Exists", zap.String("restaurant's name", fmt.Sprintf("%#v", payload.Name)))
		utils.ResponseWithError(w, http.StatusConflict, nil, "Restaurant Already Exists")
	}

	txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
		_, resEr := dbHelper.CreateRestaurantHelper(tx, payload.Name, payload.Address, payload.Latitude, payload.Longitude, createdBy)
		if resEr != nil {
			return resEr
		}
		return nil
	})
	if txErr != nil {
		logger.Error("Failed in database transaction. Error while finding restaurant", zap.String("restaurant's name", fmt.Sprintf("%#v", payload.Name)), zap.Error(txErr))
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed to add new Restaurant")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"responseBody": fmt.Sprintf("Restaurant created with name %v", payload.Name),
	})
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var payload models.DishRequest
	logger := middlewares.LoggerContext(r)
	restaurantId := chi.URLParam(r, "restaurantId")
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		logger.Error("Payload cannot be parsed. Check the payload", zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(payload, logger)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Check if restaurant is created by sub-admin
	restExist, err := dbHelper.GetRestaurantById(restaurantId)
	if err != nil {
		logger.Error("Error while searching for restaurants", zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while searching for restaurants")
		return
	}

	if restExist.CreatedBy == createdBy {
		exist, err := dbHelper.IsDishExists(payload.Name, restaurantId)
		if err != nil {
			logger.Error("Error while finding dishes", zap.Error(err))
			utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding dishes")
			return
		}
		if exist {
			logger.Error("Dish Already Exists", zap.String("dish's name", fmt.Sprintf("%#v", payload.Name)))
			utils.ResponseWithError(w, http.StatusConflict, nil, "Dish Already Exists")
		}

		txErr := database.WithTxn(logger, func(tx *sqlx.Tx) error {
			_, resEr := dbHelper.CreateDishHelper(tx, payload.Name, payload.Price, restaurantId)
			if resEr != nil {
				return resEr
			}
			return nil
		})
		if txErr != nil {
			logger.Error("Failed in database transaction. Failed to add new dish", zap.String("dish's name", fmt.Sprintf("%#v", payload.Name)), zap.Error(txErr))
			utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed to add new Restaurant")
			return
		}

		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
			"responseBody": fmt.Sprintf("Dish created with name %v", payload.Name),
		})
	} else {
		logger.Error("Restaurant is forbidden for sub-admin or does not exist", zap.String("Issue", "restaurant Id in URL param is forbidden for sub-admin or does not exist"))
		utils.ResponseWithError(w, http.StatusBadRequest, errors.New("restaurant Id in URL param is forbidden for sub-admin or does not exist"), "Restaurant is forbidden for sub-admin or does not exist")
	}
}
