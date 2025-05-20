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
)

var (
    client = resty.New()
    address = fmt.Sprintf("%s://%s:%d", proto, ip, port)
)

type chainFunc func(*testing.T, string) string

// Run testUserRegister for get token
// pass token to all given funcs and overwrite it with return value (may be same)
// at end call testUserDelete
func runChain(t *testing.T, chain ...chainFunc) {
    token := testUserRegister(t)

    for _, f := range chain {
        token = f(t, token)
    }

    testUserDelete(t, token)
}

func TestUserCRUD(t *testing.T) {
    token := testUserRegister(t)
    token = testUserLogin(t, token)
    testUserDelete(t, token)
}

// return a token string after registration
func testUserRegister(t *testing.T) string {
    req := client.R()
    req.SetBody(map[string]string{
        "login": "user",
        "password": "user",
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

    return token
}

func testUserLogin(t *testing.T, token string) string {
    req := client.R().SetBody(
        map[string]string{
            "login": "user",
            "password": "user",
        },
    ).SetHeader("Authorization", token)

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

    return token
}

func testUserDelete(t *testing.T, token string) {
    req := client.R().SetBody(
        map[string]string{
            "login": "user",
            "password": "user",
        },
    ).SetHeader("Authorization", token)

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
        testUserLogin,
        testCreateOrder,
    )
}

func testCreateOrder(t *testing.T, token string) string {
    req := client.R().SetHeader("Authorization", token)
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

    return token
}
