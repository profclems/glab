package utils

import (
	"fmt"
	"glab/internal/config"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/gookit/color"
	"glab/internal/browser"
	"glab/internal/run"
)

// OpenInBrowser opens the url in a web browser based on OS and $BROWSER environment variable
func OpenInBrowser(url string) error {
	browseCmd, err := browser.Command(url)
	if err != nil {
		return err
	}
	return run.PrepareCmd(browseCmd).Run()
}

func RenderMarkdown(text string) (string, error) {
	// Glamour rendering preserves carriage return characters in code blocks, but
	// we need to ensure that no such characters are present in the output.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	renderStyle := glamour.WithStandardStyle("dark")
	if config.GetEnv("GLAMOUR_STYLE") != "" {
		renderStyle = glamour.WithEnvironmentConfig()
	}

	tr, err := glamour.NewTermRenderer(
		renderStyle,
	)
	if err != nil {
		return "", err
	}

	return tr.Render(text)
}

func Pluralize(num int, thing string) string {
	if num == 1 {
		return fmt.Sprintf("%d %s", num, thing)
	}
	return fmt.Sprintf("%d %ss", num, thing)
}

func fmtDuration(amount int, unit string) string {
	return fmt.Sprintf("about %s ago", Pluralize(amount, unit))
}

func PrettyTimeAgo(ago time.Duration) string {
	if ago < time.Minute {
		return "less than a minute ago"
	}
	if ago < time.Hour {
		return fmtDuration(int(ago.Minutes()), "minute")
	}
	if ago < 24*time.Hour {
		return fmtDuration(int(ago.Hours()), "hour")
	}
	if ago < 30*24*time.Hour {
		return fmtDuration(int(ago.Hours())/24, "day")
	}
	if ago < 365*24*time.Hour {
		return fmtDuration(int(ago.Hours())/24/30, "month")
	}

	return fmtDuration(int(ago.Hours()/24/365), "year")
}

func TimeToPrettyTimeAgo(d time.Time) string {
	now := time.Now()
	ago := now.Sub(d)
	return PrettyTimeAgo(ago)
}

func Humanize(s string) string {
	// Replaces - and _ with spaces.
	replace := "_-"
	h := func(r rune) rune {
		if strings.ContainsRune(replace, r) {
			return ' '
		}
		return r
	}

	return strings.Map(h, s)
}

func IsURL(s string) bool {
	return strings.HasPrefix(s, "http:/") || strings.HasPrefix(s, "https:/")
}

func DisplayURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	return u.Hostname() + u.Path
}

func GreenCheck() string {
	return color.Green.Sprintf("✓")
}