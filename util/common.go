package util

import (
	"net/http"
	"time"
)

const (
	timezone = "Asia/Kolkata"
)

//ProfileIDFromHTTPRequest reads the Profile-Id field set in header
func ProfileIDFromHTTPRequest(r *http.Request) string {
	return r.Header.Get("Profile-Id")
}

//EnvironmentFromHTTPRequest reads the Profile-Id field set in header
func EnvironmentFromHTTPRequest(r *http.Request) string {
	return r.Header.Get("Environment")
}

//VersionFromHTTPRequest reads the Profile-Id field set in header
func VersionFromHTTPRequest(r *http.Request) string {
	return r.Header.Get("Payabbhi-Version")
}

//FormatTime formats unix time
func FormatTime(unixTime int64, format string) string {
	localTime := time.Unix(unixTime, 0)
	ist, err := time.LoadLocation(timezone)
	if format != "" {
		if err != nil {
			return localTime.Format(format)
		}
		return localTime.In(ist).Format(format)
	}
	if err != nil {
		return localTime.Format("_2 Jan 2006")
	}
	suffix := "th"
	switch localTime.In(ist).Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return localTime.In(ist).Format("_2" + suffix + " Jan 2006")
}

// Multiply performs multiplication arithmetic operation in views
func Multiply(num1, num2 int64) int64 {
	return num1 * num2
}

//Add returns addition
func Add(num1, num2 int) int {
	return num1 + num2
}

// Divide performs division arithmetic operation in views
func Divide(num1, num2 int64) int64 {
	if num2 == 0 {
		return 0
	}
	return num1 / num2
}

var FieldSeparator = func(c rune) bool {
	return string(c) == "[" || string(c) == "]"
}
