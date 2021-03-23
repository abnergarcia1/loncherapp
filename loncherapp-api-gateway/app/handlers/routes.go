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

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/routes"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"
	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/routes"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
)

var (
	routesServiceHost = os.Getenv("ROUTES_SERVICE_HOST")
)

// RoutesHandler struct
type routesHandler struct {
	rd    auth.AuthInterface
	tk    auth.TokenInterface
	tools tools.Tools
}

func NewRoutesHandler(rd auth.AuthInterface, tk auth.TokenInterface) *routesHandler {
	return &routesHandler{rd, tk, tools.Tools{}}
}

func (r *routesHandler) GetRouteByID(c *gin.Context) {
	routeID, _ := strconv.Atoi(c.Param("routeID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.GetRouteByID(context.Background(), &pbm.SimpleRequestByID{Id: int32(routeID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting profile: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *routesHandler) GetRoutesByProfileID(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.GetRoutesByProfileID(context.Background(), &pbm.SimpleRequestByID{Id: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting profile: %v", err))
		return
	}

	var listRoutes = make([]models.Route, len(response.Routes))
	for i, route := range response.Routes {
		listRoutes[i] = *ConvertRoutePMtoModel(route)
	}

	c.JSON(http.StatusOK, listRoutes)
}

func (r *routesHandler) CreateRoute(c *gin.Context) {
	route := &models.Route{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	routeRequest := ConvertRouteModeltoPM(route)
	routeRequest.LoncheraId = int32(profileID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.CreateRoute(context.Background(), routeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create route: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *routesHandler) UpdateRoute(c *gin.Context) {
	route := &models.Route{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	routeID, _ := strconv.Atoi(c.Param("routeID"))

	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	routeRequest := ConvertRouteModeltoPM(route)
	routeRequest.LoncheraId = int32(profileID)
	routeRequest.Id = int32(routeID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.UpdateRoute(context.Background(), routeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when updating route: %v", err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (r *routesHandler) DeleteRoute(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	routeID, _ := strconv.Atoi(c.Param("routeID"))

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	res, err := s.DeleteRoute(context.Background(), &pbm.RouteRequest{Id: int32(routeID), ProfileId: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when delete routeID %v : %v", routeID, err))
		return
	}

	c.JSON(http.StatusNotFound, sharedModels.SimpleResponse{Success: res.Success, Message: res.Message})
}

func (r *routesHandler) GetScheduleByRouteID(c *gin.Context) {
	routeID, _ := strconv.Atoi(c.Param("routeID"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.GetScheduleByRouteID(context.Background(), &pbm.SimpleRequestByID{Id: int32(routeID)})
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Error when getting schedule: %v", err))
		return
	}

	var listSchedules = make([]models.Schedule, len(response.Schedules))
	for i, schedule := range response.Schedules {
		listSchedules[i] = *ConvertSchedulePMtoModel(schedule)
	}

	c.JSON(http.StatusOK, listSchedules)
}

func (r *routesHandler) CreateSchedule(c *gin.Context) {
	schedule := &models.Schedule{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	scheduleRequest := ConvertScheduleModeltoPM(schedule)
	scheduleRequest.ProfileId = int32(profileID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.CreateSchedule(context.Background(), scheduleRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when create schedule: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *routesHandler) UpdateSchedule(c *gin.Context) {
	schedule := &models.Schedule{}
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	scheduleID, _ := strconv.Atoi(c.Param("scheduleID"))

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	scheduleRequest := ConvertScheduleModeltoPM(schedule)
	scheduleRequest.ProfileId = int32(profileID)
	scheduleRequest.Id = int32(scheduleID)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	response, err := s.UpdateSchedule(context.Background(), scheduleRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when updating schedule: %v", err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (r *routesHandler) DeleteSchedule(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	scheduleID, _ := strconv.Atoi(c.Param("scheduleID"))

	if !r.tools.ValidateDataIDToken(c, "profile_id", profileID, r.rd, r.tk) {
		c.JSON(http.StatusUnauthorized, "Not authorized")
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(routesServiceHost, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when connecting to internal service: %v", err))
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := pb.NewRoutesServiceClient(conn)

	_, err = s.DeleteRoute(context.Background(), &pbm.RouteRequest{Id: int32(scheduleID), ProfileId: int32(profileID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error when delete scheduleID %v : %v", scheduleID, err))
		return
	}

	c.JSON(http.StatusNotFound, nil)
}

func ConvertRoutePMtoModel(pbmRoute *pbm.Route) *models.Route {
	var createdAt, updatedAt time.Time
	if pbmRoute.CreatedAt != "" {
		createdAt, _ = time.Parse(time.RFC3339, pbmRoute.CreatedAt)
	}

	if pbmRoute.CreatedAt != "" {
		updatedAt, _ = time.Parse(time.RFC3339, pbmRoute.UpdatedAt)
	}

	route := &models.Route{
		ID:            pbmRoute.Id,
		LoncheraID:    pbmRoute.LoncheraId,
		Location:      pbmRoute.Location,
		Address:       pbmRoute.Address,
		Name:          pbmRoute.Name,
		Description:   pbmRoute.Description,
		Order:         pbmRoute.Order,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Latitude:      pbmRoute.Latitude,
		Longitude:     pbmRoute.Longitude,
		GooglePlaceID: pbmRoute.GooglePlaceId,
	}

	for _, schedulePbm := range pbmRoute.Schedules {
		route.Schedules = append(route.Schedules, *ConvertSchedulePMtoModel(schedulePbm))
	}

	return route
}

func ConvertRouteModeltoPM(route *models.Route) *pbm.Route {
	routeModel := &pbm.Route{
		Id:            route.ID,
		LoncheraId:    route.LoncheraID,
		Location:      route.Location,
		Address:       route.Address,
		Name:          route.Name,
		Description:   route.Description,
		Order:         route.Order,
		CreatedAt:     route.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     route.UpdatedAt.Format("2006-01-02 15:04:05"),
		Latitude:      route.Latitude,
		Longitude:     route.Longitude,
		GooglePlaceId: route.GooglePlaceID,
	}

	for _, scheduleModel := range route.Schedules {
		routeModel.Schedules = append(routeModel.Schedules, ConvertScheduleModeltoPM(&scheduleModel))
	}

	return routeModel

}

func ConvertSchedulePMtoModel(pbmSchedule *pbm.Schedule) *models.Schedule {
	arriveAt, err := time.Parse("2006-01-02 15:04:05", pbmSchedule.ArriveAt)
	if err != nil {
		log.Error(err.Error())
	}
	goneAt, _ := time.Parse("2006-01-02 15:04:05", pbmSchedule.GoneAt)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", pbmSchedule.CreatedAt)

	log.Info(pbmSchedule)

	return &models.Schedule{
		ID:        pbmSchedule.Id,
		RouteID:   pbmSchedule.RouteId,
		Weekday:   pbmSchedule.Weekday,
		ArriveAt:  arriveAt,
		GoneAt:    goneAt,
		CreatedAt: createdAt,
		Active:    pbmSchedule.Active,
	}
}

func ConvertScheduleModeltoPM(schedule *models.Schedule) *pbm.Schedule {
	return &pbm.Schedule{
		Id:       schedule.ID,
		RouteId:  schedule.RouteID,
		Weekday:  schedule.Weekday,
		ArriveAt: schedule.ArriveAt.Format("2006-01-02 15:04:05"),
		GoneAt:   schedule.GoneAt.Format("2006-01-02 15:04:05"),
		Active:   schedule.Active,
	}
}
