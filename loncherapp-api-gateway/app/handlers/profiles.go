package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"bitbucket.org/edgelabsolutions/loncherapp-core/tools"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/profiles"

	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/profiles"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
)

var (
	profilesServiceHost = os.Getenv("PROFILES_SERVICE_HOST")
)

// UserHandler struct
type profilesHandler struct {
	rd    auth.AuthInterface
	tk    auth.TokenInterface
	tools tools.Tools
}

func NewProfilesHandler(rd auth.AuthInterface, tk auth.TokenInterface) *profilesHandler {
	return &profilesHandler{rd, tk, tools.Tools{}}
}

func (p *profilesHandler) CreateProfile(c *gin.Context) {
	profile := &models.Profile{}

	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !p.tools.ValidateDataIDToken(c, "user_id", int(profile.UserID), p.rd, p.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	profileRequest := ConvertProfileModeltoPM(profile)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	response, err := s.CreateProfile(context.Background(), profileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (p *profilesHandler) UpdateProfile(c *gin.Context) {
	profile := &models.Profile{}

	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !p.tools.ValidateDataIDToken(c, "user_id", int(profile.UserID), p.rd, p.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	profileRequest := ConvertProfileModeltoPM(profile)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	response, err := s.UpdateProfile(context.Background(), profileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create profile: %v", err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (p *profilesHandler) GetProfileByUserID(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userID"))

	if !p.tools.ValidateDataIDToken(c, "user_id", int(userID), p.rd, p.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	response, err := s.GetProfileByUserID(context.Background(), &pbm.SimpleRequestByID{Id: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when getting profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (p *profilesHandler) GetProfileByID(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	response, err := s.GetProfileByID(context.Background(), &pbm.SimpleRequestByID{Id: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when getting profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (p *profilesHandler) GetProfilesByCategories(c *gin.Context) {
	categoryID, _ := strconv.Atoi(c.Param("categoryID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	response, err := s.GetProfilesByCategory(context.Background(), &pbm.ProfileQueryCategory{Category: int32(categoryID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when getting profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (p *profilesHandler) GetProfilesByUserLocation(c *gin.Context) {

	return
}

func (p *profilesHandler) GetProfilesByUserFavorites(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userID"))

	if !p.tools.ValidateDataIDToken(c, "user_id", userID, p.rd, p.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	response, err := s.GetProfilesByUserFavorites(context.Background(), &pbm.SimpleRequestByID{Id: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when getting profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (p *profilesHandler) SetUserProfileFavorite(c *gin.Context) {
	favorite := &models.Favorite{}
	userID, _ := strconv.Atoi(c.Param("userID"))

	if err := c.ShouldBindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !p.tools.ValidateDataIDToken(c, "user_id", userID, p.rd, p.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	_, err = s.SetProfileUserFavorite(context.Background(), &pbm.ProfileUserFavoriteRequest{UserId: int32(userID), ProfileId: favorite.ProfileID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create favorite profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, favorite)
}

func (p *profilesHandler) DeleteUserProfileFavorite(c *gin.Context) {
	favorite := &models.Favorite{}
	userID, _ := strconv.Atoi(c.Param("userID"))

	if err := c.ShouldBindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !p.tools.ValidateDataIDToken(c, "user_id", userID, p.rd, p.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(profilesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewProfilesServiceClient(conn)

	_, err = s.DeleteProfileUserFavorite(context.Background(), &pbm.ProfileUserFavoriteRequest{UserId: int32(userID), ProfileId: favorite.ProfileID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when delete favorite profile: %v", err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func ConvertProfilePMtoModel(pbmProfile *pbm.Profile) *models.Profile {
	membershipDate, _ := time.Parse(time.RFC3339, pbmProfile.MembershipDueDate)
	createdAt, _ := time.Parse(time.RFC3339, pbmProfile.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbmProfile.UpdatedAt)

	return &models.Profile{
		ID:                pbmProfile.Id,
		UserID:            pbmProfile.UserId,
		Description:       pbmProfile.Description,
		CategoryID:        pbmProfile.CategoryId,
		CoverImageURL:     pbmProfile.CoverImageUrl,
		Website:           pbmProfile.Website,
		Active:            pbmProfile.Active,
		MembershipDueDate: membershipDate,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Rating:            pbmProfile.Rating,
		IsFavorite:        pbmProfile.IsFavorite,
	}
}

func ConvertProfileModeltoPM(profile *models.Profile) *pbm.Profile {
	return &pbm.Profile{
		Id:                profile.ID,
		UserId:            profile.UserID,
		Description:       profile.Description,
		CategoryId:        profile.CategoryID,
		CoverImageUrl:     profile.CoverImageURL,
		Website:           profile.Website,
		Active:            profile.Active,
		MembershipDueDate: profile.MembershipDueDate.String(),
		CreatedAt:         profile.CreatedAt.String(),
		UpdatedAt:         profile.UpdatedAt.String(),
		Rating:            profile.Rating,
		IsFavorite:        profile.IsFavorite,
	}
}
