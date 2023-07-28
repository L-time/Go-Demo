package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	//新建一个Get路由

	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "Hello,World!",
		})
	})

	//在localhost:8080上启动服务

	r.Run(":8080")
}
