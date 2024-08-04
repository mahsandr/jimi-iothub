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
	RtmpMediaServerPort     string `env:"RTMP_PORT" envDefault:"1936"`
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
}

func ReadIotHubEnvironments() (*IotHubConfig, error) {
	cfg := &IotHubConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
