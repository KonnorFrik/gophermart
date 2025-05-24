package api

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createToken(userID int64) (string, error) {
    userIDstr := strconv.FormatInt(userID, 10)
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": "gophermart",
        "sub": userIDstr,
        "iat": time.Now().Unix(), // issued at
    })
    str, err := tok.SignedString(jwtSecretKey)

    // log.Printf("[createToken]: Create token with 'sub'=%q\n", userIDstr)
    return str, err
}

func validByLUHN(numbers string) bool {
    return len(numbers) > 0
}

func ToInt64(value string) (int64, error) {
    value64, err := strconv.ParseInt(value, 10, strconv.IntSize)

    if err != nil {
        return 0, err
    }

    return value64, nil
}
