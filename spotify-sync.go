package main

import (
	"dredly/spotify-sync/apiclient"
	"dredly/spotify-sync/browserautomation"
	"dredly/spotify-sync/echoserver"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	playlistIds := os.Args[1:]
	if len(playlistIds) == 0 {
		// TODO: Add usage info here
		log.Fatal("No playlist ids Provided")
	}
	if len(playlistIds)%2 != 0 {
		log.Fatal("Each source playlist must have a destionation")
	}

	fmt.Printf("Running spotify-sync with playlist ids %v", playlistIds)

	authCodeChan := make(chan string)
	e := echoserver.SpinUpTempServer(authCodeChan)

	// This is to make sure the eechoserver is up and running. Temporary solution
	time.Sleep(1 * time.Second) 

	go browserautomation.AutoLogin()
	c := *apiclient.NewHttpClient()

	authCode := <-authCodeChan
	echoserver.GracefulShutdown(e)
	t := apiclient.GetAccessToken(c, authCode)
	apiclient.GetDestinationTrackUris(c, t, "03whiAjg4TdJtDikG6wZIa?si=2b64819f13ce4c60")
}
