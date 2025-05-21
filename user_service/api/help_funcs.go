package api

import (
	"gophermart/model"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func createToken(user *model.User) (string, error) {
    userIDstr := strconv.FormatUint(uint64(user.ID), 10)
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": "gophermart",
        "sub": userIDstr,
        "iat": time.Now().Unix(), // issued at
        // "exp": // expiration time
    })
    str, err := tok.SignedString(jwtSecretKey)

    log.Printf("[createToken]: Create token with 'sub'=user ID= %q\n", userIDstr)
    return str, err
}

// postAuthCredentials - returns User type from post body.
func postAuthCredentials(c *gin.Context) (*model.User, bool) {
    var user model.User
    err := c.ShouldBindBodyWithJSON(&user)

    if err != nil || user.Login == "" || user.Password == "" {
        log.Printf("[/login]: Error on binding: %v\n", err)
        return nil, false
    }

    return &user, true
}

// userDBbyCreadentials - Returns user from DB by login
func userDBbyCreadentials(login, password string) (*model.User, bool) {
    user, err := model.UserByLogin(login)

    if err != nil {
        log.Printf("[LOGIN]: Error for login - %q: %q\n", login, err)
        return nil, false
    }

    // TODO: use crypted password and related validation
    // TODO: Create method .ValidCred for User type
    if user.Password != password {
        log.Printf("[LOGIN]: %q: wrong password\n", login)
        return nil, false
    }

    return user, true
}

func validByLUHN(numbers string) bool {
    return len(numbers) > 0
}

type cookieGetter interface {
    Cookie(string) (string, error)
}

func UintFromCookie(cg cookieGetter, key string) (uint, error) {
    value, err := cg.Cookie(key)

    if err != nil {
        return 0, err
    }

    value64, err := strconv.ParseUint(value, 10, strconv.IntSize)

    if err != nil {
        return 0, err
    }

    return uint(value64), nil
}
