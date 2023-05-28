package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo"
)

const (
	authoriseEndpoint string = "https://accounts.spotify.com/authorize"
	tokenEndpoint     string = "https://accounts.spotify.com/api/token"
	scopes            string = "playlist-modify-private playlist-modify-public"
	redirectUri       string = "http://localhost:9000/callback"
	stateVal          string = "miguel"
)

var (
	clientId        string = getenv("SPOTIFY_API_CLIENT_ID", "fakeid")
	clientSecret    string = getenv("SPOTIFY_API_CLIENT_SECRET", "fakesecret")
	spotifyUser     string = getenv("SPOTIFY_USERNAME", "fakeusername")
	spotifyPassword string = getenv("SPOTIFY_PASSWORD", "fakepassword")
)

func main() {
	playlistIds := os.Args[1:]
	if len(playlistIds) == 0 {
		// TODO: Add usage info here
		log.Fatal("No playlist ids Provided")
	}
	if len(playlistIds) % 2 != 0 {
		log.Fatal("Each source playlist must have a destionation")
	}

	fmt.Printf("Running spotify-sync with playlist ids %v", playlistIds)

	authCodeChan := make(chan string)

	e := echo.New()
	e.GET("/login", func(c echo.Context) error {
		loginUrl := getLoginUrl()
		return c.Redirect(301, loginUrl)
	})

	e.GET("/callback", func(c echo.Context) error {
		code := c.QueryParams().Get("code")
		authCodeChan <- code
		return c.String(http.StatusOK, "Got auth code " + code)
	})

	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	go autoLogin()

	authCode := <-authCodeChan
	echoServerGracefulShutdown(e)
	getAccessToken(authCode)
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

func autoLogin() {
	fmt.Println("Attempting autologin")
	chromeCtx, cancelChrome := chromedp.NewContext(context.Background())
	defer cancelChrome()

	err := chromedp.Run(chromeCtx, fillInLoginForm())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Finished hitting login url")
	}
}

func getAccessToken(code string) {
	fmt.Println("Getting access token")
	
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("redirect_uri", redirectUri)

	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientId, clientSecret)

	c := http.Client{Timeout: 5 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Body : %s", resp_body)
}

func fillInLoginForm() chromedp.Tasks {
	loginButtonQuery := "#login-button > span.ButtonInner-sc-14ud5tc-0.cJdEzG.encore-bright-accent-set > span"
	return chromedp.Tasks{
		chromedp.Navigate("http://localhost:9000/login"),
		chromedp.WaitVisible("login-username", chromedp.ByID),
		chromedp.WaitVisible("login-password", chromedp.ByID),
		chromedp.WaitVisible(loginButtonQuery, chromedp.ByQuery),
		chromedp.SendKeys("login-username", spotifyUser, chromedp.ByID),
		chromedp.SendKeys("login-password", spotifyPassword, chromedp.ByID),
		chromedp.Click(loginButtonQuery, chromedp.ByQuery),
		chromedp.WaitNotVisible("login-username", chromedp.ByID),
	}
}

func echoServerGracefulShutdown(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
