package dbHelper

import (
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func GetRestaurantById(restaurantId string) (models.Restaurant, error) {
	sqlQuery := `select Id, name, address, latitude, longitude, createdBy, createdAt from public.restaurants 
					where Id=$1 and archivedAt is NULL`

	var restData models.Restaurant
	getErr := database.RmsDB.Get(&restData, sqlQuery, restaurantId)
	if getErr != nil {
		return restData, getErr
	}

	return restData, nil
}

func IsRestaurantExists(name, address string) (bool, error) {
	sqlQuery := `select count(Id) > 0 as isExists from public.restaurants where name=$1 and address=$2 and archivedAt is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, name, address)
	return exists, err
}

func CreateRestaurantHelper(name string, address string, latitude float64, longitude float64, createdBy string) (string, error) {
	var restaurantId uuid.UUID
	sqlQuery := `insert into public.restaurants (Id, name, address, latitude, longitude, createdBy, createdAt) 
					values ($1, $2, $3, $4, $5, $6, $7) returning Id;`

	crtErr := database.RmsDB.Get(&restaurantId, sqlQuery, uuid.New(), name, address, latitude, longitude, createdBy, time.Now())
	return restaurantId.String(), crtErr
}

func GetRestaurantHelper() ([]models.Restaurant, error) {
	sqlQuery := `select Id, name, address, latitude, longitude, createdBy, createdAt
					from public.restaurants where archivedAt is null`
	restData := make([]models.Restaurant, 0)
	err := database.RmsDB.Select(&restData, sqlQuery)

	return restData, err
}

func GetRestaurantSubAdminHelper(createdBy string) ([]models.Restaurant, error) {
	sqlQuery := `select Id, name, address, latitude, longitude, createdBy, createdAt
					from public.restaurants where createdBy=$1 and archivedAt is null`
	restData := make([]models.Restaurant, 0)
	err := database.RmsDB.Select(&restData, sqlQuery, createdBy)

	return restData, err
}

func GetDishesForRestaurantHelper(resData []models.Restaurant) ([]models.Restaurant, error) {
	sqlQuery := `select Id, name, price, restaurantId, createdAt from public.dishes 
					where restaurantId=$1 and archivedAt is NULL`
	var err error
	for i := range resData {
		dishData := make([]models.DishData, 0)
		err = database.RmsDB.Select(&dishData, sqlQuery, resData[i].Id)
		resData[i].Dishes = dishData
	}

	return resData, err
}
