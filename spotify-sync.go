package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	authoriseEndpoint string = "https://accounts.spotify.com/authorize"
	scopes string = "playlist-modify-private playlist-modify-public"
)

var (
	client_id string = getenv("SPOTIFY_API_CLIENT_ID", "fakeid")
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/kill", func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		e.Shutdown(context.Background())
		return c.String(http.StatusOK, "Hello, World!")
	})

	authorise()

	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Fatal(err)
		}
	}()
	fmt.Println("Server running")

	time.Sleep(30 * time.Second)

	fmt.Println("Server timed out")
}

func authorise() {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, authoriseEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", client_id)
	q.Add("redirect_uri", "miguel")
	q.Add("state", "whatever")
	q.Add("scopes", scopes)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}