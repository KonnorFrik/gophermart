package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
    cookie, err := c.Cookie("basic")

    if err != nil {
        log.Printf("[/hello]: cookie error: %v\n", err)
        cookie = "NotSet"
        c.SetCookie("basic", "usage", 3600, "/", "", false, true)
    }

    log.Printf("[/hello]: cookie value: %q\n", cookie)
    c.JSON(http.StatusOK, gin.H{
        "status": "OK",
    })
}

func main() {
    router := gin.Default()
    router.GET("/hello", hello)
    router.Run() // :8080
}
