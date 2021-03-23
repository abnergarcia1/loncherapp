package handlers

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/routes"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"

	"bitbucket.org/edgelabsolutions/loncherapp-routes-service/app/services"
)

//RoutesAPIServer Routes API Server
type RoutesAPIServer struct {
	RouteService *services.RouteService
}

//GetRoutesByProfileID Get Routes object list from DB using ProfileID
func (r RoutesAPIServer) GetRoutesByProfileID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.Routes, error) {
	routes, err := r.RouteService.GetRoutesByProfileID(id.Id)
	if err != nil {
		return nil, err
	}

	var listRoutes = make([]*pbm.Route, len(*routes))
	for i, route := range *routes {
		listRoutes[i] = ConvertRouteModeltoPM(route)
	}

	return &pbm.Routes{Routes: listRoutes}, nil
}

//GetRouteByID Get Route object from DB using RouteID
func (r RoutesAPIServer) GetRouteByID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.Route, error) {
	route, err := r.RouteService.GetRouteByID(id.Id)
	if err != nil {
		return nil, err
	}

	return ConvertRouteModeltoPM(*route), nil
}

//CreateRoute Create new Route object in DB
func (r RoutesAPIServer) CreateRoute(ctx context.Context, route *pbm.Route) (*pbm.Route, error) {
	modelRoute := ConvertRoutePMtoModel(route)
	routeResponse, err := r.RouteService.CreateRoute(modelRoute)
	if err != nil {
		return nil, err
	}

	return ConvertRouteModeltoPM(*routeResponse), nil
}

//UpdateRoute Update Route object in DB
func (r RoutesAPIServer) UpdateRoute(ctx context.Context, route *pbm.Route) (*pbm.SimpleResponse, error) {
	modelRoute := ConvertRoutePMtoModel(route)
	routeResponse, err := r.RouteService.UpdateRoute(modelRoute)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: routeResponse.Success, Message: routeResponse.Message}, nil
}

//DeleteRoute Delete Route from DB using Route_ID
func (r RoutesAPIServer) DeleteRoute(ctx context.Context, request *pbm.RouteRequest) (*pbm.SimpleResponse, error) {
	err := r.RouteService.DeleteRoute(request.Id, request.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "Route correctly deleted"}, nil
}

//GetScheduleByRouteID Get Schedule objects from DB using RouteID
func (r RoutesAPIServer) GetScheduleByRouteID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.Schedules, error) {
	schedules, err := r.RouteService.GetScheduleByRouteID(id.Id)
	if err != nil {
		return nil, err
	}

	var listSchedules = make([]*pbm.Schedule, len(*schedules))
	for i, schedule := range *schedules {
		listSchedules[i] = ConvertScheduleModeltoPM(&schedule)
	}

	return &pbm.Schedules{Schedules: listSchedules}, nil
}

//CreateSchedule Create New Schedule Object in DB
func (r RoutesAPIServer) CreateSchedule(ctx context.Context, schedule *pbm.Schedule) (*pbm.Schedule, error) {
	modelSchedule := ConvertSchedulePMtoModel(schedule)
	logrus.Info("Converter return: ", modelSchedule)
	scheduleResponse, err := r.RouteService.CreateSchedule(modelSchedule)
	if err != nil {
		return nil, err
	}

	return ConvertScheduleModeltoPM(scheduleResponse), nil
}

// UpdateSchedule Update Schedule object
func (r RoutesAPIServer) UpdateSchedule(ctx context.Context, schedule *pbm.Schedule) (*pbm.SimpleResponse, error) {
	modelSchedule := ConvertSchedulePMtoModel(schedule)
	scheduleResponse, err := r.RouteService.UpdateSchedule(modelSchedule)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: scheduleResponse.Success, Message: scheduleResponse.Message}, nil
}

//DeleteSchedule Delete Schedule Object from DB using ScheduleID
func (r RoutesAPIServer) DeleteSchedule(ctx context.Context, request *pbm.RouteRequest) (*pbm.SimpleResponse, error) {
	err := r.RouteService.DeleteRoute(request.Id, request.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "schedule correctly deleted"}, nil
}

func NewRoutesAPIServer(ctx context.Context) *RoutesAPIServer {
	return &RoutesAPIServer{
		RouteService: services.NewRouteService(ctx),
	}
}

// ConvertRoutePMtoModel Convert Route PBM to Model object
func ConvertRoutePMtoModel(pbmRoute *pbm.Route) *models.Route {

	route := &models.Route{
		ID:            pbmRoute.Id,
		LoncheraID:    pbmRoute.LoncheraId,
		Location:      pbmRoute.Location,
		Address:       pbmRoute.Address,
		Name:          pbmRoute.Name,
		Description:   pbmRoute.Description,
		Order:         pbmRoute.Order,
		Latitude:      pbmRoute.Latitude,
		Longitude:     pbmRoute.Longitude,
		GooglePlaceID: pbmRoute.GooglePlaceId,
	}

	for _, schedulePbm := range pbmRoute.Schedules {
		route.Schedules = append(route.Schedules, *ConvertSchedulePMtoModel(schedulePbm))
	}

	return route
}

// ConvertRouteModeltoPM Convert Route Object to PBM
func ConvertRouteModeltoPM(route models.Route) *pbm.Route {
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

// ConvertSchedulePMtoModel Convert Schedule PBM Object to Model object
func ConvertSchedulePMtoModel(pbmSchedule *pbm.Schedule) *models.Schedule {
	arriveAt, _ := time.Parse("2006-01-02 15:04:05", pbmSchedule.ArriveAt)
	goneAt, _ := time.Parse("2006-01-02 15:04:05", pbmSchedule.GoneAt)

	return &models.Schedule{
		ID:       pbmSchedule.Id,
		RouteID:  pbmSchedule.RouteId,
		Weekday:  pbmSchedule.Weekday,
		ArriveAt: arriveAt,
		GoneAt:   goneAt,
		Active:   pbmSchedule.Active,
	}
}

// ConvertScheduleModeltoPM Convert Schedule Object to PBM
func ConvertScheduleModeltoPM(schedule *models.Schedule) *pbm.Schedule {
	return &pbm.Schedule{
		Id:        schedule.ID,
		RouteId:   schedule.RouteID,
		Weekday:   schedule.Weekday,
		ArriveAt:  schedule.ArriveAt.Format("2006-01-02 15:04:05"),
		GoneAt:    schedule.GoneAt.Format("2006-01-02 15:04:05"),
		CreatedAt: schedule.GoneAt.Format("2006-01-02 15:04:05"),
		Active:    schedule.Active,
	}
}
