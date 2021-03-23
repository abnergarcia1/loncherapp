package routers

import (
	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/handlers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"github.com/gin-gonic/gin"
)

func LoginRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var service = handlers.NewProfile(rd, tk)

	route.POST("/v1/login", service.Login)
	route.POST("/v1/logout", service.Logout)
	route.POST("/v1/token/refresh", service.Refresh)
}
