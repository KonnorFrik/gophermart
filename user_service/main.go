package main

import (
	"log"
    "gophermart/api"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    apiGroup := router.Group("/api")
    userGroup := apiGroup.Group("/user")

    router.GET("/hello", api.Hello)
    userGroup.POST("/register", api.Register)
    userGroup.POST("/login", api.Login)
    userGroup.POST("/delete", api.Delete)

    log.Println("[SERVER]: listen at:", ":8080")
    router.Run() 
}
