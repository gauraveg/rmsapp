package handlers

import (
	"errors"
	"fmt"
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UserData
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		zap.L().Error("Payload cannot be parsed. Check the payload",
			zap.Error(err))

		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	isValid := utils.CheckValidation(payload)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed")
		return
	}

	//Check if the user with the email exists in system or not
	exist, err := dbHelper.IsUserExists(payload.Email)
	if err != nil {
		utils.LogError("Error while finding user", err, "user's email", fmt.Sprintf("%#v", payload.Email))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding user")
		return
	}
	if exist {
		utils.LogError("User Already Exists", nil, "user's email", fmt.Sprintf("%#v", payload.Email))
		utils.ResponseWithError(w, http.StatusConflict, nil, "User Already Exists")
		return
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	userId, userEr := dbHelper.CreateUserHelper(payload.Email, payload.Name, hashedPwd, createdBy, payload.Role, payload.Addresses)
	if userEr != nil {
		utils.LogError("Failed to create new user", userEr, "user's email", fmt.Sprintf("%#v", payload.Email))
		utils.ResponseWithError(w, http.StatusInternalServerError, userEr, "Failed to create new user")
		return
	}

	var user models.User
	user, userErr := dbHelper.GetUserById(userId, payload.Role)
	if userErr != nil {
		utils.LogError("Failed to create and fetch sub admin user", userEr, "user's email", fmt.Sprintf("%#v", payload.Email))
		utils.ResponseWithError(w, http.StatusInternalServerError, userErr, "Failed to create and fetch sub admin user")
		return
	}
	utils.ResponseWithJson(w, http.StatusCreated, user)
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var payload models.RestaurantsRequest
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		zap.L().Error("Payload cannot be parsed. Check the payload",
			zap.Error(err))

		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload cannot be parsed. Check the payload")
		return
	}

	//Validator to check the payload's required fields
	isValid := utils.CheckValidation(payload)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed")
	}

	exist, err := dbHelper.IsRestaurantExists(payload.Name, payload.Address)
	if err != nil {
		utils.LogError("Error while finding restaurant", err, "restaurant's name", fmt.Sprintf("%#v", payload.Name))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding restaurant")
		return
	}
	if exist {
		utils.LogError("Restaurant Already Exists", nil, "restaurant's name", fmt.Sprintf("%#v", payload.Name))
		utils.ResponseWithError(w, http.StatusConflict, nil, "Restaurant Already Exists")
	}

	restaurantId, resEr := dbHelper.CreateRestaurantHelper(payload.Name, payload.Address, payload.Latitude, payload.Longitude, createdBy)
	if resEr != nil {
		utils.LogError("Failed to add new Restaurant", resEr, "restaurant's name", fmt.Sprintf("%#v", payload.Name))
		utils.ResponseWithError(w, http.StatusInternalServerError, resEr, "Failed to add new Restaurant")
		return
	}

	var restaurant models.Restaurant
	restaurant, restEr := dbHelper.GetRestaurantById(restaurantId)
	if restEr != nil {
		utils.LogError("Failed to create and fetch restaurant data", restEr, "restaurant's name", fmt.Sprintf("%#v", payload.Name))
		utils.ResponseWithError(w, http.StatusInternalServerError, restEr, "Failed to create and fetch restaurant data")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, restaurant)
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var payload models.DishRequest
	restaurantId := chi.URLParam(r, "restaurantId")
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		zap.L().Error("Payload cannot be parsed. Check the payload",
			zap.Error(err))
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	//Validator to check the payload's required fields
	isValid := utils.CheckValidation(payload)
	if !isValid {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Payload's required validation failed")
		return
	}

	//Check if restaurant is created by sub-admin
	restExist, err := dbHelper.GetRestaurantById(restaurantId)
	if err != nil {
		utils.LogError("Error while searching for restaurants", err, "dish's name", fmt.Sprintf("%#v", payload.Name))
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while searching for restaurants")
		return
	}

	if restExist.CreatedBy == createdBy {
		exist, err := dbHelper.IsDishExists(payload.Name, restaurantId)
		if err != nil {
			utils.LogError("Error while finding dishes", err, "dish's name", fmt.Sprintf("%#v", payload.Name))
			utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding dishes")
			return
		}
		if exist {
			utils.LogError("Dish Already Exists", nil, "dish's name", fmt.Sprintf("%#v", payload.Name))
			utils.ResponseWithError(w, http.StatusConflict, nil, "Dish Already Exists")
		}

		dishId, resEr := dbHelper.CreateDishHelper(payload.Name, payload.Price, restaurantId)
		if resEr != nil {
			utils.LogError("Failed to add new Restaurant", resEr, "dish's name", fmt.Sprintf("%#v", payload.Name))
			utils.ResponseWithError(w, http.StatusInternalServerError, resEr, "Failed to add new Restaurant")
			return
		}

		var dish models.Dish
		dish, dishEr := dbHelper.GetDishById(dishId)
		if dishEr != nil {
			utils.LogError("Failed to create and fetch restaurant data", dishEr, "dish's name", fmt.Sprintf("%#v", payload.Name))
			utils.ResponseWithError(w, http.StatusInternalServerError, dishEr, "Failed to create and fetch restaurant data")
			return
		}

		utils.ResponseWithJson(w, http.StatusCreated, dish)
	} else {
		utils.ResponseWithError(w, http.StatusBadRequest, errors.New("restaurant Id in URL param can't be accessed by sub-admin or does not exist"), "Restaurant can't be accessed by sub-admin or does not exist")
	}
}
