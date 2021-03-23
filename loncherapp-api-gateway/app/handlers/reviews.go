package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	sharedModels "bitbucket.org/edgelabsolutions/loncherapp-core/models/shared"

	"bitbucket.org/edgelabsolutions/loncherapp-core/tools"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/reviews"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"
	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/reviews"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
)

var (
	reviewsServiceHost = os.Getenv("REVIEWS_SERVICE_HOST")
)

// RoutesHandler struct
type reviewsHandler struct {
	rd    auth.AuthInterface
	tk    auth.TokenInterface
	tools tools.Tools
}

func NewReviewsHandler(rd auth.AuthInterface, tk auth.TokenInterface) *reviewsHandler {
	return &reviewsHandler{rd, tk, tools.Tools{}}
}

func (r *reviewsHandler) GetReviewsByProfileID(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(reviewsServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewReviewsServiceClient(conn)

	response, err := s.GetReviewsByProfileID(context.Background(), &pbm.SimpleRequestByID{Id: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting reviews items: %v", err))
		return
	}

	var listReviews = make([]models.Review, len(response.Reviews))
	for i, review := range response.Reviews {
		listReviews[i] = *ConvertReviewPMtoModel(review)
	}

	c.JSON(http.StatusOK, listReviews)
}

func (r *reviewsHandler) GetAverageRatingByProfileID(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(reviewsServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewReviewsServiceClient(conn)

	response, err := s.GetAverageRatingByProfileID(context.Background(), &pbm.SimpleRequestByID{Id: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting review average: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *reviewsHandler) CreateReview(c *gin.Context) {
	review := &models.Review{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviewsRequest := ConvertReviewModelToPM(review)
	reviewsRequest.LoncheraId = int32(profileID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(reviewsServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewReviewsServiceClient(conn)

	response, err := s.CreateReview(context.Background(), reviewsRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create review item: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *reviewsHandler) DeleteReview(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	reviewID, _ := strconv.Atoi(c.Param("reviewID"))

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(reviewsServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewReviewsServiceClient(conn)

	res, err := s.DeleteReview(context.Background(), &pbm.ReviewRequest{Id: int32(reviewID), ProfileId: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when delete review item %v : %v", reviewID, err))
		return
	}

	c.JSON(http.StatusNotFound, sharedModels.SimpleResponse{Success: res.Success, Message: res.Message})
}

// ConvertReviewPMtoModel Convert Review PBM to Model object
func ConvertReviewPMtoModel(pbmReview *pbm.Review) *models.Review {
	var createdAt time.Time
	if pbmReview.CreatedAt != "" {
		createdAt, _ = time.Parse("2006-01-02 15:04:05", pbmReview.CreatedAt)
	}

	review := &models.Review{
		ID:         pbmReview.Id,
		LoncheraID: pbmReview.LoncheraId,
		Comment:    pbmReview.Comment,
		UserID:     pbmReview.UserId,
		UserName:   pbmReview.UserName,
		Rating:     pbmReview.Rating,
		CreatedAt:  createdAt,
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
	}

	return reviewModel
}
