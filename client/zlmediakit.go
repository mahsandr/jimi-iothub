package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ZLMediaKitClient wraps the small subset of the ZLMediaKit HTTP API we
// need to publish a JT/T 1078 stream through it: openRtpServer to grab a
// per-stream listening port, closeRtpServer to release it.
//
// The flow for a 0x9101 / 0x9201 request becomes:
//  1. openRtpServer(stream_id, tcp_mode=1) -> dynamic port
//  2. send 0x9101/0x9201 with tcpPort = dynamic port
//  3. device opens a TCP connection to ZLMediaKit on that port and
//     publishes the stream under the stream_id we provided
//  4. browser plays it at <base>/<app>/<stream_id>/hls.m3u8
//
// All HTTP errors surface as Go errors; callers should fall back to the
// static-port LKM path when the client returns one so a misconfigured
// ZLM doesn't take video offline.
type ZLMediaKitClient struct {
	baseURL string
	secret  string
	app     string
	hc      *http.Client
}

// NewZLMediaKitClient builds the client from IotHubConfig fields. Returns
// nil when ZLMediaKitURL is empty - callers should treat nil as "not
// configured" and skip ZLM entirely.
func NewZLMediaKitClient(cfg *IotHubConfig) *ZLMediaKitClient {
	if cfg == nil || strings.TrimSpace(cfg.ZLMediaKitURL) == "" {
		return nil
	}
	app := cfg.ZLMediaKitApp
	if app == "" {
		app = "live"
	}
	return &ZLMediaKitClient{
		baseURL: strings.TrimRight(cfg.ZLMediaKitURL, "/"),
		secret:  cfg.ZLMediaKitSecret,
		app:     app,
		hc:      &http.Client{Timeout: 10 * time.Second},
	}
}

// App returns the application name (URL prefix) that streams are
// published under. Defaults to "live".
func (c *ZLMediaKitClient) App() string { return c.app }

// BaseURL returns the configured ZLMediaKit HTTP endpoint.
func (c *ZLMediaKitClient) BaseURL() string { return c.baseURL }

// zlmStreamID is the canonical ZLMediaKit stream id we publish a Jimi
// device's JT/T 1078 feed under. It mirrors the "<sim>/<channel>" id
// LKM uses, but flattens the slash so the result is safe inside a
// ZLMediaKit URL path (e.g. /live/<sim>_<channel>/hls.m3u8).
func zlmStreamID(imei, channel string) string {
	if channel == "" {
		channel = "0"
	}
	return fmt.Sprintf("%s_%s", imei, channel)
}

// IsZLMediaKitEnabled reports whether the underlying IotHubClient was
// built with ZLMediaKit credentials. Other code can use this to expose
// HLS / FLV URLs that point at ZLMediaKit instead of LKM.
func (cli *IotHubClient) IsZLMediaKitEnabled() bool {
	return cli.zlm != nil
}

// ZLM returns the lazily-initialised ZLMediaKit client (nil when ZLM
// is not configured).
func (cli *IotHubClient) ZLM() *ZLMediaKitClient { return cli.zlm }

// OpenRtpServer asks ZLMediaKit to start a listener for a future JT/T
// 1078 push under the given stream id and returns the allocated TCP
// port. The stream id should be unique per concurrent device session;
// the SIM number + channel pair is a good fit.
func (c *ZLMediaKitClient) OpenRtpServer(ctx context.Context, streamID string) (int, error) {
	if streamID == "" {
		return 0, errors.New("zlmediakit: empty stream id")
	}
	q := url.Values{
		"secret":    {c.secret},
		"port":      {"0"}, // 0 = auto-pick from port_range
		"stream_id": {streamID},
		"tcp_mode":  {"1"}, // TCP listening mode; JT1078 is TCP-only here
		"app":       {c.app},
		"vhost":     {"__defaultVhost__"},
	}
	type resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Port int    `json:"port"`
	}
	var out resp
	if err := c.get(ctx, "/index/api/openRtpServer", q, &out); err != nil {
		return 0, err
	}
	if out.Code != 0 {
		return 0, fmt.Errorf("zlmediakit openRtpServer code=%d msg=%s", out.Code, out.Msg)
	}
	if out.Port == 0 {
		return 0, errors.New("zlmediakit: openRtpServer returned port 0")
	}
	return out.Port, nil
}

// CloseRtpServer releases the listener for a stream id. Safe to call
// even if the listener was already closed (ZLM returns code=0 either
// way for not-found).
func (c *ZLMediaKitClient) CloseRtpServer(ctx context.Context, streamID string) error {
	if streamID == "" {
		return nil
	}
	q := url.Values{
		"secret":    {c.secret},
		"stream_id": {streamID},
	}
	type resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var out resp
	if err := c.get(ctx, "/index/api/closeRtpServer", q, &out); err != nil {
		return err
	}
	if out.Code != 0 {
		return fmt.Errorf("zlmediakit closeRtpServer code=%d msg=%s", out.Code, out.Msg)
	}
	return nil
}

func (c *ZLMediaKitClient) get(ctx context.Context, path string, q url.Values, out interface{}) error {
	u := c.baseURL + path + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("zlmediakit %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zlmediakit %s: HTTP %d: %s", path, resp.StatusCode, body)
	}
	return json.Unmarshal(body, out)
}
