package dbHelper

import (
	"time"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/models"
	"github.com/google/uuid"
)

func CreateSubAdminHelper(email, name, hashpwd, createdby, role string) (string, error) {
	userId := uuid.New()
	sqlQuery := `insert into public.users (userid, name, email, password, role, createdby, createdat, updatedby, updatedat) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning userid;`

	crtErr := database.RmsDB.Get(&userId, sqlQuery, userId, name, email, hashpwd, role, createdby, time.Now(), createdby, time.Now())
	return userId.String(), crtErr
}

func GetSubAdminsHelper() ([]models.User, error) {
	sqlquery := `select userid, name, email, role, createdby, createdat, updatedby, updatedat, archivedat 
					from public.users where role='sub-admin' and archivedat is null`

	subadmindata := make([]models.User, 0)
	err := database.RmsDB.Select(&subadmindata, sqlquery)
	return subadmindata, err
}
