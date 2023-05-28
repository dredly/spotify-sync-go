package browserautomation

import (
	"context"
	"dredly/spotify-sync/utils"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

var (
	spotifyUser     string = utils.GetEnvWithFallback("SPOTIFY_USERNAME", "fakeusername")
	spotifyPassword string = utils.GetEnvWithFallback("SPOTIFY_PASSWORD", "fakepassword")
)

func AutoLogin() {
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