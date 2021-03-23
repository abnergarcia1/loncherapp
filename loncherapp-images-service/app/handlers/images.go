package handlers

import (
	"context"
	"net/http"
	"os"
	filepath2 "path/filepath"
	"strconv"

	"bitbucket.org/edgelabsolutions/loncherapp-images-service/app/services"

	uuid "github.com/satori/go.uuid"

	"bitbucket.org/edgelabsolutions/loncherapp-core/tools"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
)

type ImagesHandler struct {
	rd           auth.AuthInterface
	tk           auth.TokenInterface
	tools        tools.Tools
	RouteService *services.ImagesService
}

func NewImagesHandler(rd auth.AuthInterface, tk auth.TokenInterface) *ImagesHandler {
	return &ImagesHandler{rd, tk, tools.Tools{}, services.NewImagesService(context.Background())}
}

func (s *ImagesHandler) UploadCoverImage(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	if !s.tools.ValidateDataIDToken(c, "profile_id", profileID, s.rd, s.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	filename := "Cover-" + uuid.NewV4().String()

	newURL, err := uploadToS3(c, filename)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "Failed to upload cover image file: ",
			"error_description": err,
		})
		return
	}

	if err := s.RouteService.UpdateProfileImageURL(int32(profileID), newURL); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "Failed to update profile image url in DB: ",
			"error_description": err,
		})
		return
	}

	c.JSON(http.StatusOK, "Image uploaded correctly")
}

func (s *ImagesHandler) UploadMenuImage(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	menuID, _ := strconv.Atoi(c.Param("menuID"))

	if !s.tools.ValidateDataIDToken(c, "profile_id", profileID, s.rd, s.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	filename := "Menu-" + uuid.NewV4().String()

	newURL, err := uploadToS3(c, filename)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "Failed to upload menu image file: ",
			"error_description": err,
		})
		return
	}

	if err := s.RouteService.UpdateMenuImageURL(int32(profileID), int32(menuID), newURL); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "Failed to update menu image url in DB: ",
			"error_description": err,
		})
		return
	}

	c.JSON(http.StatusOK, "Image uploaded correctly")
}

func uploadToS3(c *gin.Context, imageName string) (string, error) {
	sess := c.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)

	MyBucket := os.Getenv("BUCKET_NAME")

	file, header, err := c.Request.FormFile("image")
	filename := header.Filename
	fileExt := filepath2.Ext(filename)
	newFilename := imageName + fileExt

	//upload to the s3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String("images/" + newFilename),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	filepath := "https://d1jjbjbyvtug09.cloudfront.net/" + newFilename

	return filepath, nil
}
