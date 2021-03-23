package main

import (
	"fmt"
	"os"

	"bitbucket.org/edgelabsolutions/loncherapp-api-gateway/app/routers"
	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"
	"bitbucket.org/edgelabsolutions/loncherapp-core/db/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//TODO: delete this
type stubmaindata struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Hours       string `json:"hours"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func main() {

	appAddr := ":" + os.Getenv("GATEWAY_PORT")

	//redis details
	redis_dsn := os.Getenv("REDIS_DSN")
	redis_password := os.Getenv("REDIS_PASSWORD")

	redisClient := redis.NewRedisDB(redis_dsn, redis_password)

	var rd = auth.NewAuth(redisClient)
	var tk = auth.NewToken()

	r := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	//set the routers to gin server
	routers.UserRouter(r, rd, tk)
	routers.LoginRouter(r, rd, tk)
	routers.ProfilesRouter(r, rd, tk)
	routers.RoutesRouter(r, rd, tk)
	routers.MenusRouter(r, rd, tk)
	routers.ReviewsRouter(r, rd, tk)

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:    []string{"Origin"},
		ExposeHeaders:   []string{"Content-Length"},
	}))

	if err := r.Run(appAddr); err != nil {
		fmt.Println("API Gateway service failed, err:%v\n", err)
	}
}
