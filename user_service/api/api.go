package api

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
    "gophermart/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Hello - just for simple ping and check cookie
func Hello(c *gin.Context) {
    cookie, err := c.Cookie("basic")

    if err != nil {
        log.Printf("[/hello]: cookie error: %v\n", err)
        cookie = "NotSet"
        c.SetCookie("basic", "usage", 3600, "/", "", false, true)
    }

    log.Printf("[/hello]: cookie value: %q\n", cookie)
    c.JSON(http.StatusOK, gin.H{
        "status": "OK",
        "cookie": cookie,
    })
}

var secretKey = []byte("aBoBa")

// Register - create a new user
// TODO: now it create a user in memory, need to use db
func Register(c *gin.Context) {
    var user model.User
    err := c.ShouldBindBodyWithJSON(&user)

    if err != nil || user.Login == "" || user.Password == "" {
        log.Printf("[REGISTER]: Error on binding: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    userDB, err := model.NewUser(user.Login, user.Password, "")

    if err != nil {
        log.Printf("[REGISTER]: %v: %q\n", err, user.Login)
        c.JSON(http.StatusConflict, gin.H{})
        return
    }

    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": "gophermart",
        "sub": user.Login,
    })
    str, err := tok.SignedString(secretKey)

    if err != nil {
        log.Printf("[REGISTER]: Error on signing: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{})
        return
    }

    userDB.Token = str
    c.SetCookie("jwt", str, 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{})
}

func userDBbyCreadentials(login, password string) (*model.User, bool) {
    user, err := model.UserByLogin(login)

    if err != nil || user.Password != password {
        log.Printf("[LOGIN]: %v: %q\n", err, login)
        return nil, false
    }

    return user, true
}

// basicAuthCredentials - returns User type from pair login:password from basic auth or false
func basicAuthCredentials(headerAuthString string) (*model.User, bool) {
    if len(headerAuthString) == 0 {
        return nil, false
    }

    headerSplitted := strings.Split(headerAuthString, " ")

    if len(headerSplitted) != 2 || headerSplitted[0] != "Basic" {
        log.Printf("[LOGIN/basicAuth]: Wrong data: %q\n", headerAuthString)
        return nil, false
    }

    // BUG: can't decode pair "admin:admin"
    credPair, err := base64.RawStdEncoding.DecodeString(headerSplitted[1])

    if err != nil {
        log.Printf("[LOGIN/basicAuth]: Error on decode: %q %v\n", headerAuthString, err)
        return nil, false
    }

    credPairSplitted := strings.Split(string(credPair), ":")

    if len(credPairSplitted) != 2 {
        log.Printf("[LOGIN/basicAuth]: Wrong pair: %q\n", credPair)
        return nil, false
    }

    login, password := credPairSplitted[0], credPairSplitted[1]
    return userDBbyCreadentials(login, password)
}

// basicAuthCredentials - returns User type from login:password of post body or false
func postAuthCredentials(c *gin.Context) (*model.User, bool) {
    var user model.User
    err := c.ShouldBindBodyWithJSON(&user)

    if err != nil || user.Login == "" || user.Password == "" {
        log.Printf("[LOGIN]: Error on binding: %v\n", err)
        return nil, false
    }

    return userDBbyCreadentials(user.Login, user.Password)
}

func Login(c *gin.Context) {
    var userDB *model.User
    var ok bool

    userDB, ok = postAuthCredentials(c)

    if !ok {
        userDB, ok = basicAuthCredentials(c.GetHeader("Authorization"))
    }

    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }

    if userDB.Token == "" {
        // unreachable in normal flow
        // TODO: rewrite this to normal behaviour
        log.Printf("[LOGIN]: PANIC: Token not exist, but must\n")
        c.JSON(http.StatusInternalServerError, gin.H{})
        return
    }

    c.SetCookie("jwt", userDB.Token, 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{})
}

func Delete(c *gin.Context) {
    cookie, err := c.Cookie("jwt")

    if err != nil {
        log.Printf("[/user/delete]: no cookie: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    userDB, ok := postAuthCredentials(c)

    if !ok {
        log.Printf("[/user/delete]: no credentials\n")
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    if userDB.Token != cookie {
        log.Printf("[/user/delete]: no credentials\n")
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    model.DeleteUser(userDB.Login)
    c.JSON(http.StatusOK, gin.H{})
}
