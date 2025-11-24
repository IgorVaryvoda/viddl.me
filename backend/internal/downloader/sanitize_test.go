package downloader

import (
	"testing"
)

var testAllowedDomains = []string{
	"youtube.com",
	"youtu.be",
	"twitter.com",
	"x.com",
	"instagram.com",
	"facebook.com",
	"tiktok.com",
	"vimeo.com",
	"reddit.com",
	"twitch.tv",
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid youtube URL",
			input:   "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			want:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "valid youtu.be short URL",
			input:   "https://youtu.be/dQw4w9WgXcQ",
			want:    "https://youtu.be/dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "valid twitter URL",
			input:   "https://twitter.com/user/status/123456789",
			want:    "https://twitter.com/user/status/123456789",
			wantErr: false,
		},
		{
			name:    "valid x.com URL",
			input:   "https://x.com/user/status/123456789",
			want:    "https://x.com/user/status/123456789",
			wantErr: false,
		},
		{
			name:    "valid instagram URL",
			input:   "https://www.instagram.com/p/ABC123/",
			want:    "https://www.instagram.com/p/ABC123/",
			wantErr: false,
		},
		{
			name:    "valid tiktok URL",
			input:   "https://www.tiktok.com/@user/video/123456789",
			want:    "https://www.tiktok.com/@user/video/123456789",
			wantErr: false,
		},
		{
			name:    "valid vimeo URL",
			input:   "https://vimeo.com/123456789",
			want:    "https://vimeo.com/123456789",
			wantErr: false,
		},
		{
			name:    "valid reddit URL",
			input:   "https://www.reddit.com/r/videos/comments/abc123/test/",
			want:    "https://www.reddit.com/r/videos/comments/abc123/test/",
			wantErr: false,
		},
		{
			name:    "valid twitch URL",
			input:   "https://www.twitch.tv/videos/123456789",
			want:    "https://www.twitch.tv/videos/123456789",
			wantErr: false,
		},
		{
			name:    "valid facebook URL",
			input:   "https://www.facebook.com/watch?v=123456789",
			want:    "https://www.facebook.com/watch?v=123456789",
			wantErr: false,
		},
		{
			name:    "URL without www prefix",
			input:   "https://youtube.com/watch?v=test",
			want:    "https://youtube.com/watch?v=test",
			wantErr: false,
		},
		{
			name:    "subdomain URL",
			input:   "https://m.youtube.com/watch?v=test",
			want:    "https://m.youtube.com/watch?v=test",
			wantErr: false,
		},
		{
			name:    "domain not allowed",
			input:   "https://malicious-site.com/video",
			wantErr: true,
			errMsg:  "domain not allowed",
		},
		{
			name:    "invalid protocol - ftp",
			input:   "ftp://youtube.com/video",
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name:    "invalid protocol - javascript",
			input:   "javascript:alert(1)",
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name:    "invalid protocol - file",
			input:   "file:///etc/passwd",
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name:    "URL too long",
			input:   "https://youtube.com/" + string(make([]byte, 2050)),
			wantErr: true,
			errMsg:  "URL too long",
		},
		{
			name:    "empty URL",
			input:   "",
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name:    "malformed URL",
			input:   "not-a-url",
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name:    "similar domain attack",
			input:   "https://youtube.com.malicious.com/watch",
			wantErr: true,
			errMsg:  "domain not allowed",
		},
		{
			name:    "subdomain spoof attempt",
			input:   "https://fakeyoutube.com/watch",
			wantErr: true,
			errMsg:  "domain not allowed",
		},
		{
			name:    "http protocol allowed",
			input:   "http://youtube.com/watch?v=test",
			want:    "http://youtube.com/watch?v=test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeURL(tt.input, testAllowedDomains)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SanitizeURL() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("SanitizeURL() error = %q, want %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("SanitizeURL() unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("SanitizeURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSanitizeFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid format ID",
			input: "137",
			want:  "137",
		},
		{
			name:  "alphanumeric format",
			input: "bestvideo",
			want:  "bestvideo",
		},
		{
			name:  "mixed case format",
			input: "BestVideo",
			want:  "BestVideo",
		},
		{
			name:  "empty format returns best",
			input: "",
			want:  "best",
		},
		{
			name:  "format with special chars returns best",
			input: "137+140",
			want:  "best",
		},
		{
			name:  "format with spaces returns best",
			input: "137 140",
			want:  "best",
		},
		{
			name:  "format with injection attempt returns best",
			input: "137;rm -rf /",
			want:  "best",
		},
		{
			name:  "format with shell metacharacters returns best",
			input: "$(whoami)",
			want:  "best",
		},
		{
			name:  "format with backticks returns best",
			input: "`id`",
			want:  "best",
		},
		{
			name:  "format too long returns best",
			input: "abcdefghijklmnopqrstuvwxyz",
			want:  "best",
		},
		{
			name:  "format exactly 20 chars allowed",
			input: "abcdefghijklmnopqrst",
			want:  "abcdefghijklmnopqrst",
		},
		{
			name:  "format 21 chars returns best",
			input: "abcdefghijklmnopqrstu",
			want:  "best",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFormat(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
