package client

import "github.com/caarlos0/env/v6"

type IotHubConfig struct {
	EndPoint                string `env:"IOTHUB_ENDPOINT,notEmpty"`
	Proxy                   string `env:"CLIENT_PROXY"`
	Token                   string `env:"IOTHUB_TOKEN" envDefault:"123456"`
	JimiGatewayPort         string `env:"JIMI_GATEWAY_PORT" envDefault:"21100"`
	JTGatewayPort           string `env:"JT_GATEWAY_PORT" envDefault:"21122"`
	FileStoragePort         string `env:"FILE_STORAGE_PORT" envDefault:"23010"`
	HttpFlvMediaServerPort  string `env:"FLV_HTTP_PORT" envDefault:"8881"`
	HttpsFlvMediaServerPort string `env:"FLV_HTTPS_PORT" envDefault:"8890"`
	HLSPort                 string `env:"HLS_PORT" envDefault:"8080"`
	HLSSecurePort           string `env:"HLS_HTTPS_PORT" envDefault:"8088"`
	RtmpMediaServerPort     string `env:"RTMP_PORT" envDefault:"1936"`
	// MediaServerHost overrides the IOTHUB_ENDPOINT host for HLS / FLV / RTMP
	// link generation. Set this to the SRS (or other media origin) hostname
	// when the streaming server runs separately from the Jimi command gateway.
	MediaServerHost string `env:"MEDIA_SERVER_HOST"`
	LiveVideoPort           string `env:"LIVE_VIDEO_PORT" envDefault:"10002"`
	HistoryVideoPort        string `env:"HISTORY_VIDEO_PORT" envDefault:"10003"`
	APIPort                 string `env:"API_PORT" envDefault:"9080"`
	VideoIP                 string `env:"IOTHUB_VIDEO_IP"`
	InstructionServicePort  string `env:"INSTRUCTION_SERVICE_PORT" envDefault:"10088"`
	RedisAddress            string `env:"IOTHUB_REDIS_ADDRESS"`
	RedisPassword           string `env:"IOTHUB_REDIS_PASSWORD"`
	RedisDB                 int    `env:"IOTHUB_REDIS_DB" envDefault:"0"`
	Timeout                 int    `env:"JIMI_REQUEST_TIMEOUT" envDefault:"30"`
	OfflineFlag             bool   `env:"JIMI_OFFLINE_FLAG" envDefault:"false"`
	Sync                    bool   `env:"JIMI_REQUEST_SYNC" envDefault:"true"`

	// DefaultLiveCodeStreamType picks the bitrate the device emits for
	// 0x9101 (real-time video) when the caller leaves CodeStreamType
	// unset. 0 = MainStream (high bitrate, ~1.5-3 Mbps), 1 = SubStream
	// (~0.3-0.8 Mbps). Sub is the right default for cellular-billed
	// fleets - clients can still ask for main explicitly per request.
	DefaultLiveCodeStreamType uint8 `env:"JIMI_LIVE_CODE_STREAM_TYPE" envDefault:"1"`

	// DefaultPlaybackCodeType picks the bitrate for 0x9201 (history
	// playback). 0 = AllStream, 1 = MainStream, 2 = SubStream. Defaults
	// to SubStream for the same SIM-cost reason.
	DefaultPlaybackCodeType uint8 `env:"JIMI_PLAYBACK_CODE_TYPE" envDefault:"2"`

	// ZLMediaKitURL toggles the dynamic-port ZLMediaKit publishing path
	// for 0x9101/0x9201. Set it to the ZLM HTTP API origin (e.g.
	// "http://zlmediakit:18088") to opt in; leave empty to keep the
	// static LIVE_VIDEO_PORT / HISTORY_VIDEO_PORT LKM flow.
	ZLMediaKitURL    string `env:"ZLMEDIAKIT_URL"`
	ZLMediaKitSecret string `env:"ZLMEDIAKIT_SECRET"`
	ZLMediaKitApp    string `env:"ZLMEDIAKIT_APP" envDefault:"live"`
}

func ReadIotHubEnvironments() (*IotHubConfig, error) {
	cfg := &IotHubConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
