package main

import (
	"context"
	"ngaymai/common/cache"
	"ngaymai/common/env"
	"ngaymai/common/sql"
	"ngaymai/repository"
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
		Host:      env.GetStringENV("DB_HOST", ""),
		Database:  env.GetStringENV("DB_DATABASE", ""),
		Username:  env.GetStringENV("DB_USERNAME", ""),
		Password:  env.GetStringENV("DB_PASSWORD", ""),
		Port:      env.GetIntENV("DB_PORT", 0),
	}
	repository.DBConn = sql.NewSqlClient(sqlClientConfig)

	config = cfg
}

/*
* author: AnhLe
 */
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repository.InitTables(ctx, repository.DBConn)
	repository.InitRepositories()

	r := gin.Default()

	// Define swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + config.Port)
}
