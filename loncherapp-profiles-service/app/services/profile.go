package services

import (
	"context"
	"errors"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/profiles"
	sharedModels "bitbucket.org/edgelabsolutions/loncherapp-core/models/shared"
)

type ProfileService struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewProfileService(context context.Context) *ProfileService {
	return &ProfileService{
		ctx: context,
		db:  sql.NewClient(),
	}
}

func (p *ProfileService) CreateProfile(profile *models.Profile) (*models.Profile, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	_, err := ses.Execute(`INSERT INTO Loncheras(User_ID, Description, Category_ID, Cover_Image_URL, Website, Active, Membership_Due_Date, Created_At, Updated_At) VALUE(?,?,?,?,?,0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) `,
		profile.UserID, profile.Description, profile.CategoryID, profile.CoverImageURL, profile.Website)
	if err != nil {
		log.WithFields(log.Fields{
			"profile": profile,
		}).Errorf("Error when trying to create profile in DB: %v", err)
		return nil, err
	}

	return profile, nil
}

func (p *ProfileService) UpdateProfile(profile *models.Profile) (*sharedModels.SimpleResponse, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	_, err := ses.Execute(`UPDATE Loncheras SET Description=?, Category_ID=?, Website=?, Updated_At=CURRENT_TIMESTAMP WHERE User_ID=?`,
		profile.Description, profile.CategoryID, profile.Website, profile.UserID)
	if err != nil {
		log.WithFields(log.Fields{
			"profile": profile,
		}).Errorf("Error when trying to update profile in DB: %v", err)
		return nil, err
	}

	return &sharedModels.SimpleResponse{Success: true, Message: "Updated profile correctly"}, nil
}

func (p *ProfileService) GetProfileByUserID(userID int32) (*models.Profile, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	profile := &models.Profile{}
	err := ses.QueryOne(profile, `SELECT * FROM Loncheras WHERE User_ID=? `, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"userid": userID,
		}).Errorf("Error when trying to get profile from DB: %v", err)
		return nil, err
	}

	return profile, nil
}

func (p *ProfileService) GetProfileByID(profileID int32) (*models.Profile, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	profile := &models.Profile{}
	err := ses.QueryOne(profile, `SELECT * FROM Loncheras WHERE ID=? `, profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
		}).Errorf("Error when trying to get profile from DB: %v", err)
		return nil, err
	}

	return profile, nil
}

func (p *ProfileService) GetProfilesByCategory(categoryID int32) (*[]models.Profile, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	profile := []models.Profile{}
	err := ses.Select(&profile, `SELECT ID, Description, Category_ID, Cover_Image_URL, Website, Created_At FROM Loncheras WHERE Active=1 AND Category_ID=? `, categoryID)
	if err != nil {
		log.WithFields(log.Fields{
			"categoryID": categoryID,
		}).Errorf("Error when trying to get profile from DB: %v", err)
		return nil, err
	}

	return &profile, nil
}

func (p *ProfileService) GetProfilesByUserLocation(ctx context.Context, profileQueryUserLocation *pbm.ProfileQueryUserLocation) (*pbm.Profiles, error) {

	return nil, errors.New("implement me")
}

func (p *ProfileService) GetProfilesByUserFavorites(userID int32) (*[]models.Profile, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		ses.Disconnect()
	}()

	profile := []models.Profile{}
	err := ses.Select(&profile, `SELECT lon.ID, lon.Description, lon.Category_ID, lon.Cover_Image_URL, lon.Website, lon.Created_At 
	FROM Loncheras lon
	INNER JOIN Users_Favorites uf
	ON(lon.ID=uf.Lonchera_ID)
	WHERE lon.Active=1 AND uf.User_ID =?`, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"userID": userID,
		}).Errorf("Error when trying to get profile from DB: %v", err)
		return nil, err
	}

	return &profile, nil
}

func (p *ProfileService) DeleteProfileUserFavorite(userID int32, profileID int32) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		p.db.Disconnect()
	}()

	_, err := ses.Execute(`DELETE FROM Users_Favorites WHERE User_ID=? AND Lonchera_ID=? `,
		userID, profileID)

	if err != nil {
		log.WithFields(log.Fields{
			"userID":    userID,
			"profileID": profileID,
		}).Errorf("Error when trying to create profile in DB: %v", err)
		return err
	}

	return nil
}

func (p *ProfileService) SetProfileUserFavorite(userID int32, profileID int32) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		p.db.Disconnect()
	}()

	_, err := ses.Execute(`INSERT INTO Users_Favorites(User_ID, Lonchera_ID, Created_At) VALUE(?,?,CURRENT_TIMESTAMP) `,
		userID, profileID)

	if err != nil {
		log.WithFields(log.Fields{
			"userID":    userID,
			"profileID": profileID,
		}).Errorf("Error when trying to create profile in DB: %v", err)
		return err
	}

	return nil
}
