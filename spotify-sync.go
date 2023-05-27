package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo"
)

const (
	authoriseEndpoint string = "https://accounts.spotify.com/authorize"
	scopes string = "playlist-modify-private playlist-modify-public"
	redirectUri string = "http://localhost:9000/callback"
	stateVal string = "miguel"
)

var (
	clientId string = getenv("SPOTIFY_API_CLIENT_ID", "fakeid")
)

func main() {
	authCodeChan := make(chan string)

	e := echo.New()
	e.GET("/login", func(c echo.Context) error {
		loginUrl := getLoginUrl()
		fmt.Println("loginUrl = " + loginUrl)
		return c.Redirect(301, loginUrl)
	})

	e.GET("/callback", func(c echo.Context) error {
		code := c.QueryParams().Get("code")
		fmt.Println("code = " + code)
		authCodeChan <- code
		return c.String(http.StatusOK, "Got auth code " + code)
	})

	go func() {
		if err := e.Start(":9000"); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	go hitLoginUrl()

	authCode := <- authCodeChan
	fmt.Println("Got auth code " + authCode)

	time.Sleep(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func getLoginUrl() string {
	req, err := http.NewRequest(http.MethodGet, authoriseEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", clientId)
	q.Add("redirect_uri", redirectUri)
	q.Add("state", stateVal)
	q.Add("scopes", scopes)

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func hitLoginUrl() {
	chromeCtx, cancelChrome := chromedp.NewContext(context.Background())
	defer cancelChrome()

	_, err := chromedp.RunResponse(chromeCtx, emulation.SetUserAgentOverride("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"), chromedp.Navigate("http://localhost:9000/login"))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Finished hitting login url")
	}
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}