package dependency

import (
	"context"
	"database/sql"
	"net/smtp"
	"time"

	"github.com/gomajido/hospital-cms-golang/pkg/mailer/ses"
	mailerSmtp "github.com/gomajido/hospital-cms-golang/pkg/mailer/smtp"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	mailerSesv2 "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
	db "github.com/gomajido/hospital-cms-golang/pkg/db/mysql"
	"github.com/gomajido/hospital-cms-golang/pkg/db/redis"
	http_wrapper "github.com/gomajido/hospital-cms-golang/pkg/http_client/http"
	s3Driver "github.com/gomajido/hospital-cms-golang/pkg/storage/s3"
)

type Drivers struct {
	Db       *sql.DB
	Redis    *redis.Redis
	S3       *s3.Client
	Http     http_wrapper.IHTTPClientWrapper
	SES      *mailerSesv2.Client
	SMTPAuth *smtp.Auth
}

func InitDrivers(cfg *config.Config) *Drivers {
	return &Drivers{
		Db:       initDB(&cfg.Database),
		Redis:    initRedis(&cfg.Redis),
		S3:       initS3(&cfg.S3),
		Http:     initHTTP(&cfg.Http),
		SES:      initSES(&cfg.SES),
		SMTPAuth: initSMTPAuth(&cfg.SMTP),
	}
}

func initDB(cfg *config.DatabaseConfig) *sql.DB {
	db, err := db.InitDatabase(*cfg)
	if err != nil {
		app_log.Fatalf("Fail connect to database: %s\n", err)
	}
	return db
}

func initRedis(cfg *config.RedisConfig) *redis.Redis {
	return redis.InitRedis(cfg)
}

func initS3(cfg *config.S3Config) *s3.Client {
	return s3Driver.InitS3(cfg)
}

func initHTTP(cfg *config.HttpConfig) http_wrapper.IHTTPClientWrapper {
	return http_wrapper.NewHTTPClientWrapper(time.Duration(time.Second * cfg.ReadTimeout))
}

func initSES(cfg *config.SESConfig) *mailerSesv2.Client {
	ctx := context.TODO()
	return ses.InitSES(ctx, cfg)
}

func initSMTPAuth(cfg *config.SMTPConfig) *smtp.Auth {
	ctx := context.TODO()
	return mailerSmtp.InitSMTPAuth(ctx, cfg)
}
