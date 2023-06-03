package apiclient

import (
	"dredly/spotify-sync/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	redirectUri   string = "http://localhost:9000/callback"
	tokenEndpoint string = "https://accounts.spotify.com/api/token"
)

var (
	clientId     string = utils.GetEnvWithFallback("SPOTIFY_API_CLIENT_ID", "fakeid")
	clientSecret string = utils.GetEnvWithFallback("SPOTIFY_API_CLIENT_SECRET", "fakesecret")
)

type accessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetAccessToken(c http.Client, code string) string {
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

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var atr accessTokenResponse
	err = json.Unmarshal(respBody, &atr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(atr.AccessToken)
	return atr.AccessToken
}
