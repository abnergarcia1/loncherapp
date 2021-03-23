package routers

import (
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt/middleware"
	"bitbucket.org/edgelabsolutions/loncherapp-images-service/app/handlers"
	"github.com/gin-gonic/gin"
)

func ImagesRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var handler = handlers.NewImagesHandler(rd, tk)

	route.POST("/v1/images/profiles/:profileID", middleware.TokenAuthMiddleware(), handler.UploadCoverImage)
	route.POST("/v1/images/profiles/:profileID/menus/:menuID", middleware.TokenAuthMiddleware(), handler.UploadMenuImage)

}
