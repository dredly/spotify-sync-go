package apiclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dredly/spotify-sync-go/utils"
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

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Returns the access token and also saves the refresh token to a file
func GetAccessToken(c http.Client, code string) string {
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

	saveToken(atr.RefreshToken)

	return atr.AccessToken
}

func RefreshAccessToken(c http.Client, refreshToken string) string {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)
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

	var rtr refreshTokenResponse

	err = json.Unmarshal(respBody, &rtr)
	if err != nil {
		log.Fatal(err)
	}

	return rtr.AccessToken
}