package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"unicode"

	"gophermart/model"
	"gophermart/model/models"

	"github.com/gin-gonic/gin"
)

const (
    jwtTokenCookieName = "token"
)

var (
    jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func init() {
    if len(jwtSecretKey) == 0 {
        panic("jwt secret key is missing")
    }
}

// Register - create a new user and token
func UserRegister(c *gin.Context) {
    user, ok := postAuthCredentials(c)

    if !ok || user.Login == "" || user.Password == "" {
        log.Printf("[POST /register]: Error on get user credentials\n")
        c.Status(http.StatusBadRequest)
        return
    }

    var err error
    userDB, err := model.NewUser(user.Login, user.Password)

    // if errors.Is(err, model.ErrUserAlreadyExist) {
    //     log.Printf("[/register]: already exist %q\n", userDB.Login)
    //     c.Status(http.StatusConflict)
    //     return
    // }
    if errors.Is(err, model.ErrDataBaseNotConnected) {
        log.Printf("[POST /register]: db not connected\n")
        c.Status(http.StatusInternalServerError)
        return
    }

    if err != nil {
        log.Printf("[POST /register]: Unknown error: %q\n", err)
        c.Status(http.StatusInternalServerError)
        return
    }

    strToken, err := createToken(userDB.ID)

    if err != nil {
        log.Printf("[POST /register]: Error on create token: %q\n", err)
        c.Status(http.StatusInternalServerError)
        return
    }

    uid := strconv.FormatInt(userDB.ID, 10)
    c.SetCookie("uid", uid, 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{
        "token": strToken,
    })
}

func UserLogin(c *gin.Context) {
    user, ok := postAuthCredentials(c)

    if !ok {
        c.Status(http.StatusUnauthorized)
        return
    }

    userDB, ok := userDBbyCreadentials(user.Login, user.Password)

    if !ok {
        c.Status(http.StatusUnauthorized)
        return
    }

    strToken, err := createToken(userDB.ID)

    if err != nil {
        log.Printf("[POST /login]: Error on token creation: %q\n", err)
        c.Status(http.StatusInternalServerError)
        return
    }

    uid := strconv.FormatInt(userDB.ID, 10)
    c.SetCookie("uid", uid, 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{
        "token": strToken,
    })
}

func UserDelete(c *gin.Context) {
    cookieUID, err := c.Cookie("uid")

    if err != nil {
        log.Printf("[DELETE /user]: missed cookie\n")
        c.Status(http.StatusBadRequest)
        return
    }

    userID, err := ToInt64(cookieUID)

    if err != nil {
        log.Printf("[DELETE /user]: invalid uid from cookie: %q\n", err)
        c.Status(http.StatusUnauthorized)
        return
    }

    user, ok := postAuthCredentials(c)

    if !ok {
        log.Printf("[DELETE /user]: no credentials\n")
        c.Status(http.StatusBadRequest)
        return
    }

    // validate creadentials
    userDB, ok := userDBbyCreadentials(user.Login, user.Password)

    if !ok {
        log.Printf("[DELETE /user]: Not found in DB: %q\n", user.Login)
        c.Status(http.StatusUnauthorized)
        return
    }

    if userDB.ID != userID {
        log.Printf("[DELETE /user]: Conflict: different id: DB(%d) != cookie(%d)\n", userDB.ID, userID)
        c.Status(http.StatusBadRequest)
        return
    }

    err = model.DeleteUser(userID)

    if err != nil {
        log.Printf("[DELETE /user]: error from model.DeleteUser: %q\n", err)
        c.Status(http.StatusInternalServerError)
        return
    }

    c.Status(http.StatusOK)
}

// Create a new order
// accept text/plain
func NewOrder(c *gin.Context) {
    cookieUID, err := c.Cookie("uid")

    if err != nil {
        log.Printf("[POST /orders]: missed cookie\n")
        c.Status(http.StatusUnauthorized)
        return
    }

    userID, err := ToInt64(cookieUID)

    if err != nil {
        log.Printf("[POST /orders]: invalid uid from cookie: %q\n", err)
        c.Status(http.StatusUnauthorized)
        return
    }

    bodyBytes, err := io.ReadAll(c.Request.Body)

    if err != nil {
        log.Printf("[POST /orders]: Error on read data: %q\n", err)
        c.Status(http.StatusBadRequest)
        return
    }

    if len(bodyBytes) == 0 {
        log.Printf("[POST /orders]: Empty Body\n")
        c.Status(http.StatusBadRequest)
        return
    }

    orderString := string(bodyBytes)

    for _, r := range orderString {
        if !unicode.IsDigit(r) {
            log.Printf("[POST /orders]: Invalid input data: %q\n", orderString)
            c.Status(http.StatusBadRequest)
            return
        }
    }
    
    if !validByLUHN(orderString) {
        log.Printf("[POST /orders]: Invalid by LUHN input\n")
        c.Status(http.StatusUnprocessableEntity)
        return
    }

    order := models.Order{
        Number: orderString,
    }
    err = model.NewOrder(&order, &models.User{ID: userID})
    /*
    TODO:
    error can be like: duplicate - violates contraint
    need to catch this and return code 200
    */

    switch {
    case errors.Is(err, model.ErrDataBaseNotConnected):
        c.Status(http.StatusInternalServerError)
    case errors.Is(err, model.ErrOrderNothingToCreate):
        c.Status(http.StatusInternalServerError)
    case err != nil:
        log.Printf("[POST /orders]: Got err: %q\n", err)
        c.Status(http.StatusInternalServerError)

    default:
        c.Status(http.StatusAccepted)
    }

    // c.Status(http.StatusAccepted)
}

func AllOrders(c *gin.Context) {
    cookieUID, err := c.Cookie("uid")

    if err != nil {
        log.Printf("[GET /orders]: missed cookie\n")
        c.Status(http.StatusUnauthorized)
        return
    }

    userID, err := ToInt64(cookieUID)

    if err != nil {
        log.Printf("[GET /orders]: invalid uid from cookie: %q\n", err)
        c.Status(http.StatusUnauthorized)
        return
    }

    orders, err := model.OrdersRelated(&models.User{ID: userID})

    if err != nil {
        log.Printf("[GET /orders]: Error on fetch data: %q\n", err)
        c.Status(http.StatusInternalServerError)
        return
    }

    c.JSON(http.StatusOK, orders)
}

