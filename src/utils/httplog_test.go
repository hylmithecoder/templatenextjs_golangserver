package utils

import (
	"strings"
	"testing"
)

func TestRedactSensitive(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "password masked",
			in:   `{"username":"admin","password":"s3cr3t"}`,
			want: `{"username":"admin","password":"***"}`,
		},
		{
			name: "token and secret masked",
			in:   `{"token":"abc","secret":"xyz","role":"staff"}`,
			want: `{"token":"***","secret":"***","role":"staff"}`,
		},
		{
			name: "spaces around colon",
			in:   `{"password" : "hunter2"}`,
			want: `{"password" : "***"}`,
		},
		{
			name: "non-sensitive untouched",
			in:   `{"name":"template","status":"ok"}`,
			want: `{"name":"template","status":"ok"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := RedactSensitive(tc.in); got != tc.want {
				t.Errorf("RedactSensitive(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSanitizeBodyEmpty(t *testing.T) {
	if got := SanitizeBody([]byte("")); got != "-" {
		t.Errorf("empty body = %q, want \"-\"", got)
	}
	if got := SanitizeBody([]byte("   \n\t ")); got != "-" {
		t.Errorf("whitespace body = %q, want \"-\"", got)
	}
}

func TestSanitizeBodyRedacts(t *testing.T) {
	got := SanitizeBody([]byte(`{"password":"leakme"}`))
	if strings.Contains(got, "leakme") {
		t.Errorf("SanitizeBody leaked a password: %q", got)
	}
	if !strings.Contains(got, `"***"`) {
		t.Errorf("SanitizeBody did not mask password: %q", got)
	}
}

func TestSanitizeBodyTruncates(t *testing.T) {
	long := strings.Repeat("a", maxBodyLogLength+500)
	got := SanitizeBody([]byte(long))
	if len(got) >= len(long) {
		t.Errorf("expected truncation, got length %d (input %d)", len(got), len(long))
	}
	if !strings.Contains(got, "truncated") {
		t.Errorf("expected truncated marker, got %q", got[len(got)-40:])
	}
}

func TestStatusColor(t *testing.T) {
	cases := map[int]string{
		200: "\x1b[32m",       // green
		201: "\x1b[32m",       // green
		301: "\x1b[36m",       // cyan
		404: "\x1b[38;5;214m", // orange
		500: "\x1b[31m",       // red
	}
	for status, want := range cases {
		if got := statusColor(status); got != want {
			t.Errorf("statusColor(%d) = %q, want %q", status, got, want)
		}
	}
}
