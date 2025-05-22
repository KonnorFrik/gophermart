package api

import (
	"gophermart/model"
	"gophermart/model/models"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func createToken(user *models.User) (string, error) {
    userIDstr := strconv.FormatInt(user.ID, 10)
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": "gophermart",
        "sub": userIDstr,
        "iat": time.Now().Unix(), // issued at
        // "exp": // expiration time
    })
    str, err := tok.SignedString(jwtSecretKey)

    log.Printf("[createToken]: Create token with 'sub'=%q\n", userIDstr)
    return str, err
}

// postAuthCredentials - returns User type from post body.
func postAuthCredentials(c *gin.Context) (*models.User, bool) {
    var user models.User
    err := c.ShouldBindBodyWithJSON(&user)

    if err != nil || user.Login == "" || user.Password == "" {
        log.Printf("[/login]: Error on binding: %v\n", err)
        return nil, false
    }

    return &user, true
}

// userDBbyCreadentials - Returns user from DB by login
func userDBbyCreadentials(login, password string) (*models.User, bool) {
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

func Int64(value string) (int64, error) {
    value64, err := strconv.ParseInt(value, 10, strconv.IntSize)

    if err != nil {
        return 0, err
    }

    return value64, nil
}
