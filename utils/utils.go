package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	validEndpointRegex = `^(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}|(?:\d{1,3}\.){3}\d{1,3})(?::\d{1,5})?$`
)

// GetEndpointURL - construct a new endpoint.
func GetEndpointURL(endpoint string) (*url.URL, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// Validate incoming endpoint URL.
	if e := isValidEndpoint(endpointURL); e != nil {
		return nil, e
	}
	return endpointURL, nil
}

func isValidEndpoint(u *url.URL) error {
	if !regexp.MustCompile(validEndpointRegex).MatchString(u.Host) || !(u.Scheme == "http" || u.Scheme == "https") || u.Path != "" {
		return fmt.Errorf("%s is an invalid endpoint", u)
	}
	return nil
}

func GenerateUniqueInstructionID() string {
	// Generate a random number between 0 and 999999999

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := rnd.Intn(1000000000)

	return fmt.Sprintf("%09d", randomNumber)
}

func FormatTime(t time.Time) string {
	return fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		t.Year()%100, t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func FormatUnixTime(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	return FormatTime(t)
}

// FormatUnixTimeIn formats a UTC Unix timestamp as YYMMDDhhmmss in the
// requested timezone. JT/T 808 / 1078 cameras interpret the BCD time
// payload as their LOCAL clock, so a UTC client request has to be
// converted into the device's local time before it goes out on the
// wire - otherwise the device looks at the wrong slice of its
// recording. Pass nil to fall back to the server's local zone, matching
// FormatUnixTime.
func FormatUnixTimeIn(unixTime int64, loc *time.Location) string {
	t := time.Unix(unixTime, 0)
	if loc != nil {
		t = t.In(loc)
	}
	return FormatTime(t)
}

// ParseDeviceTimeToUnix parses a YYMMDDhhmmss string emitted by a device
// (interpreted in the device's local timezone) and returns the
// corresponding UTC Unix timestamp. Loc nil means "treat the string as
// UTC", which is rarely what you want for a Jimi-style device.
func ParseDeviceTimeToUnix(s string, loc *time.Location) (int64, error) {
	if loc == nil {
		loc = time.UTC
	}
	t, err := time.ParseInLocation("060102150405", s, loc)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func AddOffsetToUnixTime(offset string, unixTime int64) (int64, error) {
	parts := strings.Split(offset, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid offset format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	// Calculate the duration in seconds
	offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	timestamp := time.Unix(unixTime, 0)

	return timestamp.Add(offsetDuration).Unix(), nil
}
