package routers

import (
	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/handlers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt/middleware"
	"github.com/gin-gonic/gin"
)

func ProfilesRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var service = handlers.NewProfilesHandler(rd, tk)

	route.GET("/v1/profiles/:profileID", middleware.TokenAuthMiddleware(), service.GetProfileByID)
	route.POST("/v1/profiles", middleware.TokenAuthMiddleware(), service.CreateProfile)
	route.PATCH("/v1/profiles/:profileID", middleware.TokenAuthMiddleware(), service.UpdateProfile)
	//route.DELETE("/v1/profiles/:profileID", middleware.TokenAuthMiddleware(), service.)
	route.GET("/v1/users/:userID/profile", middleware.TokenAuthMiddleware(), service.GetProfileByUserID)

	route.GET("/v1/users/:userID/favorites", middleware.TokenAuthMiddleware(), service.GetProfilesByUserFavorites)
	route.POST("/v1/users/:userID/favorites", middleware.TokenAuthMiddleware(), service.SetUserProfileFavorite)
	route.DELETE("/v1/users/:userID/favorites", middleware.TokenAuthMiddleware(), service.DeleteUserProfileFavorite)

	route.GET("/v1/categories/:categoryID/profiles", middleware.TokenAuthMiddleware(), service.GetProfilesByCategories)
	route.GET("/v1/profiles", middleware.TokenAuthMiddleware(), service.GetProfilesByUserLocation) //queryparms

}
