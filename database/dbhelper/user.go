package dbHelper

import (
	"strings"
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func IsUserExists(email string) (bool, error) {
	sqlQuery := `select count(Id) > 0 as isExists from public.users where email=$1 and archivedAt is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, email)
	return exists, err
}

func GetUserInfoForLogin(payload models.LoginRequest) (string, string, string, error) {
	sqlQuery := `select Id, password, role from public.users where email=$1 and archivedAt is NULL`

	var loginResp models.LoginData
	getErr := database.RmsDB.Get(&loginResp, sqlQuery, payload.Email)
	if getErr != nil {
		return "", "", "", getErr
	}

	return loginResp.UserID, loginResp.PasswordHash, string(loginResp.Role), nil
}

func CreateUserSession(userId string) (string, error) {
	sqlQuery := `insert into public.user_session (Id, userId)
				 values ($1, $2);`
	sessionId := uuid.New()

	_, crtErr := database.RmsDB.Exec(sqlQuery, sessionId, userId)
	return sessionId.String(), crtErr
}

func FetchUserDataBySessionId(sessionId string) (models.SessionData, error) {
	var userData models.SessionData

	sqlQuery := `select u.email, s.archivedAt
				from public.user_session s
						inner join public.users u on s.userId = u.Id
				where s.Id = $1`

	getErr := database.RmsDB.Get(&userData, sqlQuery, sessionId)
	return userData, getErr
}

func DeleteUserSession(sessionId string) error {
	sqlQuery := `UPDATE public.user_session set archivedAt=NOW() 
					where id=$1 and archivedAt is null`

	_, err := database.RmsDB.Exec(sqlQuery, sessionId)
	return err
}

func CreateUserHelper(tx *sqlx.Tx, email, name, hashPwd, createdBy, role string, address []models.AddressData) (string, error) {
	var userId, addressId uuid.UUID
	sqlQuery := `insert into public.users (Id, name, email, password, role, createdBy, updatedBy) 
					values ($1, $2, $3, $4, $5, $6, $7) returning Id;`

	userId = uuid.New()
	crtErr := tx.Get(&userId, sqlQuery, uuid.New(), name, email, hashPwd, role, createdBy, createdBy)

	if crtErr == nil && strings.EqualFold(role, string(models.RoleUser)) {
		sqlQuery = `insert into public.addresses (Id, address, latitude, longitude, userId) 
						values ($1, $2, $3, $4, $5) returning Id;`

		for i := range address {
			crtErr = tx.Get(&addressId, sqlQuery, uuid.New(), address[i].Address, address[i].Latitude, address[i].Longitude, userId)
		}
	}

	return userId.String(), crtErr
}

func CreateSignUpHelper(email, name, hashPwd, role string, address []models.AddressData) (string, error) {
	var userId, addressId uuid.UUID
	sqlQuery := `insert into public.users (Id, name, email, password, role, createdBy, createdAt, updatedBy, updatedAt) 
					values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning Id;`

	userId = uuid.New()
	crtErr := database.RmsDB.Get(&userId, sqlQuery, userId, name, email, hashPwd, role, userId, time.Now(), userId, time.Now())

	if crtErr == nil && strings.EqualFold(role, string(models.RoleUser)) {
		sqlQuery = `insert into public.addresses (Id, address, latitude, longitude, userId, createdAt) 
						values ($1, $2, $3, $4, $5, $6) returning Id;`

		for i := range address {
			crtErr = database.RmsDB.Get(&addressId, sqlQuery, uuid.New(), address[i].Address, address[i].Latitude, address[i].Longitude, userId, time.Now())
		}
	}

	return userId.String(), crtErr
}

func GetUserDataHelper(tx *sqlx.Tx, role string) ([]models.User, error) {
	sqlQuery := `select Id, name, email, role, createdBy, createdAt, updatedBy, updatedAt 
					from public.users where role=$1 and archivedAt is null`
	userData := make([]models.User, 0)
	err := tx.Select(&userData, sqlQuery, role)

	return userData, err
}

func GetAddressForUserHelper(tx *sqlx.Tx, userData []models.User) ([]models.User, error) {
	sqlQuery := `select Id, address, latitude, longitude, userId, createdAt 
							from public.addresses where archivedAt is NULL`

	var err error
	addressData := make([]models.AddressData, 0)
	err = tx.Select(&addressData, sqlQuery)

	addressMap := make(map[string][]models.AddressData)
	for _, addr := range addressData {
		addressMap[*addr.UserId] = append(addressMap[*addr.UserId], addr)
	}
	for i := range userData {
		userAddress := addressMap[userData[i].ID]
		userData[i].Address = userAddress
	}

	return userData, err
}

func GetUsersSubAdminHelper(tx *sqlx.Tx, role, createdBy string) ([]models.User, error) {
	sqlQuery := `select Id, name, email, role, createdBy, createdAt, updatedBy, updatedAt 
					from public.users where role=$1 and createdBy=$2 and archivedAt is null`
	userData := make([]models.User, 0)
	err := tx.Select(&userData, sqlQuery, role, createdBy)

	return userData, err
}

func GetUserLatitudeAndLongitude(tx *sqlx.Tx, userId string) ([]models.Coordinates, error) {
	sqlQuery := `select a.latitude, a.longitude, a.address
				from public.addresses a
						inner join public.users u on u.id = a.userid
				where u.Id = $1`
	coordinates := make([]models.Coordinates, 0)
	err := tx.Select(&coordinates, sqlQuery, userId)
	return coordinates, err
}
