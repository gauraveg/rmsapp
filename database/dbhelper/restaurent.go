package dbHelper

import (
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func GetRestaurantById(restaurantId string) (models.Restaurant, error) {
	sqlQuery := `select restaurantid, name, addressline, latitude, longitude, createdby, createdat, archivedat from public.restaurants 
					where restaurantid=$1 and archivedat is NULL`

	var restdata models.Restaurant
	getErr := database.RmsDB.Get(&restdata, sqlQuery, restaurantId)
	if getErr != nil {
		return restdata, getErr
	}

	return restdata, nil
}

func IsRestaurantExists(name, addressline string) (bool, error) {
	sqlQuery := `select count(restaurantid) > 0 as isExists from public.restaurants where name=$1 and addressline=$2 and archivedat is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, name, addressline)
	return exists, err
}

func CreateRestaurantHelper(name string, addressline string, latitude float64, longitude float64, createdBy string) (string, error) {
	var restaurantId uuid.UUID
	sqlQuery := `insert into public.restaurants (restaurantid, name, addressline, latitude, longitude, createdby, createdat) 
					values ($1, $2, $3, $4, $5, $6, $7) returning restaurantid;`

	crtErr := database.RmsDB.Get(&restaurantId, sqlQuery, uuid.New(), name, addressline, latitude, longitude, createdBy, time.Now())
	return restaurantId.String(), crtErr
}

func GetRestaurantHelper() ([]models.Restaurant, error) {
	sqlquery := `select restaurantid, name, addressline, latitude, longitude, createdby, createdat, archivedat 
					from public.restaurants where archivedat is null`
	restdata := make([]models.Restaurant, 0)
	err := database.RmsDB.Select(&restdata, sqlquery)

	return restdata, err
}

func GetRestaurantSubAdminHelper(createdBy string) ([]models.Restaurant, error) {
	sqlquery := `select restaurantid, name, addressline, latitude, longitude, createdby, createdat, archivedat 
					from public.restaurants where createdby=$1 and archivedat is null`
	restdata := make([]models.Restaurant, 0)
	err := database.RmsDB.Select(&restdata, sqlquery, createdBy)

	return restdata, err
}
