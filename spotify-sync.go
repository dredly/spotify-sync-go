package main

import (
	"dredly/spotify-sync/apiclient"
	"dredly/spotify-sync/browserautomation"
	"dredly/spotify-sync/echoserver"
	"fmt"
	"log"
	"os"
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
	go browserautomation.AutoLogin()

	authCode := <-authCodeChan
	echoserver.GracefulShutdown(e)
	apiclient.GetAccessToken(authCode)
}