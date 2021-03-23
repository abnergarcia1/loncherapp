package main

import (
	"fmt"
	"os"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	"bitbucket.org/edgelabsolutions/loncherapp-core/aws/s3"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/redis"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-images-service/app/routers"

	"github.com/gin-gonic/gin"
)

func main() {

	appAddr := os.Getenv("IMAGE_SERVICE_HOST")

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

	//set the routers to gin server
	routers.ImagesRouter(r, rd, tk)

	/*r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:    []string{"Origin"},
		ExposeHeaders:   []string{"Content-Length"},
	}))*/

	if err := r.Run(appAddr); err != nil {
		fmt.Println("API Gateway service failed, err:%v\n", err)
	}
}
