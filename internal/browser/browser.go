package browser

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

//go:embed inject-token.js
var injectTokenScript string

func CreateBrowser() (*rod.Browser, error) {
	return rod.New().ControlURL(
		launcher.New().
			Headless(true).
			Bin("/usr/bin/chromium-browser").
			Set("no-sandbox").
			Set("disable-gpu").
			Set("disable-dev-shm-usage").
			Set("disable-web-security").
			MustLaunch()).MustConnect(), nil
}

func AuthenticateWithToken(browser *rod.Browser, token string) error {
	page := browser.MustPage("https://discord.com/login").MustWaitLoad()

	// inject token
	script := strings.Replace(injectTokenScript, "__TOKEN__", token, 1)
	if _, err := page.Eval("() => {" + script + "}"); err != nil {
		return err
	}

	// wait for redirect
	for range 30 {
		if !strings.Contains(page.MustInfo().URL, "/login") {
			log.Println("authenticated successfully")
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("authentication timeout")
}
