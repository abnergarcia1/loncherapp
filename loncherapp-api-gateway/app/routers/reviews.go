package routers

import (
	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/handlers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt/middleware"
	"github.com/gin-gonic/gin"
)

func ReviewsRouter(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var handler = handlers.NewReviewsHandler(rd, tk)

	route.GET("/v1/profiles/:profileID/reviews", middleware.TokenAuthMiddleware(), handler.GetReviewsByProfileID)
	route.GET("/v1/profiles/:profileID/rating", middleware.TokenAuthMiddleware(), handler.GetAverageRatingByProfileID)
	route.POST("/v1/profiles/:profileID/reviews", middleware.TokenAuthMiddleware(), handler.CreateReview)
	route.DELETE("/v1/profiles/:profileID/reviews/:reviewID", middleware.TokenAuthMiddleware(), handler.DeleteReview)

}
