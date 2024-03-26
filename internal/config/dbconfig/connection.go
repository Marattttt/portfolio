package dbconfig

import (
	"regexp"
	"strings"
	"time"
)

type ConnParams struct {
	MaxConns          int
	MaxConnLifeTime   time.Duration
	HealthCheckPeriod time.Duration
}

// Safe for printing/logging
type dbConnStr string

func (c dbConnStr) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

func (c dbConnStr) String() string {
	if isPGConnURL(string(c)) {
		return sanitizeConnURL(string(c))
	}

	return sanitizeDSN(string(c))
}

func sanitizeDSN(dsn string) string {
	r := regexp.MustCompile(`[ ;]?[pP]assword=[^ ;]+`)

	return r.ReplaceAllString(dsn, "<password redacted>")
}

func isPGConnURL(connstr string) bool {
	r := regexp.MustCompile(`^[a-zA-Z]+://`)
	location := r.FindStringIndex(connstr)

	// Match is found and start is at the beginning of the string
	return location != nil && location[0] == 0
}

func sanitizeConnURL(url string) string {
	// From postgres://... match postgres://
	cleanrgx := regexp.MustCompile(`^.*\:\/\/`)

	// Trim the postgres:// part
	clean := cleanrgx.ReplaceAllString(url, "")

	// Isn't a connection string
	if clean == url {
		return url
	}

	sensitiveEnd := strings.IndexRune(clean, '@')
	passStart := strings.IndexRune(clean, ':')
	parametersStart := strings.IndexRune(clean, '/')

	if parametersStart == -1 {
		parametersStart = len(url) - 1
	}

	noSensitiveChars := sensitiveEnd == -1 || passStart == -1
	noSensitiveData := sensitiveEnd > parametersStart || passStart > parametersStart || passStart > sensitiveEnd

	if noSensitiveChars || noSensitiveData {
		return url
	}

	if sensitiveEnd > parametersStart && passStart > parametersStart {
		return url
	}

	userPass := clean[:sensitiveEnd]
	userNoPass := userPass[:strings.IndexRune(userPass, ':')]
	sanitizedUserPass := userNoPass + ":xxxxx"

	clean = strings.Replace(url, userPass, sanitizedUserPass, 1)

	return clean
}
