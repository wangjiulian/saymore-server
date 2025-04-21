package wechat

import (
	"com.say.more.server/config"
	"testing"
)

func Test(t *testing.T) {
	return
	w := NewWechat(&config.Wechat{AppID: "", APPSecret: ""})
	session, err := w.GetSessionKey("0c1OlJFa1GMFkJ0eVTGa1bJtzQ3OlJFV")
	if err != nil {
		t.Fatal(err)
	}

	encryptedData := "kq4Yiv/V0OI+pW8i1o73CbT2hkXEgXTrpq2g/V5t9z7nuG3wlynmm9rC+bRAthEbXgLXy2cqqJJm0907FGo4DsICV0KP2/wQSOJJ4VOPLPTDL59RInUWDVfE3TWXNn4ut2G193r35nXE58y9idJcpNek5/GVfiYuqgAvSNpHx55sFPAIIp0YtGgKNs98dKLoafcXeb/Mg0xcJlAd8a/fPA=="
	iv := "ludvMPlMy952iS/3bEcsoQ=="
	phone, err := w.DecryptWeChatPhoneData(encryptedData, iv, session.SessionKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(phone)

}
