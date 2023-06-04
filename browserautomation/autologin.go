package browserautomation

import (
	"context"
	"fmt"
	"log"

	"github.com/dredly/spotify-sync-go/utils"

	"github.com/chromedp/chromedp"
)

const userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36"

var (
	spotifyUser     string = utils.GetEnvWithFallback("SPOTIFY_USERNAME", "fakeusername")
	spotifyPassword string = utils.GetEnvWithFallback("SPOTIFY_PASSWORD", "fakepassword")
)

func AutoLogin() {
	fmt.Println("Attempting autologin for spotify user ", spotifyUser)

	parentCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", true), chromedp.UserAgent(userAgent))...)
	defer cancel()

	chromeCtx, cancelChrome := chromedp.NewContext(parentCtx)
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
		chromedp.EmulateViewport(1920, 1080),
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
