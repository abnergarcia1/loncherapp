package api

import (
	"context"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/reviews"
)

type ReviewsService struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewReviewsService(context context.Context) *ReviewsService {
	return &ReviewsService{
		ctx: context,
		db:  sql.NewClient(),
	}
}

func (r *ReviewsService) CreateReview(review *models.Review) (*models.Review, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		r.db.Disconnect()
	}()

	res, err := ses.Execute("INSERT INTO Reviews(Lonchera_ID, Comment, User_ID, User_Name, Rating, Created_At) VALUES(?,?,?,?,?,CURRENT_TIMESTAMP) ",
		review.LoncheraID, review.Comment, review.UserID, review.UserName, review.Rating)
	if err != nil {
		log.WithFields(log.Fields{
			"reviewRequest": review,
		}).Errorf("Error when trying to create review in DB: %v", err)
		return nil, err
	}

	objectID, err := res.LastInsertId()
	if err == nil {
		review.ID = int32(objectID)
	}

	return review, nil
}

func (r *ReviewsService) GetReviewsByProfileID(profileID int32) (*[]models.Review, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		r.db.Disconnect()
	}()

	var listRoutes []models.Review

	err := ses.Select(&listRoutes, `SELECT * FROM Reviews WHERE Lonchera_ID=?`, profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
		}).Errorf("Error when trying to get reviews by LoncheraID  in DB: %v", err)
		return nil, err
	}

	return &listRoutes, nil
}

func (r *ReviewsService) GetAverageRatingByProfileID(profileID int32) (*models.RatingAverage, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		r.db.Disconnect()
	}()

	var rating models.RatingAverage

	err := ses.QueryOne(&rating, `SELECT ? AS Lonchera_ID , AVG(r.Rating) as Rating FROM Reviews r WHERE Lonchera_ID=?  `, profileID, profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
		}).Errorf("Error when trying to get average rating in DB: %v", err)
		return nil, err
	}

	return &rating, nil
}

func (r *ReviewsService) DeleteReview(reviewID int32, profileID int32) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		r.db.Disconnect()
	}()

	result, err := ses.Execute(`DELETE FROM Reviews WHERE ID=? AND Lonchera_ID=? `, reviewID, profileID)

	if err != nil {
		log.WithFields(log.Fields{
			"reviewID":  reviewID,
			"profileID": profileID,
		}).Errorf("Error when trying to delete Review in DB: %v", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"reviewID":  reviewID,
			"profileID": profileID,
		}).Errorf("Error when trying to delete menu in DB: %v", sql.DbErrNoDocuments)
		return sql.DbErrNoDocuments
	}

	return nil
}
