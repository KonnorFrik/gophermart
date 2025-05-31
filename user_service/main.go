package main

import (
	"gophermart/api"
	"gophermart/middleware"
	"log"
	"os"

	_ "gophermart/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title gophermart
// @version 0.1

func main() {
    router := gin.Default()
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    apiGroup := router.Group("/api")
    userGroup := apiGroup.Group("/user")

    userGroup.POST("/register", api.UserRegister)
    userGroup.POST("/login", api.UserLogin)

    UserAuthGroup := apiGroup.Group("/user")
    UserAuthGroup.Use(middleware.JWTAuthenticate)
    UserAuthGroup.DELETE("/delete", api.UserDelete)
    UserAuthGroup.POST("/orders", api.NewOrder)
    UserAuthGroup.GET("/orders", api.AllOrders)

    log.Printf("[SERVER]: listen at: :%s\n", os.Getenv("PORT"))
    router.Run() 
}
