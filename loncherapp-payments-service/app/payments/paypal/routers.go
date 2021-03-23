package paypal

import (
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"github.com/gin-gonic/gin"
)

func Routers(route *gin.Engine, rd auth.AuthInterface, tk auth.TokenInterface) {
	var handler = NewHandler(rd, tk)

	route.POST("/v1/payments/paypal/webhooks", handler.HandleWebhookRequest)
	//route.POST("/v1/payments/paypal/profiles/:profileID/orders", middleware.TokenAuthMiddleware(), handler.CreatePaypalOrder)
	//route.POST("/v1/payments/paypal/profiles/:profileID/order/:orderID/capture", middleware.TokenAuthMiddleware(), handler.CapturePayPalOrder)
	route.POST("/v1/payments/paypal/profiles/:profileID/orders", handler.CreatePaypalOrder)
	route.POST("/v1/payments/paypal/profiles/:profileID/orders/:orderID/capture", handler.CapturePayPalOrder)

}
