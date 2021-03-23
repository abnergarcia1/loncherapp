package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/auth"

	"github.com/go-redis/redis"
)

type AuthInterface interface {
	CreateAuth(string, string, *TokenDetails) error
	FetchAuth(string) (*models.TokenPayload, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type service struct {
	client *redis.Client
}

var _ AuthInterface = &service{}

func NewAuth(client *redis.Client) *service {
	return &service{client: client}
}

type AccessDetails struct {
	TokenUuid string
	UserId    string
	ProfileId string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

//Save token metadata to Redis
func (tk *service) CreateAuth(userId string, profileId string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	objectUserProfile, _ := json.Marshal(&models.TokenPayload{UserID: userId, ProfileID: profileId})

	atCreated, err := tk.client.Set(td.TokenUuid, objectUserProfile, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.client.Set(td.RefreshUuid, objectUserProfile, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

//Check the metadata saved
func (tk *service) FetchAuth(tokenUuid string) (*models.TokenPayload, error) {
	redisResponse, err := tk.client.Get(tokenUuid).Result()
	if err != nil {
		return nil, err
	}
	response := models.TokenPayload{}

	if err = json.Unmarshal([]byte(redisResponse), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

//Once a user row in the token table
func (tk *service) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := tk.client.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (tk *service) DeleteRefresh(refreshUuid string) error {
	//delete refresh token
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
