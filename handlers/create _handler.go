package handlers

import (
	"net/http"

	dbHelper "github.com/gauraveg/rmsapp/database/dbhelper"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/models"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UserData
	userctx := middlewares.UserContext(r)
	createdBy := userctx.UserID
	//role := "user"

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	exist, err := dbHelper.IsUserExists(payload.Email)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding user")
		return
	}
	if exist {
		utils.ResponseWithError(w, http.StatusConflict, nil, "User Already Exists")
	}

	//Password hashing
	hashedPwd := utils.HashingPwd(payload.Password)

	userId, userEr := dbHelper.CreateUserHelper(payload.Email, payload.Name, hashedPwd, createdBy, payload.Role, payload.Address)
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
}

func CreateRestaurent(w http.ResponseWriter, r *http.Request) {
	var payload models.RestaurantsRequest
	userctx := middlewares.UserContext(r)
	createdBy := userctx.UserID

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	exist, err := dbHelper.IsRestaurentExists(payload.Name, payload.AddressLine)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, "Error while finding restaurent")
		return
	}
	if exist {
		utils.ResponseWithError(w, http.StatusConflict, nil, "Restaurent Already Exists")
	}

	restaurantId, resEr := dbHelper.CreateRestaurentHelper(payload.Name, payload.AddressLine, payload.Latitude, payload.Longitude, createdBy)
	if resEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, resEr, "Failed to add new Restaurent")
		return
	}

	var restaurent models.Restaurant
	restaurent, restEr := dbHelper.GetRestaurentById(restaurantId)
	if restEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, restEr, "Failed to create and fetch restaurent data")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, restaurent)
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var payload models.DishRequest
	restaurantId := chi.URLParam(r, "restaurantId")

	err := utils.ParsePayload(r.Body, &payload)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err, err.Error())
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
		utils.ResponseWithError(w, http.StatusInternalServerError, resEr, "Failed to add new Restaurent")
		return
	}

	var dish models.Dish
	dish, dishEr := dbHelper.GetDishById(dishId)
	if dishEr != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, dishEr, "Failed to create and fetch restaurent data")
		return
	}

	utils.ResponseWithJson(w, http.StatusCreated, dish)
}
