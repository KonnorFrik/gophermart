package apitest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
)

const (
    proto = "http"
    ip = "localhost"
    port = 8080

    userLogin = "user"
    userPassw = "user"
)

var (
    client = resty.New()
    address = fmt.Sprintf("%s://%s:%d", proto, ip, port)
)

// funcs types
// producers : only return args (token, cookies)
// middlers : consume args and return data (token, cookies)
// full-consumers : only accept args (token cookies)
type chainFunc func(*testing.T, string, []*http.Cookie) (string, []*http.Cookie)

// Run testUserRegister for get token
// pass token to all given funcs and overwrite it with return value (may be same)
// at end call testUserDelete
func runChain(t *testing.T, chain ...chainFunc) {
    token, cookies := testUserRegister(t)

    for _, f := range chain {
        token, cookies = f(t, token, cookies)
    }

    testUserDelete(t, token, cookies)
}

func TestUserCRUD(t *testing.T) {
    testUserRegister(t)
    token, cookies := testUserLogin(t)
    testUserDelete(t, token, cookies)
}

// return a token string after registration
func testUserRegister(t *testing.T) (string, []*http.Cookie) {
    req := client.R()
    req.SetBody(map[string]string{
        "login": userLogin,
        "password": userPassw,
    })
    response, err := req.Post(address + "/api/user/register")

    if err != nil {
        t.Errorf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = 200

    if statusCode != wantStatus && statusCode != 409 {
        t.Errorf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }

    body := response.Body()

    if len(body) == 0 {
        t.Errorf("Got no body = len=%d\n", len(body))
    }

    var result = map[string]string{}

    if err := json.Unmarshal(body, &result); err != nil {
        t.Errorf("Can't Unmarshal = %q - %q\n", string(body), err)
    }

    token, ok := result["token"]

    if !ok {
        t.Errorf("Wrong body struct %v\n", result)
    }

    return token, response.Cookies()
}

func testUserLogin(t *testing.T) (string, []*http.Cookie) {
    req := client.R().SetBody(
        map[string]string{
            "login": userLogin,
            "password": userPassw,
        },
    )

    response, err := req.Post(address + "/api/user/login")

    if err != nil {
        t.Errorf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = 200

    if statusCode != wantStatus {
        t.Errorf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }

    body := response.Body()

    if len(body) == 0 {
        t.Errorf("Got no body = len=%d\n", len(body))
    }

    var result = map[string]string{}

    if err := json.Unmarshal(body, &result); err != nil {
        t.Errorf("Can't Unmarshal = %q - %q\n", string(body), err)
    }

    token, ok := result["token"]

    if !ok {
        t.Errorf("Wrong body struct %v\n", result)
    }

    return token, response.Cookies()
}

func testUserDelete(t *testing.T, token string, cookies []*http.Cookie) {
    req := client.R().SetBody(
        map[string]string{
            "login": userLogin,
            "password": userPassw,
        },
    ).SetHeader("Authorization", token).SetCookies(cookies)

    response, err := req.Delete(address + "/api/user/delete")

    if err != nil {
        t.Errorf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = 200

    if statusCode != wantStatus {
        t.Errorf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }
}

func TestOrderCRUD(t *testing.T) {
    runChain(
        t,
        testCreateOrder,
    )
}

func testCreateOrder(t *testing.T, token string, cookies []*http.Cookie) (string, []*http.Cookie) {
    req := client.R().SetHeader("Authorization", token).SetCookies(cookies)
    req = req.SetBody("1234")
    response, err := req.Post(address + "/api/user/orders")

    if err != nil {
        t.Errorf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = http.StatusAccepted

    if statusCode != wantStatus {
        t.Errorf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }

    return token, response.Cookies()
}
