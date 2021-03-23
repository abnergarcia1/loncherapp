package services

import (
	"context"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/routes"
	sharedModels "bitbucket.org/edgelabsolutions/loncherapp-core/models/shared"
)

type RouteService struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewRouteService(context context.Context) *RouteService {
	return &RouteService{
		ctx: context,
		db:  sql.NewClient(),
	}
}

func (p *RouteService) CreateRoute(route *models.Route) (*models.Route, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	res, err := ses.Execute("INSERT INTO Routes(Lonchera_ID, Location, Address, Name, `Description`, `Order`, Created_At, Updated_At, Latitude, Longitude, Google_Place_ID) VALUES(?,?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?,?,?) ",
		route.LoncheraID, route.Location, route.Address, route.Name, route.Description, route.Order, route.Latitude, route.Longitude, route.GooglePlaceID)
	if err != nil {
		log.WithFields(log.Fields{
			"route": route,
		}).Errorf("Error when trying to create route in DB: %v", err)
		return nil, err
	}

	objectID, err := res.LastInsertId()
	if err == nil {
		route.ID = int32(objectID)
	}

	return route, nil
}

func (p *RouteService) UpdateRoute(route *models.Route) (*sharedModels.SimpleResponse, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	result, err := ses.Execute("UPDATE Routes SET Location=?, Address=?, Name=?, Description=?, `Order`=?, Updated_At=CURRENT_TIMESTAMP, Latitude=?, Longitude=?, Google_Place_ID=?  WHERE ID=?",
		route.Location, route.Address, route.Name, route.Description, route.Order, route.Latitude, route.Longitude, route.GooglePlaceID, route.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"route": route,
		}).Errorf("Error when trying to update profile in DB: %v", err)
		return nil, err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"route": route,
		}).Errorf("Error when trying to update route in DB: %v", sql.DbErrNoDocuments)
		return nil, sql.DbErrNoDocuments
	}

	return &sharedModels.SimpleResponse{Success: true, Message: "Updated route correctly"}, nil
}

func (p *RouteService) GetRoutesByProfileID(profileID int32) (*[]models.Route, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	var listRoutes []models.Route

	err := ses.Select(&listRoutes, `SELECT * FROM Routes WHERE Lonchera_ID=?`, profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
		}).Errorf("Error when trying to get routes by LoncheraID  in DB: %v", err)
		return nil, err
	}

	return &listRoutes, nil
}

func (p *RouteService) GetRouteByID(routeID int32) (*models.Route, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	var route models.Route

	err := ses.QueryOne(&route, `SELECT * FROM Routes WHERE ID=?`, routeID)
	if err != nil {
		log.WithFields(log.Fields{
			"routeID": routeID,
		}).Errorf("Error when trying to get route by ID in DB: %v", err)
		return nil, err
	}

	return &route, nil
}

func (p *RouteService) DeleteRoute(routeID int32, profileID int32) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		p.db.Disconnect()
	}()

	result, err := ses.Execute(`DELETE FROM Routes WHERE ID=? AND Lonchera_ID=? `, routeID, profileID)

	if err != nil {
		log.WithFields(log.Fields{
			"routeID":   routeID,
			"profileID": profileID,
		}).Errorf("Error when trying to delete route in DB: %v", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"routeID":   routeID,
			"profileID": profileID,
		}).Errorf("Error when trying to delete route in DB: %v", sql.DbErrNoDocuments)
		return sql.DbErrNoDocuments
	}

	return nil
}

//Schedules Routes

func (p *RouteService) CreateSchedule(schedule *models.Schedule) (*models.Schedule, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	log.Info("schedule object in CreateSchedule method: ", schedule)

	res, err := ses.Execute(`INSERT INTO Routes_Schedule(Route_ID, Weekday, Arrive_At, Gone_At, Created_At, Active) VALUE(?,?,?,?,CURRENT_TIMESTAMP, 1) `,
		schedule.RouteID, schedule.Weekday, schedule.ArriveAt.Format("2006-01-02 15:04:05"), schedule.GoneAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.WithFields(log.Fields{
			"schedule": schedule,
		}).Errorf("Error when trying to create schedule in DB: %v", err)
		return nil, err
	}

	objectID, err := res.LastInsertId()
	if err == nil {
		schedule.ID = int32(objectID)
	}

	return schedule, nil
}

func (p *RouteService) UpdateSchedule(schedule *models.Schedule) (*sharedModels.SimpleResponse, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()
	log.Info(schedule)
	result, err := ses.Execute(`UPDATE Routes_Schedule SET Arrive_At=?, Gone_At=?, Active=?  WHERE ID=?`,
		schedule.ArriveAt.Format("2006-01-02 15:04:05"), schedule.GoneAt.Format("2006-01-02 15:04:05"), schedule.Active, schedule.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"schedule": schedule,
		}).Errorf("Error when trying to update schedule in DB: %v", err)
		return nil, err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"schedule": schedule,
		}).Errorf("Error when trying to update schedule in DB: %v", sql.DbErrNoDocuments)
		return nil, sql.DbErrNoDocuments
	}

	return &sharedModels.SimpleResponse{Success: true, Message: "Updated schedule correctly"}, nil
}

func (p *RouteService) GetScheduleByRouteID(routeID int32) (*[]models.Schedule, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		p.db.Disconnect()
	}()

	var listSchedule []models.Schedule

	err := ses.Select(&listSchedule, `SELECT * FROM Routes_Schedule WHERE Route_ID=?`, routeID)
	if err != nil {
		log.WithFields(log.Fields{
			"routeID": routeID,
		}).Errorf("Error when trying to get schedule by RouteID in DB: %v", err)
		return nil, err
	}

	return &listSchedule, nil
}

func (p *RouteService) DeleteSchedule(scheduleID int32, profileID int32) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		p.db.Disconnect()
	}()

	result, err := ses.Execute(`DELETE RS FROM Routes_Schedule RS INNER JOIN Routes ON RS.Route_ID = Routes.ID WHERE RS.ID=? AND Routes.Lonchera_ID=? `, scheduleID, profileID)

	if err != nil {
		log.WithFields(log.Fields{
			"scheduleID": scheduleID,
			"profileID":  profileID,
		}).Errorf("Error when trying to delete route in DB: %v", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		log.WithFields(log.Fields{
			"scheduleID": scheduleID,
			"profileID":  profileID,
		}).Errorf("Error when trying to delete route in DB: %v", sql.DbErrNoDocuments)
		return sql.DbErrNoDocuments
	}

	return nil
}
