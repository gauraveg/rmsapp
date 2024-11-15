package dbHelper

import (
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func GetDishById(dishId string) (models.Dish, error) {
	sqlQuery := `select dishid, name, price, restaurantid, createdat, archivedat from public.dishes 
					where dishid=$1 and archivedat is NULL`

	var dishdata models.Dish
	getErr := database.RmsDB.Get(&dishdata, sqlQuery, dishId)
	if getErr != nil {
		return dishdata, getErr
	}

	return dishdata, nil
}

func IsDishExists(name string, restaurantId string) (bool, error) {
	sqlQuery := `select count(dishid) > 0 as isExists from public.dishes where name=$1 and restaurantid=$2 and archivedat is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, name, restaurantId)
	return exists, err
}

func CreateDishHelper(name string, price int, restaurantId string) (string, error) {
	var dishId uuid.UUID
	sqlQuery := `insert into public.dishes (dishid, name, price, restaurantid, createdat) 
					values ($1, $2, $3, $4, $5) returning dishid;`

	crtErr := database.RmsDB.Get(&dishId, sqlQuery, uuid.New(), name, price, restaurantId, time.Now())
	return dishId.String(), crtErr
}

func GetAllDishHelper() ([]models.Dish, error) {
	sqlquery := `select dishid, name, price, restaurantid, createdat, archivedat 
					from public.dishes where archivedat is null`
	dishdata := make([]models.Dish, 0)
	err := database.RmsDB.Select(&dishdata, sqlquery)

	return dishdata, err
}

func GetAllDishSubAdminHelper(createdBy string) ([]models.Dish, error) {
	sqlquery := `select d.dishid, d.name, d.price, d.restaurantid, d.createdat 
					from public.dishes d INNER JOIN public.restaurants r
					on r.restaurantid = d.restaurantid 
					where d.archivedat is null and r.createdby = $1`
	dishdata := make([]models.Dish, 0)
	err := database.RmsDB.Select(&dishdata, sqlquery, createdBy)

	return dishdata, err
}
