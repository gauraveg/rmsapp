package dbHelper

import (
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func GetDishById(dishId string) (models.Dish, error) {
	sqlQuery := `select d.Id, d.name, d.price, d.restaurantId, r.name as restaurantName, d.createdAt
					from public.dishes d inner join public.restaurants r
					on r.Id = d.restaurantId 
					where d.Id = $1 and d.archivedAt is null`

	var dishData models.Dish
	getErr := database.RmsDB.Get(&dishData, sqlQuery, dishId)
	if getErr != nil {
		return dishData, getErr
	}

	return dishData, nil
}

func IsDishExists(name string, restaurantId string) (bool, error) {
	sqlQuery := `select count(Id) > 0 as isExists from public.dishes where name=$1 and restaurantId=$2 and archivedAt is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, name, restaurantId)
	return exists, err
}

func CreateDishHelper(name string, price int, restaurantId string) (string, error) {
	var dishId uuid.UUID
	sqlQuery := `insert into public.dishes (Id, name, price, restaurantId, createdAt) 
					values ($1, $2, $3, $4, $5) returning Id;`

	crtErr := database.RmsDB.Get(&dishId, sqlQuery, uuid.New(), name, price, restaurantId, time.Now())
	return dishId.String(), crtErr
}

func GetAllDishHelper() ([]models.Dish, error) {
	sqlQuery := `select d.Id, d.name, d.price, d.restaurantId, r.name as restaurantName, d.createdAt
					from public.dishes d inner join public.restaurants r
					on r.Id = d.restaurantId 
					where d.archivedAt is null`
	dishData := make([]models.Dish, 0)
	err := database.RmsDB.Select(&dishData, sqlQuery)

	return dishData, err
}

func GetAllDishSubAdminHelper(createdBy string) ([]models.Dish, error) {
	sqlQuery := `select d.Id, d.name, d.price, d.restaurantId, r.name as restaurantName, d.createdAt
					from public.dishes d inner join public.restaurants r
					on r.Id = d.restaurantId 
					where d.archivedAt is null and r.createdBy = $1`
	dishData := make([]models.Dish, 0)
	err := database.RmsDB.Select(&dishData, sqlQuery, createdBy)

	return dishData, err
}

func GetAllDishesByUserHelper() ([]models.Dish, error) {
	sqlQuery := `select d.Id, d.name, d.price, d.restaurantId, r.name as restaurantName, d.createdAt
				from public.dishes d inner join public.restaurants r
				on r.Id = d.restaurantId 
				where d.archivedAt is null`
	dishData := make([]models.Dish, 0)
	err := database.RmsDB.Select(&dishData, sqlQuery)

	return dishData, err
}

func GetDishesByRestIdHelper(restaurantId string) ([]models.DishData, error) {
	sqlQuery := `select Id, name, price, restaurantId, createdAt from public.dishes where restaurantId=$1`

	dishData := make([]models.DishData, 0)
	err := database.RmsDB.Select(&dishData, sqlQuery, restaurantId)

	return dishData, err
}
