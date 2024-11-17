package handlers

import (
	"errors"
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UserData
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID
	//role := "user"

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed.")
		return
	}

	//payload validations
	valid := utils.ValidateUserPayload(payload)

	if valid {
		//Validate address. If invalid, remove that particular address from payload
		payload = utils.ValidateAddress(payload)
		exist, err := dbHelper.IsUserExists(payload.Email)
		if err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding user")
			return
		}
		if exist {
			utils.ResponseWithError(w, http.StatusConflict, nil, "User Already Exists")
			return
		}

		//Password hashing
		hashedPwd := utils.HashingPwd(payload.Password)

		userId, userEr := dbHelper.CreateUserHelper(payload.Email, payload.Name, hashedPwd, createdBy, payload.Role, payload.Addresses)
		if userEr != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, userEr, "Failed to create new user")
			return
		}

		var user models.User
		user, userErr := dbHelper.GetUserById(userId, payload.Role)
		if userErr != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, userErr, "Failed to create and fetch sub admin user")
			return
		}

		utils.ResponseWithJson(w, http.StatusCreated, user)
	} else {
		utils.ResponseWithError(w, http.StatusBadRequest, errors.New("payload has incorrect data"), "Failed as user data has incorrect details")
	}
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var payload models.RestaurantsRequest
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed.")
		return
	}

	exist, err := dbHelper.IsRestaurantExists(payload.Name, payload.Address)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding restaurant")
		return
	}
	if exist {
		utils.ResponseWithError(w, http.StatusConflict, nil, "Restaurant Already Exists")
	}

	restaurantId, resEr := dbHelper.CreateRestaurantHelper(payload.Name, payload.Address, payload.Latitude, payload.Longitude, createdBy)
	if resEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, resEr, "Failed to add new Restaurant")
		return
	}

	var restaurant models.Restaurant
	restaurant, restEr := dbHelper.GetRestaurantById(restaurantId)
	if restEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, restEr, "Failed to create and fetch restaurant data")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, restaurant)
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var payload models.DishRequest
	restaurantId := chi.URLParam(r, "restaurantId")

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	//Validator to check the payload's required fields
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed.")
		return
	}

	exist, err := dbHelper.IsDishExists(payload.Name, restaurantId)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding dishes")
		return
	}
	if exist {
		utils.ResponseWithError(w, http.StatusConflict, nil, "Dish Already Exists")
	}

	dishId, resEr := dbHelper.CreateDishHelper(payload.Name, payload.Price, restaurantId)
	if resEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, resEr, "Failed to add new Restaurant")
		return
	}

	var dish models.Dish
	dish, dishEr := dbHelper.GetDishById(dishId)
	if dishEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, dishEr, "Failed to create and fetch restaurant data")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, dish)
}
