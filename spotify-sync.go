package main

import (
	"dredly/spotify-sync/apiclient"
	"dredly/spotify-sync/browserautomation"
	"dredly/spotify-sync/cli"
	"dredly/spotify-sync/echoserver"
)

func main() {
	playlistIdPairs := cli.GetPlaylistIdPairs()

	authCodeChan := make(chan string)
	e := echoserver.SpinUpTempServer(authCodeChan)

	go browserautomation.AutoLogin()
	c := *apiclient.NewHttpClient()

	authCode := <-authCodeChan
	echoserver.GracefulShutdown(e)
	t := apiclient.GetAccessToken(c, authCode)
	for _, pip := range playlistIdPairs {
		apiclient.Sync(c, t, pip)
	}
}
