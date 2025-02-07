package main

import (
	"context"
	v1 "ngaymai/api/v1"
	"ngaymai/common/cache"
	"ngaymai/common/env"
	"ngaymai/common/sql"
	"ngaymai/repository"
	"ngaymai/service"
	"time"

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

	var err error
	if cache.Redis, err = cache.NewRedis(cache.Config{
		Addr:     env.GetStringENV("REDIS_ADDRESS", ""),
		Password: env.GetStringENV("REDIS_PASSWORD", ""),
		DB:       env.GetIntENV("REDIS_DATABASE", 0),
	}); err != nil {
		panic(err)
	}

	sqlClientConfig := sql.SqlConfig{
		SecretKey: env.GetStringENV("SECRET_KEY", ""),
		Host:      env.GetStringENV("PGSQL_HOST", ""),
		Database:  env.GetStringENV("PGSQL_DATABASE", ""),
		Username:  env.GetStringENV("PGSQL_USERNAME", ""),
		Password:  env.GetStringENV("PGSQL_PASSWORD", ""),
		Port:      env.GetIntENV("PGSQL_PORT", 0),
	}
	repository.DBConn = sql.NewSqlClient(sqlClientConfig)

	config = cfg
}

/*
* author: AnhLe
 */
func main() {
	engine := gin.Default()
	initInfo(engine)

	// Define swagger endpoint
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.Run(":" + config.Port)
}

func initInfo(engine *gin.Engine) {
	// Service
	service.NewVideo()

	// Repository
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repository.InitTables(ctx, repository.DBConn)
	repository.InitRepositories()

	// Handler
	v1.NewVideoHandler(engine, service.NewVideo())
}
