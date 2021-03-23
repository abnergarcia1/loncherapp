package routers

import (
	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/handlers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt/middleware"
	"github.com/gin-gonic/gin"
)

func MenusRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var handler = handlers.NewMenusHandler(rd, tk)

	route.GET("/v1/profiles/:profileID/menus", middleware.TokenAuthMiddleware(), handler.GetMenuByProfileID)
	route.GET("/v1/profiles/:profileID/menus/:menuID", middleware.TokenAuthMiddleware(), handler.GetMenuByID)
	route.POST("/v1/profiles/:profileID/menus", middleware.TokenAuthMiddleware(), handler.CreateMenu)
	route.PATCH("/v1/profiles/:profileID/menus/:menuID", middleware.TokenAuthMiddleware(), handler.UpdateMenu)
	route.DELETE("/v1/profiles/:profileID/menus/:menuID", middleware.TokenAuthMiddleware(), handler.DeleteMenu)

}
