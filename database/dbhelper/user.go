package dbHelper

import (
	"strings"
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func GetUserById(userId, role string) (models.User, error) {
	sqlQuery := `select userid, name, email, role, createdby, createdat, updatedby, updatedat, archivedat from public.users 
					where userid=$1 and archivedat is NULL`

	var userdata models.User
	getErr := database.RmsDB.Get(&userdata, sqlQuery, userId)
	if getErr != nil {
		return userdata, getErr
	}

	return userdata, nil
}

func IsUserExists(email string) (bool, error) {
	sqlQuery := `select count(userid) > 0 as isExists from public.users where email=$1 and archivedat is null`
	var exists bool
	err := database.RmsDB.Get(&exists, sqlQuery, email)
	return exists, err
}

func GetUserInfo(payload models.LoginRequest) (string, string, error) {
	sqlQuery := `select userid, password, role from public.users where email=$1 and archivedat is NULL`

	var loginResp models.LoginData
	getErr := database.RmsDB.Get(&loginResp, sqlQuery, payload.Email)
	if getErr != nil {
		return "", "", getErr
	}

	// passwordErr := utils.CheckPassword(payload.Password, loginResp.PasswordHash)
	// if passwordErr != nil {
	// 	return "", "", passwordErr
	// }

	return loginResp.UserID, loginResp.Role, nil
}

func CreateUserSession(userId string) (string, error) {
	sqlQuery := `insert into public.usersession values ($1, $2);`
	sessionId := uuid.New()

	_, crtErr := database.RmsDB.Exec(sqlQuery, sessionId, userId)
	return sessionId.String(), crtErr
}

func GetArchivedAt(sessionId string) (*time.Time, error) {
	var archivedAt *time.Time

	SQL := `SELECT archivedat FROM public.usersession WHERE sessionid = $1`

	getErr := database.RmsDB.Get(&archivedAt, SQL, sessionId)
	return archivedAt, getErr
}

func DeleteUserSession(sessionId string) error {
	sqlQuery := `UPDATE public.usersession set archivedat=NOW() 
					where sessionid=$1 and archivedat is null`

	//loc, _ := time.LoadLocation("Asia/Shanghai")
	//archivedat := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)

	_, err := database.RmsDB.Exec(sqlQuery, sessionId)
	return err
}

// Add User and address
func CreateUserHelper(email, name, hashpwd, createdby, role string, address []models.AddressData) (string, error) {
	var userId, addressId uuid.UUID
	sqlQuery := `insert into public.users (userid, name, email, password, role, createdby, createdat, updatedby, updatedat) 
					values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning userid;`

	crtErr := database.RmsDB.Get(&userId, sqlQuery, uuid.New(), name, email, hashpwd, role, createdby, time.Now(), createdby, time.Now())

	if crtErr == nil && strings.EqualFold(role, "user") {
		sqlQuery = `insert into public.address (addressid, addressline, latitude, longitude, user_id, createdat) 
						values ($1, $2, $3, $4, $5, $6) returning addressid;`

		for i := range address {
			crtErr = database.RmsDB.Get(&addressId, sqlQuery, uuid.New(), address[i].AddressLine, address[i].Latitude, address[i].Longitude, userId, time.Now())
		}
	}

	return userId.String(), crtErr
}

// Fetch Users
func GetUsersHelper(role string) ([]models.User, error) {
	sqlquery := `select userid, name, email, role, createdby, createdat, updatedby, updatedat, archivedat 
					from public.users where role=$1 and archivedat is null`
	userdata := make([]models.User, 0)
	err := database.RmsDB.Select(&userdata, sqlquery, role)

	return userdata, err
}

// Fetch Address
func GetAddressForUser(userdata []models.User) ([]models.User, error) {
	sqlqueryaddress := `select addressid, addressline, latitude, longitude, user_id, createdat, archivedat 
							from public.address where user_id=$1`

	var err error
	for i := range userdata {
		addressdata := make([]models.AddressData, 0)
		err = database.RmsDB.Select(&addressdata, sqlqueryaddress, userdata[i].UserID)
		userdata[i].Address = addressdata
	}

	return userdata, err
}

func GetUsersSubAdminHelper(role, createdBy string) ([]models.User, error) {
	sqlquery := `select userid, name, email, role, createdby, createdat, updatedby, updatedat, archivedat 
					from public.users where role=$1 and createdby=$2 and archivedat is null`
	userdata := make([]models.User, 0)
	err := database.RmsDB.Select(&userdata, sqlquery, role, createdBy)

	return userdata, err
}
