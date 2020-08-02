package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/odmishien/go-auth-playground/models"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Debug    bool
}

func InitDatabase(config DatabaseConfig) (db *gorm.DB, err error) {
	param := config.User + ":" + config.Password +
		"@tcp(" + config.Host + ":" + config.Port + ")/" +
		config.Database + "?parseTime=true"

	if db, err = gorm.Open("mysql", param); err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{})
	db.LogMode(config.Debug)
	return db, nil
}
