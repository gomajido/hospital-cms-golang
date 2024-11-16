package mysql

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gomajido/hospital-cms-golang/config"
)

func InitDatabase(config config.DatabaseConfig) (*sql.DB, error) {
	var db *sql.DB
	cfg := mysql.Config{
		User:                 config.User,
		Passwd:               config.Password,
		Net:                  config.Network,
		Addr:                 config.Address,
		DBName:               config.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime))
	return db, nil

}
