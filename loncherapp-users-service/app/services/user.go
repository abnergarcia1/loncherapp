package services

import (
	"context"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/users"
)

type UserService struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewUserService(context context.Context) *UserService {
	return &UserService{
		ctx: context,
	}
}

func (u *UserService) CreateUser(user *models.User) (*models.User, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	_, err := ses.Execute(`INSERT INTO Users(Type_ID, FirstName, LastName, Email, CreationDate, UpdatedDate, Active, Password) VALUE(?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1,?) `, user.TypeID, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		user.Password = ""
		log.WithFields(log.Fields{
			"user": user,
		}).Errorf("Error when trying to create user in DB: %v", err)
		return nil, err
	}

	return user, nil
}

func (u *UserService) GetUser(userid int) (user *models.User, err error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	user = &models.User{}
	err = ses.QueryOne(user, `SELECT * FROM Users WHERE ID=? `, userid)
	if err != nil {
		log.WithFields(log.Fields{
			"userID": userid,
		}).Errorf("Error when trying to get user from DB: %v", err)
	}

	return
}

func (u *UserService) GetUserByEmail(email string) (user *models.User, err error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	user = &models.User{}
	err = ses.QueryOne(user, `SELECT * FROM Users WHERE Email=? `, email)
	if err != nil {
		log.WithFields(log.Fields{
			"Email": email,
		}).Errorf("Error when trying to get user email from DB: %v", err)
	}

	return
}

func (u *UserService) GetUserByEmailPassword(email string, password string) (user *models.User, err error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	user = &models.User{}
	err = ses.QueryOne(user, `SELECT Users.ID, Users.Type_ID, Users.FirstName, Users.LastName, Users.Email, Users.CreationDate, Users.UpdatedDate, IFNULL(Loncheras.ID,0) Profile_ID, IFNULL(Loncheras.Active,0) Profile_Active FROM Users LEFT JOIN Loncheras ON(Users.ID=Loncheras.User_ID) WHERE Users.Active=1 AND Users.Email=? AND Users.Password=?`, email, password)
	if err != nil {
		log.WithFields(log.Fields{
			"Email": email,
		}).Errorf("Error when trying to get user from DB using email and password: %v", err)
	}

	return
}
