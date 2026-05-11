package downloader

import (
	"os"
	"path/filepath"
	"testing"
)

func containsArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}

func TestBuildDownloadArgsSkipsMissingCookiesFile(t *testing.T) {
	missingCookies := filepath.Join(t.TempDir(), "missing-cookies.txt")
	d := New(t.TempDir(), missingCookies, "20M")

	args := d.buildDownloadArgs(
		"https://www.youtube.com/watch?v=jNQXAC9IVRw",
		"best",
		filepath.Join(t.TempDir(), "%(title)s.%(ext)s"),
		0,
	)

	if containsArg(args, "--cookies") {
		t.Fatalf("buildDownloadArgs included --cookies for missing cookies file: %v", args)
	}
	if !containsArg(args, "youtube:player_client=web_safari") {
		t.Fatalf("buildDownloadArgs should use cookieless YouTube player args when cookies are missing: %v", args)
	}
}

func TestBuildDownloadArgsIncludesUsableCookiesFile(t *testing.T) {
	cookies := filepath.Join(t.TempDir(), "cookies.txt")
	if err := os.WriteFile(cookies, []byte("# Netscape HTTP Cookie File\n"), 0600); err != nil {
		t.Fatal(err)
	}
	d := New(t.TempDir(), cookies, "20M")

	args := d.buildDownloadArgs(
		"https://www.youtube.com/watch?v=jNQXAC9IVRw",
		"best",
		filepath.Join(t.TempDir(), "%(title)s.%(ext)s"),
		0,
	)

	if !containsArg(args, "--cookies") || !containsArg(args, cookies) {
		t.Fatalf("buildDownloadArgs did not include usable cookies file: %v", args)
	}
	if !containsArg(args, "youtube:player_client=default,web_safari") {
		t.Fatalf("buildDownloadArgs should use cookie-aware YouTube player args when cookies are usable: %v", args)
	}
}
