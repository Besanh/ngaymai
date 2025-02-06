package main

import (
	_ "ngaymai/docs" // Import thư mục docs đã tạo từ `swag init`

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()

	// Định nghĩa endpoint Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Tạo API test
	r.GET("/ping/:name", PingHandler)

	r.Run(":8080")
}

// PingHandler godoc
// @Summary Ping Example
// @Description API để kiểm tra kết nối server
// @Tags example
// @Success 200 {string} string "pong"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
