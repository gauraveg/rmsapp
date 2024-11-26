package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gauraveg/rmsapp/logger"

	"github.com/gauraveg/rmsapp/database"
	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UserData
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Payload cannot be parsed. Check the payload"})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Check if the user with the email exists in system or not
	exist, err := dbHelper.IsUserExists(payload.Email)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while finding user", "email": payload.Email, "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding user")
		return
	}
	if exist {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "User Already Exists", "email": payload.Email})
		utils.ResponseWithError(w, http.StatusConflict, nil, "User Already Exists")
		return
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		_, userEr := dbHelper.CreateUserHelper(tx, payload.Email, payload.Name, hashedPwd, createdBy, string(payload.Role), payload.Addresses)
		if userEr != nil {
			return userEr
		}
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": txErr.Error()})
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"responseBody": fmt.Sprintf("User created with email %v", payload.Email),
	})
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var payload models.RestaurantsRequest
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction", "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	exist, err := dbHelper.IsRestaurantExists(payload.Name, payload.Address)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while finding restaurant", "name in payload": payload.Name, "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding restaurant")
		return
	}
	if exist {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Restaurant Already Exists", "name in payload": payload.Name})
		utils.ResponseWithError(w, http.StatusConflict, nil, "Restaurant Already Exists")
	}

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		_, resEr := dbHelper.CreateRestaurantHelper(tx, payload.Name, payload.Address, payload.Latitude, payload.Longitude, createdBy)
		if resEr != nil {
			return resEr
		}
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction. Error while creating restaurant", "name in payload": payload.Name, "error": txErr.Error()})
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed to add new Restaurant")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
		"responseBody": fmt.Sprintf("Restaurant created with name %v", payload.Name),
	})
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var payload models.DishRequest
	loggers := logger.GetLogContext(r)
	restaurantId := chi.URLParam(r, "restaurantId")
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Payload cannot be parsed. Check the payload", "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Check if restaurant is created by sub-admin
	restExist, err := dbHelper.GetRestaurantById(restaurantId)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while searching for restaurants", "name in payload": payload.Name, "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while searching for restaurants")
		return
	}

	roleCheck := string(models.RoleSubAdmin)
	if restExist.CreatedBy == createdBy || string(userCtx.Role) != roleCheck {
		exist, err := dbHelper.IsDishExists(payload.Name, restaurantId)
		if err != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while finding dishes", "name in payload": payload.Name, "error": err.Error()})
			utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding dishes")
			return
		}
		if exist {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Dish Already Exists", "name in payload": payload.Name})
			utils.ResponseWithError(w, http.StatusConflict, nil, "Dish Already Exists")
		}

		txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
			_, resEr := dbHelper.CreateDishHelper(tx, payload.Name, payload.Price, restaurantId)
			if resEr != nil {
				return resEr
			}
			return nil
		})
		if txErr != nil {
			loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction. Error while creating dish", "name in payload": payload.Name, "error": txErr.Error()})
			utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed to add new Restaurant")
			return
		}

		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
			"responseBody": fmt.Sprintf("Dish created with name %v", payload.Name),
		})
	} else {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Restaurant is forbidden for sub-admin or does not exist"})
		utils.ResponseWithError(w, http.StatusBadRequest, errors.New("restaurant Id in URL param is forbidden for sub-admin or does not exist"), "Restaurant is forbidden for sub-admin or does not exist")
	}
}

func AddAddressForUserByAdmins(w http.ResponseWriter, r *http.Request) {
	var payload models.AddressData
	userId := chi.URLParam(r, "userId")
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	role := userCtx.Role

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Payload cannot be parsed. Check the payload", "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, strings.Join(errMsg, "|"))
		return
	}

	//Check if the user with the email exists in system or not
	exist, err := dbHelper.IsAddressExists(payload.Address)
	if err != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while finding address", "address in payload": payload.Address, "error": err.Error()})
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding address")
		return
	}
	if exist {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Address already Exists", "address in payload": payload.Address, "error": err.Error()})
		utils.ResponseWithError(w, http.StatusConflict, nil, "Address already Exists")
		return
	}

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		addressId, userEr := dbHelper.CreateAddressForUserHelper(tx, payload.Address, payload.Latitude, payload.Longitude, userId, createdBy, string(role))
		if userEr != nil {
			return userEr
		}
		utils.ResponseWithJson(w, http.StatusCreated, map[string]string{
			"responseBody": fmt.Sprintf("Address added with Id: %v", addressId),
		})
		return nil
	})
	if txErr != nil {
		loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed in database transaction. Error while adding address", "address in payload": payload.Address, "error": txErr.Error()})
		utils.ResponseWithError(w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}
