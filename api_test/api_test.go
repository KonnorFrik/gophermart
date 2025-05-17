package apitest_test

import (
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

func TestHello(t *testing.T) {
    response, err := client.R().Get(address + "/hello")

    if err != nil {
        t.Fatalf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = 200

    if statusCode != wantStatus {
        t.Fatalf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }

    var body string = string(response.Body())
    var wantBody string = `{"cookie":"NotSet","status":"OK"}`

    if body != wantBody {
        t.Fatalf("Got = %q, Want = %q\n", body, wantBody)
    }
}

func TestCRUD(t *testing.T) {
    cooks := testUserRegister(t)
    cooks = testUserLogin(t, cooks)
    testUserDelete(t, cooks)
}

func testUserRegister(t *testing.T) []*http.Cookie {
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

    return response.Cookies()
}

func testUserLogin(t *testing.T, cookies []*http.Cookie) []*http.Cookie {
    req := client.R().SetBody(
        map[string]string{
            "login": "user",
            "password": "user",
        },
    ).SetCookies(cookies)

    response, err := req.Post(address + "/api/user/login")

    if err != nil {
        t.Errorf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = 200

    if statusCode != wantStatus {
        t.Errorf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }

    return response.Cookies()
}

func testUserDelete(t *testing.T, cookies []*http.Cookie) {
    req := client.R().SetBody(
        map[string]string{
            "login": "user",
            "password": "user",
        },
    ).SetCookies(cookies)

    response, err := req.Post(address + "/api/user/delete")

    if err != nil {
        t.Errorf("Got = %v\n", err)
    }

    var statusCode int = response.StatusCode()
    var wantStatus int = 200

    if statusCode != wantStatus {
        t.Errorf("Got = %d, Want = %d\n", statusCode, wantStatus)
    }
}
