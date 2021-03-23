package services

import (
	"context"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"
)

type ImagesService struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewImagesService(context context.Context) *ImagesService {
	return &ImagesService{
		ctx: context,
		db:  sql.NewClient(),
	}
}

func (i *ImagesService) UpdateProfileImageURL(profileID int32, imageURL string) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		i.db.Disconnect()
	}()

	_, err := ses.Execute(`UPDATE Loncheras SET Cover_Image_URL=? WHERE ID=?`,
		imageURL, profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
			"imageURL":  imageURL,
		}).Errorf("Error when trying to update profile cover URL in DB: %v", err)
		return err
	}

	return nil
}

func (i *ImagesService) UpdateMenuImageURL(profileID int32, menuID int32, imageURL string) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		i.db.Disconnect()
	}()

	result, err := ses.Execute("UPDATE Menus SET Image_URL=? WHERE Lonchera_ID=? AND ID=?",
		imageURL, profileID, menuID)
	if err != nil {
		log.WithFields(log.Fields{
			"menuID":    menuID,
			"imageURL":  imageURL,
			"profileID": profileID,
		}).Errorf("Error when trying to update menu in DB: %v", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"menuID":    menuID,
			"imageURL":  imageURL,
			"profileID": profileID,
		}).Errorf("Error when trying to update menu in DB: %v", sql.DbErrNoDocuments)
		return sql.DbErrNoDocuments
	}

	return nil
}
