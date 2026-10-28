package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// maxBodyLogLength caps how many characters of a request/response body are
// logged so large payloads don't flood the console.
const maxBodyLogLength = 2048

// sensitiveJSONKeyRe matches credential-like JSON string fields so their values
// can be masked before logging.
var sensitiveJSONKeyRe = regexp.MustCompile(`("(?:password|token|secret|pin|access_token|refresh_token)"\s*:\s*)"[^"]*"`)

// RedactSensitive masks the values of credential-like JSON keys (password,
// token, secret, pin, ...) so they never appear in verbose logs.
func RedactSensitive(s string) string {
	return sensitiveJSONKeyRe.ReplaceAllString(s, `$1"***"`)
}

// SanitizeBody prepares a request/response body for logging: it trims
// whitespace, redacts credentials, collapses empty bodies to "-", and
// truncates anything longer than maxBodyLogLength.
func SanitizeBody(body []byte) string {
	s := strings.TrimSpace(string(body))
	if s == "" {
		return "-"
	}
	s = RedactSensitive(s)
	if len(s) > maxBodyLogLength {
		return s[:maxBodyLogLength] + fmt.Sprintf("... (%d bytes total, truncated)", len(s))
	}
	return s
}

// statusColor returns an ANSI color code for an HTTP status class:
// green 2xx, cyan 3xx, orange 4xx, red 5xx.
func statusColor(status int) string {
	switch {
	case status >= 500:
		return "\x1b[31m" // red
	case status >= 400:
		return "\x1b[38;5;214m" // orange
	case status >= 300:
		return "\x1b[36m" // cyan
	default:
		return "\x1b[32m" // green
	}
}

// LogHTTP prints a verbose, colored log entry for a single API request and its
// response: method, path, status, latency, plus the (redacted, truncated)
// request and response bodies.
func LogHTTP(method, path string, status int, latency time.Duration, reqBody, respBody []byte) {
	timestamp := time.Now().Format("15:04:05")
	color := statusColor(status)
	fmt.Printf("\x1b[35m[%s] [HTTP]\x1b[0m \x1b[1m%-6s\x1b[0m %s \x1b[90m->\x1b[0m %s%d\x1b[0m \x1b[90m(%v)\x1b[0m\n",
		timestamp, method, path, color, status, latency)
	fmt.Printf("    \x1b[90m-> request :\x1b[0m %s\n", SanitizeBody(reqBody))
	fmt.Printf("    \x1b[90m<- response:\x1b[0m %s\n", SanitizeBody(respBody))
}
