package utils

import "time"

const JavascriptISOString = "2006-01-02T15:04:05.999Z07:00"

// JavascriptISOStringToTime is a utility function which will take a string
// formatted using the Javascript `toUTCString()` function and return a Golang
// time object.
func JavascriptISOStringToTime(s string) (time.Time, error) {
	// Special Thanks: https://stackoverflow.com/a/36582947
	dt, err := time.ParseInLocation(JavascriptISOString, s, time.UTC)
	if err != nil {
		return time.Now().UTC(), err
	}
	return dt, err
}
