package routers

import (
	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/handlers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt/middleware"
	"github.com/gin-gonic/gin"
)

func RoutesRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var handler = handlers.NewRoutesHandler(rd, tk)

	route.GET("/v1/profiles/:profileID/routes", middleware.TokenAuthMiddleware(), handler.GetRoutesByProfileID)
	route.GET("/v1/profiles/:profileID/routes/:routeID", middleware.TokenAuthMiddleware(), handler.GetRouteByID)
	route.POST("/v1/profiles/:profileID/routes", middleware.TokenAuthMiddleware(), handler.CreateRoute)
	route.PATCH("/v1/profiles/:profileID/routes/:routeID", middleware.TokenAuthMiddleware(), handler.UpdateRoute)
	route.DELETE("/v1/profiles/:profileID/routes/:routeID", middleware.TokenAuthMiddleware(), handler.DeleteRoute)

	route.GET("/v1/profiles/:profileID/routes/:routeID/schedules", middleware.TokenAuthMiddleware(), handler.GetScheduleByRouteID)
	route.POST("/v1/profiles/:profileID/routes/:routeID/schedules", middleware.TokenAuthMiddleware(), handler.CreateSchedule)
	route.PATCH("/v1/profiles/:profileID/routes/:routeID/schedules/:scheduleID", middleware.TokenAuthMiddleware(), handler.UpdateSchedule)
	route.DELETE("/v1/profiles/:profileID/routes/:routeID/schedules/:scheduleID", middleware.TokenAuthMiddleware(), handler.DeleteSchedule)

}
