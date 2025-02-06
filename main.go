package main

import (
	"ngaymai/common/env"

	_ "ngaymai/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Config struct {
	Port    string
	LogFile string
}

var config Config

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := Config{
		Port:    env.GetStringENV("PORT", "8000"),
		LogFile: env.GetStringENV("LOG_FILE", "log/console.log"),
	}

	// var err error
	// if cache.Redis, err = cache.NewRedis(cache.Config{
	// 	Addr:     env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
	// 	Password: env.GetStringENV("REDIS_PASSWORD", ""),
	// 	DB:       env.GetIntENV("REDIS_DATABASE", 0),
	// }); err != nil {
	// 	panic(err)
	// }

	config = cfg
}

/*
* author: AnhLe
 */
func main() {
	r := gin.Default()

	// Define swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + config.Port)
}
