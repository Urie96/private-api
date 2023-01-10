package main

import (
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tts/v20190823"
)

type TTSParams struct {
	Text      string `binding:"required"`
	VoiceType int64
}

func TencentTTS(p TTSParams) ([]byte, error) {
	envs := GetEnvs()
	credential := common.NewCredential(envs.TencentSDK.SecretID, envs.TencentSDK.SecretKey)
	client, _ := tts.NewClient(credential, regions.Beijing, profile.NewClientProfile())

	req := tts.NewTextToVoiceRequest()
	req.VoiceType = &p.VoiceType

	req.ProjectId = P(envs.TencentSDK.TTSProjectID)
	req.Text = &p.Text
	req.SessionId = P(MD5(p.Text))
	resp, err := client.TextToVoice(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	b, err := DecodeBase64(*resp.Response.Audio)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, nil
}
