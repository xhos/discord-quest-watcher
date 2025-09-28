package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func createBrowser() (*rod.Browser, error) {
	launcher := launcher.New().
		Headless(true).
		Bin("/usr/bin/chromium-browser").
		Set("no-sandbox").
		Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-web-security")

	browser := rod.New().ControlURL(launcher.MustLaunch()).MustConnect()
	return browser, nil
}

func authenticateWithToken(browser *rod.Browser, token string) error {
	page := browser.MustPage("https://discord.com/login").MustWaitLoad()

	// inject token
	_, err := page.Eval(fmt.Sprintf(`
		() => {
			const iframe = document.createElement('iframe');
			document.body.appendChild(iframe);
			iframe.contentWindow.localStorage.token = '"%s"';
			setTimeout(() => location.reload(), 2000);
		}
	`, token))

	if err != nil {
		return err
	}

	// wait for redirect
	for range 30 {
		if !strings.Contains(page.MustInfo().URL, "/login") {
			log.Info("authenticated successfully")
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("authentication timeout")
}
