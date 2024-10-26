package client

import (
	"context"
	"os"
	"testing"

	"github.com/openfms/jimi-iothub/utils"
	"gotest.tools/v3/assert"
)

func TestIotHubClient_HistoryVideoPlaybackRequest(t *testing.T) {
	env, err := ReadIotHubEnvironments()
	assert.NilError(t, err)
	deviceImei := os.Getenv("IOTHUB_DEVICE_IMEI")
	ctx := context.Background()
	iothubcli, err := NewIotHubClient(env)
	assert.NilError(t, err)
	req, err := iothubcli.HistoryVideoPlaybackRequest(ctx, deviceImei, DeviceModelJC450, &PlaybackCmdContent{
		InstructionID: utils.GenerateUniqueInstructionID(),
		TCPPort:       "10003",
		UDPPort:       0,
		Channel:       1,
		ResourceType:  PlaybackResourceAudioAndVideo,
		CodeType:      PlaybackAllStream,
		StorageType:   PlaybackStorageAll,
		ForwardRewind: ForwardRewindInvalid,
		PlayMethod:    PlayNormal,
		BeginTime:     "230826113555",
		EndTime:       "230826113854",
		ServerAddress: "192.168.1.1",
	})
	assert.NilError(t, err)
	resp, err := iothubcli.SendDeviceInstruction(ctx, req)
	assert.NilError(t, err)
	t.Log(resp)
	assert.Assert(t, resp.Code == 0)
	assert.Assert(t, resp.Data.Code == ResponseCodeSuccess)
}
