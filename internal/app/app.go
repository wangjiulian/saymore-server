package app

import (
	"com.say.more.server/config"
	"com.say.more.server/internal/app/job"
	repository "com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/repository/ali_oss"
	"com.say.more.server/internal/app/repository/ali_textmsg"
	"com.say.more.server/internal/app/routers"
	"com.say.more.server/internal/pkg/initialization"
	"github.com/gin-gonic/gin"
	"log"
)

func Run(cfg *config.Config) {
	logger := initialization.Logger(cfg)
	logger.Info("Starting initialization...")

	db, err := initialization.Mysql(cfg.MySQL, logger)
	if err != nil {
		log.Fatalf("init db err:%s", err)
		return
	}

	redis := initialization.Redis(cfg.Redis)

	// Initialize OSS
	aliOss, err := ali_oss.NewAliOSSBucket(cfg.AliOSS)
	if err != nil {
		logger.Fatalf("init oss err:%s", err)
		return
	}

	// Initialize AliTextMsg
	aliTextMsg, err := ali_textmsg.NewAliTextMsg(cfg.AliTextMsg)
	if err != nil {
		logger.Fatalf("init aliTextMsg err:%s", err)
		return
	}
	// Initialize repositories
	repository.InitRepositories(logger, cfg, db, redis, aliOss, aliTextMsg)

	// Initialize scheduler
	job.NewExecutor(cfg.Redis.Address, cfg.Redis.Password).Start()

	// Set Gin mode
	gin.SetMode(cfg.App.Mode)
	// Set up router
	r := routers.SetupRouter()

	// Start server
	serverAddr := ":" + cfg.App.Port
	logger.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
