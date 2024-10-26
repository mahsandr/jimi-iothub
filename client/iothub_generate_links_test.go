package client

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateVideoLiveLink(t *testing.T) {
	tests := map[string]struct {
		endpoint string
		prefix   string
		channel  int
		secure   bool
		imei     string
		expected string
	}{
		"with prefix": {
			endpoint: "example.com:8881",
			prefix:   "live",
			channel:  0,
			secure:   false,
			imei:     "86271111111111",
			expected: "http://example.com:8881/live/0/86271111111111.flv",
		},
		"without prefix": {
			endpoint: "example.com:8881",
			prefix:   "",
			secure:   false,
			channel:  0,
			imei:     "86271111111111",
			expected: "http://example.com:8881/0/86271111111111.flv",
		},
		"with prefix secure": {
			endpoint: "example.com:8881",
			prefix:   "live",
			channel:  0,
			secure:   true,
			imei:     "86271111111111",
			expected: "https://example.com:8881/live/0/86271111111111.flv",
		},
		"without prefix secure": {
			endpoint: "example.com:8881",
			prefix:   "",
			secure:   true,
			channel:  0,
			imei:     "86271111111111",
			expected: "https://example.com:8881/0/86271111111111.flv",
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GenerateFLVLiveLink(test.secure, test.endpoint, test.prefix, test.channel, test.imei)
			assert.NilError(t, err)
			assert.Equal(t, test.expected, result, "Unexpected result")
		})
	}
}

func TestGenerateRtmpLink(t *testing.T) {
	tests := map[string]struct {
		endpoint string
		prefix   string
		channel  int
		secure   bool
		imei     string
		expected string
	}{
		"with prefix": {
			endpoint: "example.com:8881",
			prefix:   "live",
			channel:  0,
			secure:   false,
			imei:     "862711111111111",
			expected: "rtmp://example.com:8881/live/0/862711111111111",
		},
		"without prefix": {
			endpoint: "example.com:8881",
			prefix:   "",
			secure:   false,
			channel:  0,
			imei:     "862711111111111",
			expected: "rtmp://example.com:8881/0/862711111111111",
		},
		"with prefix secure": {
			endpoint: "example.com:8881",
			prefix:   "live",
			channel:  0,
			secure:   true,
			imei:     "862711111111111",
			expected: "rtmps://example.com:8881/live/0/862711111111111",
		},
		"without prefix secure": {
			endpoint: "example.com:8881",
			prefix:   "",
			secure:   true,
			channel:  0,
			imei:     "862711111111111",
			expected: "rtmps://example.com:8881/0/862711111111111",
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GenerateRtmpLiveLink(test.secure, test.endpoint, test.prefix, test.channel, test.imei)
			assert.NilError(t, err)
			assert.Equal(t, test.expected, result, "Unexpected result")
		})
	}
}

func TestGenerateHttpFLVHistoryLink(t *testing.T) {
	tests := map[string]struct {
		endpoint string
		channel  int
		secure   bool
		imei     string
		expected string
	}{
		"not secure": {
			endpoint: "example.com:8881",
			secure:   false,
			channel:  2,
			imei:     "86271111111111",
			expected: "http://example.com:8881/2/86271111111111.history.flv",
		},
		"secure": {
			endpoint: "example.com:8881",
			secure:   true,
			channel:  1,
			imei:     "86271111111111",
			expected: "https://example.com:8881/1/86271111111111.history.flv",
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GenerateFLVHistoryLink(test.secure, test.endpoint, "", test.channel, test.imei)
			assert.NilError(t, err)
			assert.Equal(t, test.expected, result, "Unexpected result")
		})
	}
}

func TestGenerateHttpFLVReplayLink(t *testing.T) {
	tests := map[string]struct {
		endpoint string
		prefix   string
		secure   bool
		imei     string
		expected string
	}{
		"with prefix": {
			endpoint: "example.com:8881",
			prefix:   "live",
			secure:   false,
			imei:     "86271111111111",
			expected: "http://example.com:8881/live/86271111111111.flv",
		},
		"without prefix": {
			endpoint: "example.com:8881",
			prefix:   "",
			secure:   false,
			imei:     "86271111111111",
			expected: "http://example.com:8881/86271111111111.flv",
		},
		"with prefix secure": {
			endpoint: "example.com:8881",
			prefix:   "live",
			secure:   true,
			imei:     "86271111111111",
			expected: "https://example.com:8881/live/86271111111111.flv",
		},
		"without prefix secure": {
			endpoint: "example.com:8881",
			prefix:   "",
			secure:   true,
			imei:     "86271111111111",
			expected: "https://example.com:8881/86271111111111.flv",
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GenerateFLVReplayLink(test.secure, test.endpoint, test.prefix, test.imei)
			assert.NilError(t, err)
			assert.Equal(t, test.expected, result, "Unexpected result")
		})
	}
}
