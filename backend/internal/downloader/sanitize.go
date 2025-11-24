package downloader

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var formatRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func SanitizeURL(inputURL string, allowedDomains []string) (string, error) {
	if len(inputURL) > 2048 {
		return "", fmt.Errorf("URL too long")
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid protocol")
	}

	hostname := strings.TrimPrefix(parsedURL.Hostname(), "www.")
	allowed := false
	for _, domain := range allowedDomains {
		if hostname == domain || strings.HasSuffix(hostname, "."+domain) {
			allowed = true
			break
		}
	}

	if !allowed {
		return "", fmt.Errorf("domain not allowed")
	}

	return parsedURL.String(), nil
}

func SanitizeFormat(format string) string {
	if format == "" || !formatRegex.MatchString(format) {
		return "best"
	}
	if len(format) > 20 {
		return "best"
	}
	return format
}
