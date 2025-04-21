package ali_textmsg

import (
	"com.say.more.server/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type AliTextMsg struct {
	cfg    *config.AliTextMsg
	client *dysmsapi.Client
}

func NewAliTextMsg(cfg *config.AliTextMsg) (*AliTextMsg, error) {
	if !cfg.Enable {
		return &AliTextMsg{
			cfg: cfg,
		}, nil
	}
	client, err := dysmsapi.NewClientWithAccessKey("", cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &AliTextMsg{
		cfg:    cfg,
		client: client,
	}, nil
}

func (a *AliTextMsg) SendSms(phone, code, signName, templateCode string) error {
	if !a.cfg.Enable {
		return nil
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = signName
	request.TemplateCode = templateCode
	request.TemplateParam = "{\"code\":\"" + code + "\"}"
	_, err := a.client.SendSms(request)
	if err != nil {
		return err
	}

	return nil
}
