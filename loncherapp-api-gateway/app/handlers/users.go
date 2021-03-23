package handlers

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"strconv"
	"time"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/users"

	"google.golang.org/grpc"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"
	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/users"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	port = os.Getenv("USER_SERVICE_HOST")
)

// UserHandler struct
type userHandler struct {
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewUserHandlers(rd auth.AuthInterface, tk auth.TokenInterface) *userHandler {
	return &userHandler{rd, tk}
}

func (u *userHandler) GetUser(c *gin.Context) {

	// CHECK DATA FROM JWT
	metadata, err := u.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("something occur when extractTokenMeta auth :%s", err)
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	tokenUserID, err := u.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		log.WithFields(log.Fields{
			"redis-dsn": os.Getenv("REDIS_DSN"),
			"uuid":      metadata.TokenUuid,
			"userid":    metadata.UserId,
			"error":     err,
		}).Errorf("something occur when Fetch auth :%s", err)
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, _ := strconv.Atoi(tokenUserID.UserID)

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(port, grpc.WithInsecure())
	if err != nil {

		log.WithFields(log.Fields{
			"GRPC-server-host": port,
		}).Errorf("did not connect: %s", err)
		return
	}
	defer conn.Close()

	s := pb.NewUsersServiceClient(conn)

	response, err := s.GetUser(context.Background(), &pbm.UserQuery{Id: int32(userId)})
	if err != nil {
		log.WithFields(log.Fields{
			"user-id": userId,
		}).Errorf("Error trying to get user :%s", err)
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("error when calling rfc: %v", err))
		return
	}

	c.SecureJSON(http.StatusOK, ConvertUserPMtoModel(response))

}

func (u *userHandler) GetUserByEmail(c *gin.Context) {

	email := c.Param("email")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewUsersServiceClient(conn)

	response, err := s.GetUserByEmail(context.Background(), &pbm.UserEmailQuery{Email: email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("error when calling rfc: %v", err))
		return
	}

	c.SecureJSON(http.StatusOK, ConvertUserPMtoModel(response))

}

func (u *userHandler) CreateUser(c *gin.Context) {
	user := &models.User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRequest := ConvertUserModelToPM(user)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewUsersServiceClient(conn)

	response, err := s.CreateUser(context.Background(), userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("error when create user: %v", err))
		return
	}

	c.JSON(http.StatusOK, ConvertUserPMtoModel(response))
}

func ConvertUserPMtoModel(pmuser *pbm.User) *models.User {
	return &models.User{
		ID:            pmuser.Id,
		FirstName:     pmuser.FirstName,
		LastName:      pmuser.LastName,
		Email:         pmuser.Email,
		CreationDate:  time.Unix(pmuser.CreationDate, 0),
		UpdatedDate:   time.Unix(pmuser.UpdatedDate, 0),
		Active:        pmuser.Active,
		TypeID:        uint64(pmuser.TypeId),
		Password:      pmuser.Password,
		ProfileID:     pmuser.ProfileId,
		ProfileActive: pmuser.ProfileActive,
	}
}

func ConvertUserModelToPM(user *models.User) *pbm.User {
	return &pbm.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Active:    user.Active,
		TypeId:    int32(user.TypeID),
		Password:  user.Password,
	}
}
