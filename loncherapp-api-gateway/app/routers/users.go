package routers

import (
	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/handlers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var service = handlers.NewUserHandlers(rd, tk)

	route.GET("/v1/users/:userID", middleware.TokenAuthMiddleware(), service.GetUser)
	//route.GET("/v1/users/email/:email", middleware.TokenAuthMiddleware(), service.GetUserByEmail)
	route.POST("/v1/users", service.CreateUser)
	route.DELETE("/v1/users/:userID", middleware.TokenAuthMiddleware(), nil)

}
