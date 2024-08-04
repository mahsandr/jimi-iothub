package client

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/go-querystring/query"
)

var (
	serverFlagID   int64
	flagIDMutex    sync.Mutex
	requestID      int64
	requestIDMutex sync.Mutex
)

func getServerFlagID() int64 {
	flagIDMutex.Lock()
	defer flagIDMutex.Unlock()
	serverFlagID++
	return serverFlagID
}

func getRequestID() int64 {
	requestIDMutex.Lock()
	defer requestIDMutex.Unlock()
	requestID++
	return requestID
}

type InstructRequest struct {
	DeviceIMEI   string             `url:"deviceImei,required"`
	ServerFlagID int64              `url:"serverFlagId,required"`
	ProNo        ProNumber          `url:"proNo,required"`
	Platform     RequestPlatform    `url:"platform,required"`
	RequestID    int64              `url:"requestId,required"`
	CmdContent   string             `url:"cmdContent,required"`
	CmdType      RequestCommandType `url:"cmdType"`
	Language     string             `url:"language"`
	Sync         bool               `url:"sync"`
	OfflineFlag  bool               `url:"offlineFlag"`
	Timeout      int                `url:"timeOut"`
	Token        string             `url:"token"`
}

func (cli *IotHubClient) DeviceInstructionRequest(ctx context.Context, imei string, command string) (*InstructRequest, error) {
	if len(command) == 0 {
		return nil, ErrEmptyCmdContent
	}
	var (
		reqID  = getRequestID()
		flagID = getServerFlagID()
		err    error
	)
	if cli.redis != nil {
		reqID, err = cli.redis.Incr(ctx, RedisRequestIDKey).Result()
		if err != nil {
			return nil, err
		}
		flagID, err = cli.redis.Incr(ctx, RedisServerFlagIDKey).Result()
		if err != nil {
			return nil, err
		}
	}
	return &InstructRequest{
		DeviceIMEI:   imei,
		ProNo:        ProNoOnlineCommand,
		Platform:     RequestPlatformWeb,
		CmdType:      NormallnsCommandType,
		Token:        cli.config.Token,
		OfflineFlag:  cli.config.OfflineFlag,
		Timeout:      cli.config.Timeout,
		Sync:         cli.config.Sync,
		RequestID:    reqID,
		ServerFlagID: flagID,
		CmdContent:   command,
	}, nil
}

func (cli *IotHubClient) SendDeviceInstruction(ctx context.Context, request *InstructRequest) (*Response, error) {
	values, err := query.Values(request)
	if err != nil {
		return nil, err
	}
	// Send the POST request with x-www-form-urlencoded data
	resp, err := cli.client.R().
		SetContext(ctx).
		SetBody(values.Encode()).
		Post(cli.endPointURL + "/api/device/sendInstruct")

	if err != nil {
		return nil, err
	}
	apiResponse := &Response{}
	err = json.Unmarshal(resp.Body(), apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, nil
}
