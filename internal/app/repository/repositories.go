package repository

import (
	"com.say.more.server/config"
	"com.say.more.server/internal/app/repository/ali_oss"
	"com.say.more.server/internal/app/repository/ali_textmsg"
	redis2 "com.say.more.server/internal/app/repository/redis"
	"com.say.more.server/internal/pkg/wechat"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Repos *Repositories

type Repositories struct {
	Config     *config.Config
	DB         *gorm.DB
	Logger     *log.Logger
	Redis      *redis2.RedisClient
	AliOss     *ali_oss.AliOSSBucket
	WeChat     *wechat.Wechat
	AliTextMsg *ali_textmsg.AliTextMsg
}

func InitRepositories(logger *log.Logger, config *config.Config, db *gorm.DB, rdb *redis.Client, aliOss *ali_oss.AliOSSBucket, aliTextMsg *ali_textmsg.AliTextMsg) {
	Repos = &Repositories{
		Config:     config,
		DB:         db,
		Logger:     logger,
		Redis:      redis2.NewRedisClient(rdb, config.Redis.Prefix),
		AliOss:     aliOss,
		WeChat:     wechat.NewWechat(config.Wechat),
		AliTextMsg: aliTextMsg,
	}
}
