package config

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App         *App         `mapstructure:"app"`
	MySQL       *MySQL       `mapstructure:"mysql"`
	Redis       *Redis       `mapstructure:"redis"`
	AccessToken *AccessToken `mapstructure:"access_token"`
	Course      *Course      `mapstructure:"course"`
	AliOSS      *AliOSS      `mapstructure:"ali_oss"`
	AliTextMsg  *AliTextMsg  `mapstructure:"ali_textmsg"`
	Wechat      *Wechat      `mapstructure:"wechat"`
	Crons       *Crons       `mapstructure:"cron"`
}

type App struct {
	Port     string `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`
	LogLevel string `mapstructure:"log_level"`
}

type MySQL struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifetime  int    `mapstructure:"max_life_time"`
}

type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"`
}

type AccessToken struct {
	TokenExpire int64 `mapstructure:"token_expire"`
}

type Course struct {
	CancelInterval uint    `mapstructure:"cancel_interval"`
	CancelRefund   float64 `mapstructure:"cancel_refund"`
	CancelRule     string  `mapstructure:"cancel_rule"`
	CourseUnit     float64 `mapstructure:"course_unit"`
}

type Wechat struct {
	AppID     string `mapstructure:"app_id"`
	APPSecret string `mapstructure:"app_secret"`
}

type Crons struct {
	CronStartCourse      string `mapstructure:"cron_start_course"`
	CronScanCourse       string `mapstructure:"cron_scan_course"`
	CronScanCourseBefore int64  `mapstructure:"cron_start_course_before"`
}

type AliOSS struct {
	Enabled           bool   `mapstructure:"enabled"`
	Endpoint          string `mapstructure:"endpoint"`
	AccessKeyId       string `mapstructure:"accesskeyid"`
	AccessKeySecret   string `mapstructure:"accesskeysecret"`
	BucketName        string `mapstructure:"bucketname"`
	BucketUrl         string `mapstructure:"bucketurl"`
	ExpiredDay        int64  `mapstructure:"expiredday"`         // Image expiration time
	ExportFileExpired int64  `mapstructure:"export_file_expire"` // Exported file download link validity period
}
type AliTextMsg struct {
	Enable               bool   `mapstructure:"enable"`
	AccessKeyId          string `mapstructure:"access_key_id"`
	AccessKeySecret      string `mapstructure:"access_key_secret"`
	Endpoint             string `mapstructure:"endpoint"`
	IdentityTemplateCode string `mapstructure:"identity_template_code"` // Verification code template
}

var (
	configFile string
	AppConfig  Config
)

func NewConfig() *Config {
	var config = &Config{}
	return config
}

// LoadConfig loads configuration from the TOML file
func LoadConfig() error {
	configPath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	log.Printf("Loading config from: %s", configPath)
	if _, err := toml.DecodeFile(configPath, &AppConfig); err != nil {
		return err
	}

	return nil
}
