package main

import (
	"dredly/spotify-sync/apiclient"
	"dredly/spotify-sync/browserautomation"
	"dredly/spotify-sync/cli"
	"dredly/spotify-sync/echoserver"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Printf("Spotify-Sync -- %v local time\n", time.Now().Format("2006-01-02 15:04:05"))
	playlistIdPairs := cli.GetPlaylistIdPairs()

	authCodeChan := make(chan string)
	e := echoserver.SpinUpTempServer(authCodeChan)

	go browserautomation.AutoLogin()
	c := *apiclient.NewHttpClient()

	var authCode string
	select {
	case authCode = <-authCodeChan:
	case <-time.After(30 * time.Second):
		echoserver.GracefulShutdown(e)
		log.Fatal("Timed out waiting for auth code")
	}
	
	echoserver.GracefulShutdown(e)
	t := apiclient.GetAccessToken(c, authCode)
	for _, pip := range playlistIdPairs {
		apiclient.Sync(c, t, pip)
	}
}
