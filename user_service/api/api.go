package api

import (
	"log"
	"net/http"

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
    user, ok := postAuthCredentials(c)

    if !ok || user.Login == "" || user.Password == "" {
        log.Printf("[REGISTER]: Error on get user credentials\n")
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": "gophermart",
        "sub": user.Login,
    })
    str, err := tok.SignedString(secretKey)

    if err != nil {
        log.Printf("[REGISTER]: Error on signing token: %q\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{})
        return
    }

    _, err = model.NewUser(user.Login, user.Password, str)

    if err != nil {
        if err == model.ErrUserAlreadyExist {
            log.Printf("[REGISTER]: already exist %q\n", user.Login)
            c.JSON(http.StatusConflict, gin.H{})
            return
        }

        log.Printf("[REGISTER]: Unknown error: %q\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{})
        return
    }

    c.SetCookie("jwt", str, 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{})
}

// userDBbyCreadentials - Returns user from DB by login
func userDBbyCreadentials(login, password string) (*model.User, bool) {
    user, err := model.UserByLogin(login)

    if err != nil || user.Password != password {
        log.Printf("[LOGIN]: Error: %q: %q, is password good? %t\n", err, login, user.Password == password)
        return nil, false
    }

    return user, true
}


// postAuthCredentials - returns User type from post body or false
func postAuthCredentials(c *gin.Context) (*model.User, bool) {
    var user model.User
    err := c.ShouldBindBodyWithJSON(&user)

    if err != nil || user.Login == "" || user.Password == "" {
        log.Printf("[LOGIN]: Error on binding: %v\n", err)
        return nil, false
    }

    return &user, true
}

func Login(c *gin.Context) {
    var user *model.User
    var ok bool

    user, ok = postAuthCredentials(c)

    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }

    userDB, ok := userDBbyCreadentials(user.Login, user.Password)

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

    user, ok := postAuthCredentials(c)

    if !ok {
        log.Printf("[/user/delete]: no credentials\n")
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    userDB, ok := userDBbyCreadentials(user.Login, user.Password)

    if userDB.Token != cookie {
        log.Printf("[/user/delete]: no credentials\n")
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    err = model.DeleteUser(userDB)

    if err != nil {
        log.Printf("[/user/delete]: error from model: %q\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{})
        return
    }

    c.SetCookie("jwt", "", -1, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{})
}
