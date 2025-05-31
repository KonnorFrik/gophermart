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

func JWTAuthenticate(c *gin.Context) {
    if len(jwtSecretKey) == 0 {
        log.Println("[MDWR/jwtAuth]: Token not loaded")
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    authToken := c.GetHeader("Authorization")

    if authToken == "" {
        log.Printf("[MDWR/jwtAuth]: token is missed\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    var err error
    tok, err := verifyJWT(authToken)

    if err != nil {
        log.Printf("[MDWR/jwtAuth]: failed token verification: %q\n", err)
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    mapClaims, ok := tok.Claims.(jwt.MapClaims)

    if !ok {
        log.Printf("[MDWR/jwtAuth]: Can't convert claims\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }
    
    userIDstr, err := mapClaims.GetSubject()

    if err != nil {
        log.Printf("[MDWR/jwtAuth]: Can't get subject claim: %q\n", err)
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    cookieUID, err := c.Cookie("uid")

    if cookieUID != userIDstr {
        log.Printf("[MDWR/jwtAuth]: Cookie id does not math jwt id\n")
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    c.Next()
}

// verifyJWT - Parse a jwt token string for verification.
// Returns a token type
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

