package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"services.payment/internal/appservice/port"
	"services.payment/internal/common/config"
)

func Connect(cfg config.PostgreSQLConfig) (port.Repo, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DB,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return &repoImpl{db}, nil
}

type repoImpl struct {
	db *gorm.DB
}

func (r *repoImpl) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}
