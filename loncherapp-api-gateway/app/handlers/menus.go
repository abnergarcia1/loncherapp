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

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/menus"

	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/menus"
	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
)

var (
	menusServiceHost = os.Getenv("MENUS_SERVICE_HOST")
)

// RoutesHandler struct
type menusHandler struct {
	rd    auth.AuthInterface
	tk    auth.TokenInterface
	tools tools.Tools
}

func NewMenusHandler(rd auth.AuthInterface, tk auth.TokenInterface) *menusHandler {
	return &menusHandler{rd, tk, tools.Tools{}}
}

func (r *menusHandler) GetMenuByID(c *gin.Context) {
	menuID, _ := strconv.Atoi(c.Param("menuID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(menusServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewMenusServiceClient(conn)

	response, err := s.GetMenuByID(context.Background(), &pbm.SimpleRequestByID{Id: int32(menuID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting menu item: %v", err))
		return
	}

	c.JSON(http.StatusOK, ConvertMenuPMtoModel(response))
}

func (r *menusHandler) GetMenuByProfileID(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(menusServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewMenusServiceClient(conn)

	response, err := s.GetMenuByProfileID(context.Background(), &pbm.SimpleRequestByID{Id: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting menu items: %v", err))
		return
	}

	var listMenu = make([]models.Menu, len(response.Menus))
	for i, menu := range response.Menus {
		listMenu[i] = *ConvertMenuPMtoModel(menu)
	}

	c.JSON(http.StatusOK, listMenu)
}

func (r *menusHandler) CreateMenu(c *gin.Context) {
	menu := &models.Menu{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	menuRequest := ConvertMenuModelToPM(menu)
	menuRequest.LoncheraId = int32(profileID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(menusServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewMenusServiceClient(conn)

	response, err := s.CreateMenu(context.Background(), menuRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create menu item: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *menusHandler) UpdateMenu(c *gin.Context) {
	menu := &models.Menu{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	menuID, _ := strconv.Atoi(c.Param("menuID"))

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	menuRequest := ConvertMenuModelToPM(menu)
	menuRequest.LoncheraId = int32(profileID)
	menuRequest.Id = int32(menuID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(menusServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewMenusServiceClient(conn)

	response, err := s.UpdateMenu(context.Background(), menuRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when updating menu item: %v", err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (r *menusHandler) DeleteMenu(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	menuID, _ := strconv.Atoi(c.Param("menuID"))

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(menusServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewMenusServiceClient(conn)

	res, err := s.DeleteMenu(context.Background(), &pbm.MenuRequest{Id: int32(menuID), ProfileId: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when delete menu item %v : %v", menuID, err))
		return
	}

	c.JSON(http.StatusNotFound, sharedModels.SimpleResponse{Success: res.Success, Message: res.Message})
}

func ConvertMenuPMtoModel(pbmMenu *pbm.Menu) *models.Menu {
	var createdAt, updatedAt time.Time
	if pbmMenu.CreatedAt != "" {
		createdAt, _ = time.Parse("2006-01-02 15:04:05", pbmMenu.CreatedAt)
	}

	if pbmMenu.UpdatedAt != "" {
		updatedAt, _ = time.Parse("2006-01-02 15:04:05", pbmMenu.UpdatedAt)
	}

	menu := &models.Menu{
		ID:          pbmMenu.Id,
		LoncheraID:  pbmMenu.LoncheraId,
		Name:        pbmMenu.Name,
		Description: pbmMenu.Description,
		Price:       pbmMenu.Price,
		Currency:    pbmMenu.Currency,
		ImageURL:    pbmMenu.ImageUrl,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return menu
}

func ConvertMenuModelToPM(menu *models.Menu) *pbm.Menu {
	menuModel := &pbm.Menu{
		Id:          menu.ID,
		LoncheraId:  menu.LoncheraID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       menu.Price,
		Currency:    menu.Currency,
		ImageUrl:    menu.ImageURL,
	}

	return menuModel
}
