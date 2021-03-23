package app

import (
	"fmt"
	"os"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"

	"bitbucket.org/edgelabsolutions/loncherapp-core/aws/s3"
	"bitbucket.org/edgelabsolutions/loncherapp-core/db/redis"
	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"
	"bitbucket.org/edgelabsolutions/loncherapp-payments-service/app/payments/paypal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	appAddr := os.Getenv("PAYMENTS_SERVICE_HOST")

	//redis details
	redis_dsn := os.Getenv("REDIS_DSN")
	redis_password := os.Getenv("REDIS_PASSWORD")

	sql.Init(
		sql.SetConnectionString(os.Getenv("LONCHERAPP_DB_CONNECTION")),
	)

	redisClient := redis.NewRedisDB(redis_dsn, redis_password)

	var rd = auth.NewAuth(redisClient)
	var tk = auth.NewToken()

	sess := s3.ConnectAws()

	r := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
		gin.Recovery()
	})

	r.Use(cors.New(cors.Config{
		AllowMethods:  []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:  []string{"Origin", " X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOrigins:  []string{"*"},
	}))

	//set the routers to gin server
	paypal.Routers(r, rd, tk)

	if err := r.Run(appAddr); err != nil {
		fmt.Println("API Gateway service failed, err:%v\n", err)
	}
}
