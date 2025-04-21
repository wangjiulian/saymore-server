package token

import (
	"com.say.more.server/internal/app/constant"
	errors2 "com.say.more.server/internal/app/errors"
	"com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/repository/redis"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cristalhq/jwt/v4"

	"time"
)

type Claim struct {
	UserId     int64     `json:"user_id"`
	Phone      string    `json:"phone"`
	DeviceType string    `json:"device_type"`
	TimeStamp  time.Time `json:"time_stamp"`
}

type AuthToken struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func GenerateToken(userId int64, Phone string, deviceType string) (string, error) {
	claim := &Claim{
		UserId:     userId,
		Phone:      Phone,
		DeviceType: deviceType,
		TimeStamp:  time.Now(),
	}
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(constant.JwtTokenKey))
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(&claim)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func RefreshToken(userId int64, phone string, deviceType string) (string, error) {
	token, err := GenerateToken(userId, phone, deviceType)
	if err != nil {
		return "", err
	}
	if err := refreshRedisToken(token, userId, deviceType); err != nil {
		return "", err
	}

	return token, nil
}

func refreshRedisToken(token string, userId int64, deviceType string) error {
	tokenKey := fmt.Sprintf("%s%d_%s", constant.KeyToken, userId, deviceType)
	expireDuration := time.Duration(repository.Repos.Config.AccessToken.TokenExpire) * time.Hour
	expiredAt := time.Now().Add(expireDuration)
	authToken := AuthToken{
		Token:     token,
		ExpiredAt: expiredAt,
	}
	if err := redis.GetRedisClient().SetObject(tokenKey, authToken, expireDuration); err != nil {
		return err
	}
	return nil
}

func VerifyToken(token string, deviceType string) (int64, string, error) {
	var claim Claim
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(constant.JwtTokenKey))
	if err != nil {
		return 0, "", err
	}
	parseToken, err := jwt.Parse([]byte(token), verifier)
	if err != nil {
		return 0, "", err
	}

	if err := json.Unmarshal(parseToken.Claims(), &claim); err != nil {
		return 0, "", err
	}
	tokenKey := fmt.Sprintf("%s%d_%s", constant.KeyToken, claim.UserId, deviceType)
	var redisAuthToken AuthToken
	if err := redis.GetRedisClient().GetObject(tokenKey, &redisAuthToken); err != nil {
		return 0, "", err
	}
	if redisAuthToken.Token == "" {
		return 0, "", errors.New(errors2.ECTokenInvalid)
	}
	if redisAuthToken.Token != token {
		return 0, "", errors.New(errors2.ECTokenInvalid)
	}

	// refresh token expire time
	if err := refreshRedisToken(token, claim.UserId, deviceType); err != nil {
		return 0, "", err
	}
	return claim.UserId, claim.Phone, nil
}
