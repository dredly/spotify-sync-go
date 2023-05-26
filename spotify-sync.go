package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
		return c.Redirect(http.StatusPermanentRedirect, loginUrl)
	})

	e.GET("/callback", func(c echo.Context) error {
		code := c.QueryParams().Get("code")
		authCodeChan <- code
		return c.String(http.StatusOK, "Got auth code " + code)
	})

	e.GET("/kill", func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		e.Shutdown(context.Background())
		return c.String(http.StatusOK, "Hello, World!")
	})

	go func() {
		if err := e.Start(":9000"); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	authCode := <- authCodeChan
	fmt.Println("Got auth code " + authCode)

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

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}