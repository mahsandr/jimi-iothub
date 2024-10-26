package client

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

const (
	HttpFLVLiveLinkFormat    = "{{FlvAddress}}/{{Prefix}}/{{channel}}/{{IMEI}}.flv"
	HttpFLVHistoryLinkFormat = "{{FlvAddress}}/{{Prefix}}/{{channel}}/{{IMEI}}.history.flv"
	HttpFLVReplayLinkFormat  = "{{FlvAddress}}/{{Prefix}}/{{IMEI}}.flv"
	RtmpLiveLinkFormat       = "{{RtmpAddress}}/{{Prefix}}/{{channel}}/{{IMEI}}"
	validEndpointRegex       = `^(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}|(?:\d{1,3}\.){3}\d{1,3})(?::\d{1,5})?$`
)

type VideoLinks struct {
	RtmpLiveLink, FlvLiveLink     string
	FLVHistoryLink, FLVReplayLink string
}

// GenerateFLVLiveLink generates flv live video link
// Example: http://example.com:8881/live/0/86271111111111.flv
func GenerateFLVLiveLink(secure bool, endpoint, prefix string, channel int, imei string) (string, error) {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	// Construct a secured endpoint URL.
	endpointURLStr := scheme + "://" + endpoint
	endpointURL, err := url.Parse(endpointURLStr)
	if err != nil {
		return "", err
	}
	// Validate incoming endpoint URL.
	if !regexp.MustCompile(validEndpointRegex).MatchString(endpointURL.Host) || endpointURL.Path != "" {
		return "", fmt.Errorf("%s is an invalid endpoint", endpointURL.String())
	}
	replacer := strings.NewReplacer(
		"{{FlvAddress}}", endpointURL.String(),
		"{{Prefix}}", prefix,
		"{{channel}}", fmt.Sprintf("%d", channel),
		"{{IMEI}}", imei,
	)
	format := HttpFLVLiveLinkFormat
	if prefix == "" {
		format = strings.ReplaceAll(format, "{{Prefix}}/", "")
	}
	link := replacer.Replace(format)
	return link, nil
}

// GenerateFLVHistoryLink generates flv history video link
// Example: http://120.78.224.93:8881/3/868120303960873.history.flv
func GenerateFLVHistoryLink(secure bool, endpoint, prefix string, channel int, imei string) (string, error) {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	// Construct a secured endpoint URL.
	endpointURLStr := scheme + "://" + endpoint
	endpointURL, err := url.Parse(endpointURLStr)
	if err != nil {
		return "", err
	}
	// Validate incoming endpoint URL.
	if !regexp.MustCompile(validEndpointRegex).MatchString(endpointURL.Host) || endpointURL.Path != "" {
		return "", fmt.Errorf("%s is an invalid endpoint", endpointURL.String())
	}
	replacer := strings.NewReplacer(
		"{{FlvAddress}}", endpointURL.String(),
		"{{Prefix}}", prefix,
		"{{channel}}", fmt.Sprintf("%d", channel),
		"{{IMEI}}", imei,
	)
	format := HttpFLVHistoryLinkFormat
	if prefix == "" {
		format = strings.ReplaceAll(format, "{{Prefix}}/", "")
	}
	link := replacer.Replace(format)
	return link, nil
}

// GenerateFLVReplayLink generates flv history video link for replay list
// Example: http://example.com:8881/live/868120303960873.flv
func GenerateFLVReplayLink(secure bool, endpoint string, prefix string, imei string) (string, error) {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	// Construct a secured endpoint URL.
	endpointURLStr := scheme + "://" + endpoint
	endpointURL, err := url.Parse(endpointURLStr)
	if err != nil {
		return "", err
	}
	// Validate incoming endpoint URL.
	if !regexp.MustCompile(validEndpointRegex).MatchString(endpointURL.Host) || endpointURL.Path != "" {
		return "", fmt.Errorf("%s is an invalid endpoint", endpointURL.String())
	}
	replacer := strings.NewReplacer(
		"{{FlvAddress}}", endpointURL.String(),
		"{{Prefix}}", prefix,
		"{{IMEI}}", imei,
	)
	format := HttpFLVReplayLinkFormat
	if prefix == "" {
		format = strings.ReplaceAll(format, "{{Prefix}}/", "")
	}
	link := replacer.Replace(format)
	return link, nil
}

