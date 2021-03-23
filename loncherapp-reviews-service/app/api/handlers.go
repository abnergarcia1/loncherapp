package api

import (
	"context"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/reviews"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"
)

//ReviewsAPIServer Reviews API Server
type ReviewsAPIServer struct {
	ReviewsService *ReviewsService
}

//GetMenuByProfileID Get Menu object list from DB using ProfileID
func (r ReviewsAPIServer) GetReviewsByProfileID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.Reviews, error) {
	reviews, err := r.ReviewsService.GetReviewsByProfileID(id.Id)
	if err != nil {
		return nil, err
	}

	var listReviews = make([]*pbm.Review, len(*reviews))
	for i, review := range *reviews {
		listReviews[i] = ConvertReviewModelToPM(&review)
	}

	return &pbm.Reviews{Reviews: listReviews}, nil
}

//GetAverageRatingByProfileID Get Average Rating for Profile from DB using ProfileID
func (r ReviewsAPIServer) GetAverageRatingByProfileID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.AverageRatingResponse, error) {
	avgRating, err := r.ReviewsService.GetAverageRatingByProfileID(id.Id)
	if err != nil {
		return nil, err
	}

	return &pbm.AverageRatingResponse{AverageReview: avgRating.Rating, LoncheraId: avgRating.LoncheraID}, nil
}

//CreateReview Create new Review object in DB
func (r ReviewsAPIServer) CreateReview(ctx context.Context, review *pbm.Review) (*pbm.Review, error) {
	modelReview := ConvertReviewPMtoModel(review)
	reviewResponse, err := r.ReviewsService.CreateReview(modelReview)
	if err != nil {
		return nil, err
	}

	return ConvertReviewModelToPM(reviewResponse), nil
}

//DeleteReview Delete Review  Object from DB using ReviewID
func (r ReviewsAPIServer) DeleteReview(ctx context.Context, request *pbm.ReviewRequest) (*pbm.SimpleResponse, error) {
	err := r.ReviewsService.DeleteReview(request.Id, request.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "review item correctly deleted"}, nil
}

func NewReviewsAPIServer(ctx context.Context) *ReviewsAPIServer {
	return &ReviewsAPIServer{
		ReviewsService: NewReviewsService(ctx),
	}
}

// ConvertReviewPMtoModel Convert Review PBM to Model object
func ConvertReviewPMtoModel(pbmReview *pbm.Review) *models.Review {
	review := &models.Review{
		ID:         pbmReview.Id,
		LoncheraID: pbmReview.LoncheraId,
		Comment:    pbmReview.Comment,
		UserID:     pbmReview.UserId,
		UserName:   pbmReview.UserName,
		Rating:     pbmReview.Rating,
	}

	return review
}

// ConvertReviewModelToPM Convert Review Object to PBM
func ConvertReviewModelToPM(review *models.Review) *pbm.Review {
	reviewModel := &pbm.Review{
		Id:         review.ID,
		LoncheraId: review.LoncheraID,
		Comment:    review.Comment,
		UserId:     review.UserID,
		UserName:   review.UserName,
		Rating:     review.Rating,
		CreatedAt:  review.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return reviewModel
}
