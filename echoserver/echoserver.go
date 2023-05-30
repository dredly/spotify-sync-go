package echoserver

import (
	"context"
	"dredly/spotify-sync/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

const (
	authoriseEndpoint string = "https://accounts.spotify.com/authorize"
	tokenEndpoint     string = "https://accounts.spotify.com/api/token"
	scopes            string = "playlist-modify-private playlist-modify-public"
	redirectUri       string = "http://localhost:9000/callback"
	stateVal          string = "miguel"
)

var clientId string = utils.GetEnvWithFallback("SPOTIFY_API_CLIENT_ID", "fakeid")

func SpinUpTempServer(authCodeChan chan string) *echo.Echo {
	e := echo.New()
	e.GET("/login", func(c echo.Context) error {
		fmt.Println("Hit /login route")
		loginUrl := getLoginUrl()
		return c.Redirect(301, loginUrl)
	})

	e.GET("/callback", func(c echo.Context) error {
		fmt.Println("Hit /callback route")
		code := c.QueryParams().Get("code")
		authCodeChan <- code
		return c.String(http.StatusOK, "Got auth code " + code)
	})

	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	return e
}

func GracefulShutdown(e *echo.Echo) {
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