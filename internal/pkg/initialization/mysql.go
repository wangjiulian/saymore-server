package initialization

import (
	"com.say.more.server/config"
	"com.say.more.server/internal/pkg/logs"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Mysql(config *config.MySQL, logger *log.Logger) (*gorm.DB, error) {
	gormLogger := logs.NewGormLogger(logger)
	//init mysql
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(config.MaxLifetime))

	db, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{Logger: gormLogger})

	return db, err

}

func DB() *gorm.DB {
	return db
}
