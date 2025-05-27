package main

import (
	"log"
    "gophermart/api"
    "gophermart/middleware"

	"github.com/gin-gonic/gin"
    _ "gophermart/docs"
    swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title gophermart
// @version 0.1

func main() {
    router := gin.Default()

    apiGroup := router.Group("/api")
    userGroup := apiGroup.Group("/user")

    userGroup.POST("/register", api.UserRegister)
    userGroup.POST("/login", api.UserLogin)

    UserAuthGroup := apiGroup.Group("/user")
    UserAuthGroup.Use(middleware.JWTAuthenticate)
    UserAuthGroup.DELETE("/delete", api.UserDelete)
    UserAuthGroup.POST("/orders", api.NewOrder)
    UserAuthGroup.GET("/orders", api.AllOrders)

    log.Println("[SERVER]: listen at:", ":8080")
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    router.Run() 
}