// GenerateRtmpLiveLink generates rtmp video link
// Example: rtmp://example.com:1936/live/0/86271111111111
func GenerateRtmpLiveLink(secure bool, endpoint, prefix string, channel int, imei string) (string, error) {
	scheme := "rtmp"
	if secure {
		scheme = "rtmps"
	}
	// Construct a secured endpoint URL.
	endpointURLStr := scheme + "://" + endpoint
	endpointURL, err := url.Parse(endpointURLStr)
	if err != nil {
		return "", err
	}
	// Validate incoming endpoint URL.
	if !regexp.MustCompile(validEndpointRegex).MatchString(endpointURL.Host) || endpointURL.Path != "" {
		return "", fmt.Errorf("%s is an invalid endpoint", endpointURL.String())
	}
	replacer := strings.NewReplacer(
		"{{RtmpAddress}}", endpointURL.String(),
		"{{Prefix}}", prefix,
		"{{channel}}", fmt.Sprintf("%d", channel),
		"{{IMEI}}", imei,
	)
	format := RtmpLiveLinkFormat
	if prefix == "" {
		format = strings.ReplaceAll(format, "{{Prefix}}/", "")
	}
	link := replacer.Replace(format)
	return link, nil
}

func (cli *IotHubClient) GenerateFlvLiveLink(secure bool, prefix string, channel int, imei string) (string, error) {
	port := cli.config.HttpFlvMediaServerPort
	if secure {
		port = cli.config.HttpsFlvMediaServerPort
	}
	flvEndpoint := net.JoinHostPort(cli.GetEndpointHost(), port)
	return GenerateFLVLiveLink(secure, flvEndpoint, prefix, channel, imei)
}

func (cli *IotHubClient) GenerateRtmpLiveLink(secure bool, prefix string, channel int, imei string) (string, error) {
	rtmpEndpoint := net.JoinHostPort(cli.GetEndpointHost(), cli.config.RtmpMediaServerPort)
	return GenerateRtmpLiveLink(secure, rtmpEndpoint, prefix, channel, imei)
}

func (cli *IotHubClient) GenerateFLVHistoryLink(secure bool, prefix string, channel int, imei string) (string, error) {
	port := cli.config.HttpFlvMediaServerPort
	if secure {
		port = cli.config.HttpsFlvMediaServerPort
	}
	flvEndpoint := net.JoinHostPort(cli.GetEndpointHost(), port)
	return GenerateFLVHistoryLink(secure, flvEndpoint, prefix, channel, imei)
}

func (cli *IotHubClient) GenerateFLVReplayLink(secure bool, prefix string, imei string) (string, error) {
	port := cli.config.HttpFlvMediaServerPort
	if secure {
		port = cli.config.HttpsFlvMediaServerPort
	}
	flvEndpoint := net.JoinHostPort(cli.GetEndpointHost(), port)
	return GenerateFLVReplayLink(secure, flvEndpoint, prefix, imei)
}

func (cli *IotHubClient) GenerateVideoLinks(secure bool, prefix string, channel int, imei string) (*VideoLinks, error) {
	rtmpLiveLink, err := cli.GenerateRtmpLiveLink(secure, prefix, channel, imei)
	if err != nil {
		return nil, err
	}
	flvLiveLink, err := cli.GenerateFlvLiveLink(secure, prefix, channel, imei)
	if err != nil {
		return nil, err
	}
	flvReplayLink, err := cli.GenerateFLVReplayLink(secure, prefix, imei)
	if err != nil {
		return nil, err
	}
	flvHistoryLink, err := cli.GenerateFLVHistoryLink(secure, prefix, channel, imei)
	if err != nil {
		return nil, err
	}
	return &VideoLinks{
		RtmpLiveLink:   rtmpLiveLink,
		FlvLiveLink:    flvLiveLink,
		FLVReplayLink:  flvReplayLink,
		FLVHistoryLink: flvHistoryLink,
	}, nil
}

type DeviceConfigLinks struct {
	RtmpAddress    string
	HttpUploadLink string
}

func (cli *IotHubClient) GenerateDeviceConfigLinks(rtmpPrefix string) *DeviceConfigLinks {
	rtmpUrl := net.JoinHostPort(cli.GetEndpointHost(), cli.config.RtmpMediaServerPort)
	if len(rtmpPrefix) > 0 {
		rtmpUrl = rtmpUrl + "/" + rtmpPrefix
	}
	httpUploadUrl := fmt.Sprintf("http://%s/upload", net.JoinHostPort(cli.GetEndpointHost(), cli.config.FileStoragePort))
	return &DeviceConfigLinks{
		RtmpAddress:    rtmpUrl,
		HttpUploadLink: httpUploadUrl,
	}
}
