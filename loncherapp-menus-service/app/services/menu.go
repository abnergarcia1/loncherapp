package services

import (
	"context"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/menus"
	sharedModels "bitbucket.org/edgelabsolutions/loncherapp-core/models/shared"
)

type MenusService struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewMenuService(context context.Context) *MenusService {
	return &MenusService{
		ctx: context,
		db:  sql.NewClient(),
	}
}

func (m *MenusService) CreateMenu(menu *models.Menu) (*models.Menu, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		m.db.Disconnect()
	}()

	res, err := ses.Execute("INSERT INTO Menus(Lonchera_ID, Name, Description, Price, Currency, Created_At, Updated_At) VALUES(?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) ",
		menu.LoncheraID, menu.Name, menu.Description, menu.Price, menu.Currency)
	if err != nil {
		log.WithFields(log.Fields{
			"menuRequest": menu,
		}).Errorf("Error when trying to create menu item in DB: %v", err)
		return nil, err
	}

	objectID, err := res.LastInsertId()
	if err == nil {
		menu.ID = int32(objectID)
	}

	return menu, nil
}

func (m *MenusService) UpdateMenu(menu *models.Menu) (*sharedModels.SimpleResponse, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		m.db.Disconnect()
	}()

	result, err := ses.Execute("UPDATE Menus SET Name=?, Description=?, Price=?, Currency=?, Updated_At=CURRENT_TIMESTAMP WHERE ID=?",
		menu.Name, menu.Description, menu.Price, menu.Currency, menu.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"menuRequest": menu,
		}).Errorf("Error when trying to update menu in DB: %v", err)
		return nil, err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"menuRequest": menu,
		}).Errorf("Error when trying to update menu in DB: %v", sql.DbErrNoDocuments)
		return nil, sql.DbErrNoDocuments
	}

	return &sharedModels.SimpleResponse{Success: true, Message: "Updated menu item correctly"}, nil
}

func (m *MenusService) GetMenuByProfileID(profileID int32) (*[]models.Menu, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		m.db.Disconnect()
	}()

	var listRoutes []models.Menu

	err := ses.Select(&listRoutes, `SELECT * FROM Menus WHERE Lonchera_ID=?`, profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
		}).Errorf("Error when trying to get menus by LoncheraID  in DB: %v", err)
		return nil, err
	}

	return &listRoutes, nil
}

func (m *MenusService) GetMenuByID(menuID int32) (*models.Menu, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		m.db.Disconnect()
	}()

	var route models.Menu

	err := ses.QueryOne(&route, `SELECT * FROM Menus WHERE ID=?`, menuID)
	if err != nil {
		log.WithFields(log.Fields{
			"menuID": menuID,
		}).Errorf("Error when trying to get menu by ID in DB: %v", err)
		return nil, err
	}

	return &route, nil
}

func (m *MenusService) DeleteMenu(menuID int32, profileID int32) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		m.db.Disconnect()
	}()

	result, err := ses.Execute(`DELETE FROM Menus WHERE ID=? AND Lonchera_ID=? `, menuID, profileID)

	if err != nil {
		log.WithFields(log.Fields{
			"menuID":    menuID,
			"profileID": profileID,
		}).Errorf("Error when trying to delete Menu in DB: %v", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"menuID":    menuID,
			"profileID": profileID,
		}).Errorf("Error when trying to delete menu in DB: %v", sql.DbErrNoDocuments)
		return sql.DbErrNoDocuments
	}

	return nil
}
