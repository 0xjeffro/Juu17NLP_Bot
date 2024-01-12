package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var conn *gorm.DB

func Init() error {
	DbName := os.Getenv("DB_NAME")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", DbUser, DbPassword, DbHost, DbPort, DbName)

	var err error
	conn, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       DSN + "&loc=Asia%2fShanghai", // DSN data source name
		DisableDatetimePrecision:  true,                         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                        // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		return err
	}
	sqlDB, err := conn.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetConn() *gorm.DB {
	if conn == nil {
		err := Init()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(conn)
	return conn
}
