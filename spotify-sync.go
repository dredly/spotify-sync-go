package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dredly/spotify-sync-go/apiclient"
	"github.com/dredly/spotify-sync-go/browserautomation"
	"github.com/dredly/spotify-sync-go/cli"
	"github.com/dredly/spotify-sync-go/echoserver"
)

func main() {
	fmt.Printf("Spotify-Sync -- %v local time\n", time.Now().Format("2006-01-02 15:04:05"))
	playlistIdPairs := cli.GetPlaylistIdPairs()

	c := *apiclient.NewHttpClient()

	rt := apiclient.GetRefreshTokenFromFileIfPresent()
	var token string

	if rt != "" {
		fmt.Println("Using refresh token")
		token = apiclient.RefreshAccessToken(c, rt)
	} else {
		token = getTokenThroughAutoLogin(c)
	}

	for _, pip := range playlistIdPairs {
		apiclient.Sync(c, token, pip)
	}
}

func getTokenThroughAutoLogin(c http.Client) string {
	authCodeChan := make(chan string)
	e := echoserver.SpinUpTempServer(authCodeChan)

	go browserautomation.AutoLogin()
	
	var authCode string
	select {
	case authCode = <-authCodeChan:
	case <-time.After(30 * time.Second):
		echoserver.GracefulShutdown(e)
		log.Fatal("Timed out waiting for auth code")
	}

	echoserver.GracefulShutdown(e)
	return apiclient.GetAccessToken(c, authCode)
}
