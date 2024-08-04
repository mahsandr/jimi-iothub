<div align="center">
    <a href="https://www.jimilab.com/">
        <img src="https://www.jimilab.com/wp-content/uploads/2022/09/logo.png" alt="jimi Logo" width="120">
    </a>
</div>

# jimi iothub client

[![Documentation](https://img.shields.io/badge/Documentation-Read%20Now-blue.svg)](https://docs.jimicloud.com/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

The IoT Hub Client Library for Go language is a powerful and user-friendly library designed to simplify the process of interacting with Jimi IoT Hub in Golang. This library provides an abstraction layer over the complex IoT Hub protocols, allowing developers to focus on building and managing their IoT solutions.

For detailed information and usage examples, please refer to our [documentation](https://docs.jimicloud.com/).

## Features

- **Device Management:** Easily manage jimi devices within the IoT Hub.
- **Cloud-to-Device Communication:** Easily send cloud-to-device messages and commands.
- **Command Generator:** Generate commands with ease for device control and management.
- **API Integration:** Call external APIs directly from your IoT application.
- **Custom Configuration:** Configure IoT Hub client behavior to match your application's requirements.
- **Redis Integration:** Integrated Redis client for tracking unique request codes.

## Installation

```bash
go get github.com/openfms/jimi-iothub
```

here is an exmaple of sending `STATUS` command to device

```go
import iothub "github.com/openfms/jimi-iothub"

func main(){
    deviceImei := os.Getenv("IOTHUB_DEVICE_IMEI")
    ctx := context.Background()
    // prepare redis to track requests
	opts, err := redis.ParseURL(env.RedisURL)
	if err!=nil{
        panic(err)
    }
	redisCli := redis.NewClient(opts)

    // prepare iothub client
	iothubCli, err := iothub.NewIotHubClient(env, redisCli)
	if err!=nil{
        panic(err)
    }
	req, err := iothubCli.DeviceInstructionRequest(ctx, deviceImei, commands.GenerateCommand(commands.STATUS))
	if err!=nil{
        panic(err)
    }
	resp, err := iothubCli.SendDeviceInstruction(ctx, req)
    if err!=nil{
        panic(err)
    }
}

```

## Interfaces

We've designed two interfaces for seamless integration: one for making API calls and another for generating commands.

### Commands

Explore the commands available in our IoT Hub Client Library. While we're actively developing and updating the library, we've started with popular commands. You can easily generate any command using our 'Generate Command' feature.

```Go
type DeviceCommands interface {
    GenerateCommand(command DeviceCommand, params ...string) string

    CoreKitSwitch(mode CoreKitMode) string
    ModifyHttpUploadServer(url string) string
    SetHttpUploadLimit(retryCount, retryInterval uint8) string
    SetRtmpServer(url string) string
    SetFileListServer(url string) string
    SetTCPServer(host HostType, serverAddr string, port uint16) string
    CapturePicture(camera CameraType) string
    CaptureVideo(camera CameraType, seconds uint8) string
    RtmpLiveStream(OnOff OnOffState, camera CameraType, pushDuration uint8) string
    UploadPlaybackVideosList() string
    ReplayVideoList(videoNames []string) string
    RtmpLogin(userName, password string) string
    UploadHistoryVideo(timeStamp string, cameraType HistoryCameraType) string
    UploadEventVideo(timeStamp string, cameraType HistoryCameraType, lengthSecond uint8) string
    UploadEventFile(fileName string) string

    RecordAudio(state EnableDisable) string
    RecordAudioSub(state EnableDisable) string
    RecordSwitch(camera RecordSwitchCamera, state EnableDisable) string
    Mirror(state OnOffState) string
    Rotation(camera CameraInOut, rotation RotationDegree) string
    SetInwardVideoQuality(quality VideoQualityInward) string
    SetOutVideoQuality(quality VideoQualityOut) string
    VideoResolutionSub(resolution VideoResolution) string
    SetSpeedUnit(speedUnit VideoResolution) string

    ChangePassword(oldPassword, newPassword string) string
    SetDeviceVolume(level VolumeLevel) string
    SetLEDOnOff(state OnOffState) string
}
```

### Client interface

```Go
type JimiIotHub interface {
    Stop()
    EndpointURL() string
    SendDeviceInstruction(ctx context.Context, request *InstructRequest) (*Response, error)

    DeviceInstructionRequest(ctx context.Context, imei string, command string) (*InstructRequest, error)
    RealTimeAVRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *RealTimeCmdContent) (*InstructRequest, error)
    RealTimeAVControlRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *RealTimeControlCmdContent) (*InstructRequest, error)
    ListAVResourcesRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *AVResourceListCmdContent) (*InstructRequest, error)
    HistoryVideoPlaybackRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *PlaybackCmdContent) (*InstructRequest, error)
    HistoryPlaybackControlRequest(ctx context.Context, imei string, deviceModel DeviceModel, cmdContent *PlaybackControlCmdContent) (*InstructRequest, error)
}

```
