package handlers

import (
	"encoding/json"
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

	body, ok := r.Context().Value("payload").(string)
	if !ok {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload not present"), "cannot parse payload data")
		return
	}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload validation failed"), strings.Join(errMsg, "|"))
		return
	}

	//Check if the user with the email exists in system or not
	exist, err := dbHelper.IsUserExists(payload.Email)
	if err != nil {
		//loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Error while finding user", "email": payload.Email, "error": err.Error()})
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, fmt.Sprintf("Error while finding user with email %s", payload.Email))
		return
	}
	if exist {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusConflict, nil, fmt.Sprintf("User Already Exists with email %s", payload.Email))
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
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusCreated, map[string]string{
		"responseBody": fmt.Sprintf("User created with email %v", payload.Email),
	})
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var payload models.RestaurantsRequest
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	body, ok := r.Context().Value("payload").(string)
	if !ok {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload not present"), "cannot parse payload data")
		return
	}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload validation failed"), strings.Join(errMsg, "|"))
		return
	}

	exist, err := dbHelper.IsRestaurantExists(payload.Name, payload.Address)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, fmt.Sprintf("Error while finding restaurant with name %v", payload.Name))
		return
	}
	if exist {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusConflict, nil, fmt.Sprintf("Restaurant Already Exists with name %v", payload.Name))
	}

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		_, resEr := dbHelper.CreateRestaurantHelper(tx, payload.Name, payload.Address, payload.Latitude, payload.Longitude, createdBy)
		if resEr != nil {
			return resEr
		}
		return nil
	})
	if txErr != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, fmt.Sprintf("Error while finding restaurant with name %v", payload.Name))
		return
	}

	utils.ResponseWithJson(r.Context(), loggers, w, http.StatusCreated, map[string]string{
		"responseBody": fmt.Sprintf("Restaurant created with name %v", payload.Name),
	})
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var payload models.DishRequest
	loggers := logger.GetLogContext(r)
	restaurantId := chi.URLParam(r, "restaurantId")
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	body, ok := r.Context().Value("payload").(string)
	if !ok {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload not present"), "cannot parse payload data")
		return
	}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, err.Error())
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload validation failed"), strings.Join(errMsg, "|"))
		return
	}

	//Check if restaurant is created by sub-admin
	restExist, err := dbHelper.GetRestaurantById(restaurantId)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, fmt.Sprintf("Error while finding restaurant with name %v", payload.Name))
		return
	}

	roleCheck := string(models.RoleSubAdmin)
	if restExist.CreatedBy == createdBy || string(userCtx.Role) != roleCheck {
		exist, err := dbHelper.IsDishExists(payload.Name, restaurantId)
		if err != nil {
			utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, fmt.Sprintf("Error while finding dish with name %v", payload.Name))
			return
		}
		if exist {
			utils.ResponseWithError(r.Context(), loggers, w, http.StatusConflict, nil, fmt.Sprintf("Dish already exist with name %v", payload.Name))
		}

		txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
			_, resEr := dbHelper.CreateDishHelper(tx, payload.Name, payload.Price, restaurantId)
			if resEr != nil {
				return resEr
			}
			return nil
		})
		if txErr != nil {
			utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, fmt.Sprintf("Failed to add new Restaurant with name %v", payload.Name))
			return
		}

		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusCreated, map[string]string{
			"responseBody": fmt.Sprintf("Dish created with name %v", payload.Name),
		})
	} else {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("restaurant Id in URL param is forbidden for sub-admin or does not exist"), "Restaurant is forbidden for sub-admin or does not exist")
	}
}

func AddAddressForUserByAdmins(w http.ResponseWriter, r *http.Request) {
	var payload models.AddressData
	userId := chi.URLParam(r, "userId")
	loggers := logger.GetLogContext(r)
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	role := userCtx.Role

	body, ok := r.Context().Value("payload").(string)
	if !ok {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload not present"), "cannot parse payload data")
		return
	}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	errMsg, isValid := utils.CheckValidation(r.Context(), payload, loggers)
	if !isValid {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, errors.New("payload validation failed"), strings.Join(errMsg, "|"))
		return
	}

	//Check if the user with the email exists in system or not
	exist, err := dbHelper.IsAddressExists(payload.Address)
	if err != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, err, fmt.Sprintf("Error while finding address with addressline %v", payload.Address))
		return
	}
	if exist {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusBadRequest, nil, fmt.Sprintf("Address already exist with addressline %v", payload.Address))
		return
	}

	txErr := database.WithTxn(r.Context(), loggers, func(tx *sqlx.Tx) error {
		addressId, userEr := dbHelper.CreateAddressForUserHelper(tx, payload.Address, payload.Latitude, payload.Longitude, userId, createdBy, string(role))
		if userEr != nil {
			return userEr
		}
		utils.ResponseWithJson(r.Context(), loggers, w, http.StatusCreated, map[string]string{
			"responseBody": fmt.Sprintf("Address added with Id: %v", addressId),
		})
		return nil
	})
	if txErr != nil {
		utils.ResponseWithError(r.Context(), loggers, w, http.StatusInternalServerError, txErr, "Failed in database transaction")
		return
	}
}
