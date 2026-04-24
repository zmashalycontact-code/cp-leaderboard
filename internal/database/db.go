package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct{}

func ConfigFromEnv() Config {
	return Config{}
}

func Connect(cfg Config) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	var db *gorm.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{

			Logger: logger.Default.LogMode(logger.Warn),
		})

		if err == nil {
			fmt.Printf("✅ تم الاتصال بقاعدة البيانات في المحاولة رقم %d\n", i)
			break
		}

		fmt.Printf("⚠️ محاولة اتصال فاشلة (%d/5): %v. إعادة المحاولة خلال 5 ثوانٍ...\n", i, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("❌ فشل الاتصال النهائي بقاعدة البيانات: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func AutoMigrate(db *gorm.DB, dst ...interface{}) error {
	if err := db.AutoMigrate(dst...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}
