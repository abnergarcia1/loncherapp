package grpc_handlers

import (
	"context"
	"fmt"
	"time"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/users"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"

	"bitbucket.org/edgelabsolutions/loncherapp-users-service/app/services"
)

type UsersAPIServer struct {
	UserService *services.UserService
}

func NewUsersAPIServer(ctx context.Context) *UsersAPIServer {
	return &UsersAPIServer{
		UserService: services.NewUserService(ctx),
	}
}

func (u UsersAPIServer) GetUser(ctx context.Context, userQuery *pbm.UserQuery) (*pbm.User, error) {
	fmt.Println("GET USER function called...")

	user, err := u.UserService.GetUser(int(userQuery.Id))
	if err != nil {
		return nil, err
	}

	return ConvertUserModeltoPM(user), nil

}

func (u UsersAPIServer) GetUserByEmail(ctx context.Context, userQuery *pbm.UserEmailQuery) (*pbm.User, error) {

	user, err := u.UserService.GetUserByEmail(userQuery.Email)
	if err != nil {
		return nil, err
	}

	return ConvertUserModeltoPM(user), nil
}

func (u UsersAPIServer) GetUserByEmailPassword(ctx context.Context, userQuery *pbm.UserLoginRequest) (*pbm.User, error) {

	user, err := u.UserService.GetUserByEmailPassword(userQuery.Email, userQuery.Password)
	if err != nil {
		return nil, err
	}

	return ConvertUserModeltoPM(user), nil
}

func (u UsersAPIServer) CreateUser(ctx context.Context, user *pbm.User) (*pbm.User, error) {

	userResponse, err := u.UserService.CreateUser(ConvertUserPMtoModel(user))
	if err != nil {
		return nil, err
	}

	return ConvertUserModeltoPM(userResponse), nil
}

func (u UsersAPIServer) DeactivateUser(ctx context.Context, userQuery *pbm.UserQuery) (*pbm.SimpleResponse, error) {
	panic("implement me")
}

func (u UsersAPIServer) UpdateUser(ctx context.Context, user *pbm.User) (*pbm.SimpleResponse, error) {
	panic("implement me")
}

func ConvertUserPMtoModel(pmuser *pbm.User) *models.User {
	return &models.User{
		ID:           pmuser.Id,
		TypeID:       uint64(pmuser.TypeId),
		FirstName:    pmuser.FirstName,
		LastName:     pmuser.LastName,
		Email:        pmuser.Email,
		CreationDate: time.Unix(pmuser.CreationDate, 0),
		UpdatedDate:  time.Unix(pmuser.UpdatedDate, 0),
		Active:       pmuser.Active,
		Password:     pmuser.Password,
	}
}

func ConvertUserModeltoPM(user *models.User) *pbm.User {
	return &pbm.User{
		Id:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		CreationDate:  user.CreationDate.Unix(),
		UpdatedDate:   user.UpdatedDate.Unix(),
		Active:        user.Active,
		TypeId:        int32(user.TypeID),
		Password:      user.Password,
		ProfileId:     user.ProfileID,
		ProfileActive: user.ProfileActive,
	}
}
