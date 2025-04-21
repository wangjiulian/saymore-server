package wechat

import (
	"com.say.more.server/config"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Wechat struct {
	cfg *config.Wechat
}

type WeChatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// PhoneData represents the decrypted phone number structure
type PhoneData struct {
	PhoneNumber string `json:"phoneNumber"`
	PurePhone   string `json:"purePhoneNumber"`
	CountryCode string `json:"countryCode"`
}

func NewWechat(cfg *config.Wechat) *Wechat {
	return &Wechat{
		cfg: cfg,
	}
}

func (w *Wechat) GetSessionKey(code string) (*WeChatSession, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		w.cfg.AppID, w.cfg.APPSecret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result WeChatSession
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("WeChat API returned error: %s", result.ErrMsg)
	}

	return &result, nil
}

func (w *Wechat) DecryptWeChatPhoneData(encryptedData, iv, sessionKey string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(encryptedData)
	key, _ := base64.StdEncoding.DecodeString(sessionKey)
	ivBytes, _ := base64.StdEncoding.DecodeString(iv)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(data, data)

	// Remove padding
	length := len(data)
	padding := int(data[length-1])
	data = data[:(length - padding)]

	// Parse JSON
	var phoneData PhoneData
	if err := json.Unmarshal(data, &phoneData); err != nil {
		return "", err
	}

	return phoneData.PhoneNumber, nil
}
