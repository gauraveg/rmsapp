package dbHelper

import (
	"strings"
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func GetUserById(userId, role string) (models.User, error) {
	sqlQuery := `select Id, name, email, role, createdBy, createdAt, updatedBy, updatedAt from public.users 
					where Id=$1 and archivedAt is NULL`

	var userData models.User
	getErr := database.RmsDB.Get(&userData, sqlQuery, userId)
	if getErr != nil {
		return userData, getErr
	}

	if role == "user" {
		sqlQueryAddr := `select Id, address, latitude, longitude, userId, createdAt 
							from public.addresses where userId=$1`

		addressData := make([]models.AddressData, 0)
		getErr = database.RmsDB.Select(&addressData, sqlQueryAddr, userData.ID)
		if getErr != nil {
			return userData, getErr
		}
		userData.Address = addressData

	}

	return userData, nil
}

func IsUserExists(email string) (bool, error) {
	sqlQuery := `select count(Id) > 0 as isExists from public.users where email=$1 and archivedAt is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, email)
	return exists, err
}

func GetUserInfo(payload models.LoginRequest) (string, string, string, error) {
	sqlQuery := `select Id, password, role from public.users where email=$1 and archivedAt is NULL`

	var loginResp models.LoginData
	getErr := database.RmsDB.Get(&loginResp, sqlQuery, payload.Email)
	if getErr != nil {
		return "", "", "", getErr
	}

	return loginResp.UserID, loginResp.PasswordHash, loginResp.Role, nil
}

func CreateUserSession(userId string) (string, error) {
	sqlQuery := `insert into public.user_session (Id, userId) values ($1, $2);`
	sessionId := uuid.New()

	_, crtErr := database.RmsDB.Exec(sqlQuery, sessionId, userId)
	return sessionId.String(), crtErr
}

func FetchUserDetails(sessionId string) (models.SessionData, error) {
	var userData models.SessionData

	SQL := `select u.email, s.archivedAt from public.user_session s inner join public.users u on s.userId = u.Id
			where s.Id = $1`

	getErr := database.RmsDB.Get(&userData, SQL, sessionId)
	return userData, getErr
}

func DeleteUserSession(sessionId string) error {
	sqlQuery := `UPDATE public.user_session set archivedAt=NOW() 
					where id=$1 and archivedAt is null`

	_, err := database.RmsDB.Exec(sqlQuery, sessionId)
	return err
}

// Add User and address
func CreateUserHelper(email, name, hashPwd, createdBy, role string, address []models.AddressData) (string, error) {
	var userId, addressId uuid.UUID
	sqlQuery := `insert into public.users (Id, name, email, password, role, createdBy, createdAt, updatedBy, updatedAt) 
					values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning Id;`

	crtErr := database.RmsDB.Get(&userId, sqlQuery, uuid.New(), name, email, hashPwd, role, createdBy, time.Now(), createdBy, time.Now())

	if crtErr == nil && strings.EqualFold(role, "user") {
		sqlQuery = `insert into public.addresses (Id, address, latitude, longitude, userId, createdAt) 
						values ($1, $2, $3, $4, $5, $6) returning Id;`

		for i := range address {
			crtErr = database.RmsDB.Get(&addressId, sqlQuery, uuid.New(), address[i].Address, address[i].Latitude, address[i].Longitude, userId, time.Now())
		}
	}

	return userId.String(), crtErr
}

// Fetch Users
func GetUsersHelper(role string) ([]models.User, error) {
	sqlQuery := `select Id, name, email, role, createdBy, createdAt, updatedBy, updatedAt 
					from public.users where role=$1 and archivedAt is null`
	userData := make([]models.User, 0)
	err := database.RmsDB.Select(&userData, sqlQuery, role)

	return userData, err
}

// Fetch Address
func GetAddressForUser(userData []models.User) ([]models.User, error) {
	sqlQuery := `select Id, address, latitude, longitude, userId, createdAt 
							from public.addresses where userId=$1`

	var err error
	for i := range userData {
		addressData := make([]models.AddressData, 0)
		err = database.RmsDB.Select(&addressData, sqlQuery, userData[i].ID)
		userData[i].Address = addressData
	}

	return userData, err
}

func GetUsersSubAdminHelper(role, createdBy string) ([]models.User, error) {
	sqlQuery := `select Id, name, email, role, createdBy, createdAt, updatedBy, updatedAt 
					from public.users where role=$1 and createdBy=$2 and archivedAt is null`
	userData := make([]models.User, 0)
	err := database.RmsDB.Select(&userData, sqlQuery, role, createdBy)

	return userData, err
}
