package client

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
)

type IotHubClient struct {
	client      *resty.Client
	config      *IotHubConfig
	redis       *redis.Client
	endPointURL string
	wg          *sync.WaitGroup
}

//go:generate mockgen -source=$GOFILE -destination=../mock/iothub_cleint.go -package=$GOPACKAG
type JimiIotHub interface {
	Stop()
	EndpointURL() string
	GetEndpointHost() string
	Client() *resty.Client
	Config(canModify bool) *IotHubConfig

	SendDeviceInstruction(ctx context.Context, request *InstructRequest) (*Response, error)

	DeviceInstructionRequest(ctx context.Context, imei string, command string) (*InstructRequest, error)
	RealTimeAVRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *RealTimeCmdContent) (*InstructRequest, error)
	RealTimeAVControlRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *RealTimeControlCmdContent) (*InstructRequest, error)
	ListAVResourcesRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *AVResourceListCmdContent) (*InstructRequest, error)
	HistoryVideoPlaybackRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *PlaybackCmdContent) (*InstructRequest, error)
	HistoryPlaybackControlRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *PlaybackControlCmdContent) (*InstructRequest, error)

	GenerateDeviceConfigLinks(rtmpPrefix string) *DeviceConfigLinks
	GenerateRtmpLiveLink(secure bool, prefix string, channel int, imei string) (string, error)
	GenerateFlvLiveLink(secure bool, prefix string, channel int, imei string) (string, error)
	GenerateFLVReplayLink(secure bool, prefix string, imei string) (string, error)
	GenerateFLVHistoryLink(secure bool, channel int, imei string) (string, error)
	GenerateVideoLinks(secure bool, prefix string, channel int, imei string) (*VideoLinks, error)
}

var (
	_            JimiIotHub = &IotHubClient{}
	ErrNilConfig            = errors.New("config should not be nil")
)

func NewIotHubClient(config *IotHubConfig) (*IotHubClient, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	client := resty.New().
		SetBaseURL(config.EndPoint).
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		})
	if len(config.Proxy) > 0 {
		client.SetProxy(config.Proxy)
	}
	iotHub := &IotHubClient{
		client:      client,
		wg:          &sync.WaitGroup{},
		endPointURL: config.EndPoint,
		config:      config,
	}
	if len(config.RedisAddress) > 0 {
		iotHubRedis := redis.NewClient(&redis.Options{
			Addr:     config.RedisAddress,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
		})
		if _, e := iotHubRedis.Ping(context.Background()).Result(); e != nil {
			return nil, e
		}
		iotHub.redis = iotHubRedis
	}
	return iotHub, nil
}

func (cli *IotHubClient) Stop() {
	cli.wg.Wait()
}

// EndpointURL returns the URL of the endpoint.
func (cli *IotHubClient) EndpointURL() string {
	return cli.endPointURL
}

func (cli *IotHubClient) GetEndpointHost() string {
	hostParts := strings.Split(cli.EndpointURL(), ":")
	if len(hostParts) > 0 {
		return hostParts[0]
	}
	return ""
}

// Client returns the client.
func (cli *IotHubClient) Client() *resty.Client {
	return cli.client
}

// Config returns the client config.
func (cli *IotHubClient) Config(canModify bool) *IotHubConfig {
	if canModify {
		return cli.config
	}
	config := *cli.config
	return &config
}
