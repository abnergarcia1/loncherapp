package tools

import (
	"os"
	"strconv"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (t *Tools) ValidateDataIDToken(c *gin.Context, NameID string, ValueID int, rd auth.AuthInterface, tk auth.TokenInterface) bool {
	// CHECK DATA FROM JWT
	metadata, err := tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("something occur when extractTokenMeta auth :%s", err)
		return false
	}
	tokenUserID, err := rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		log.WithFields(log.Fields{
			"redis-dsn": os.Getenv("REDIS_DSN"),
			"uuid":      metadata.TokenUuid,
			"userid":    metadata.UserId,
			"error":     err,
		}).Errorf("something occur when Fetch auth :%s", err)
		return false
	}

	log.WithFields(log.Fields{
		"valueId":     ValueID,
		"tokenUserId": tokenUserID,
	}).Infof("user verification")

	switch NameID {
	case "user_id":
		if ID, _ := strconv.Atoi(tokenUserID.UserID); ID != ValueID {
			return false
		}
	case "profile_id":
		if ID, _ := strconv.Atoi(tokenUserID.ProfileID); ID != ValueID {
			return false
		}
	}

	return true
}
