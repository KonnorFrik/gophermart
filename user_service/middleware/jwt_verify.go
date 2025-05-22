package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


var (
    jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

    ErrInvalidToken = errors.New("invalid token")
)

func init() {
    if len(jwtSecretKey) == 0 {
        panic("jwt secret key is missing")
    }
}

func JWTAuthenticate(c *gin.Context) {
    authToken := c.GetHeader("Authorization")

    if authToken == "" {
        log.Printf("[jwtAuthMDLWR]: token is missed\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    var err error
    tok, err := verifyJWT(authToken)

    if err != nil {
        log.Printf("[jwtAuthMDLWR]: failed token verification\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    mapClaims, ok := tok.Claims.(jwt.MapClaims)

    if !ok {
        log.Printf("[jwtAuthMDLWR]: Can't convert claims\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }
    
    userIDstr, err := mapClaims.GetSubject()

    if err != nil {
        log.Printf("[jwtAuthMDLWR]: Can't get subject claim: %q\n", err)
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    cookieUID, err := c.Cookie("uid")

    if cookieUID != userIDstr {
        log.Printf("[jwtAuthMDLWR]: Cookie id does not math jwt id\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    c.Next()
}

func verifyJWT(tokStr string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokStr, func(t *jwt.Token) (interface{}, error) {
        return jwtSecretKey, nil
    })

    if err != nil {
        return nil, err 
    }

    if !token.Valid {
        return nil, ErrInvalidToken
    }

    return token, nil
}

