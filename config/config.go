package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	LOCAL_ENV      = "local"
	STAGING_ENV    = "staging"
	PRODUCTION_ENV = "production"

	CONTENT_TYPE                     = "Content-Type"
	CONTENT_TYPE_JSON                = "application/json"
	ENV                              = "APEXA_ENV"
	POSTGRE_MASTER_DB_SECRET_MANAGER = "APEXA_POSTGRE_DB_MASTER"
	POSTGRE_SLAVE_DB_SECRET_MANAGER  = "APEXA_POSTGRE_DB_STAGING"
	REDIS_ADDR_SECRET_MANAGER        = "APEXA_REDIS_ADDR"
	REDIS_PASSWORD_SECRET_MANAGER    = "APEXA_REDIS_PASSWORD"
	REDIS_DB_SECRET_MANAGER          = "APEXA_REDIS_DB"
)

type Config struct {
	Http       HttpConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Fonnte     FonnteConfig
	Starsender StarsenderConfig
	Privy      PrivyConfig
	S3         S3Config
	Bsi        Bsi
	SES        SESConfig
	SMTP       SMTPConfig
	Telegram   TelegramConfig
	Discord    DiscordConfig
	Secret     SecretManagerConfig
	Gotenberg  GotenbergConfig
	Media      MediaConfig
	Token      TokenConfig
}

type TokenConfig struct {
	TokenExpiration string `json:"TokenExpiration"`
}

type HttpConfig struct {
	Address      string        `json:"HTTP_Address"`
	ReadTimeout  time.Duration `json:"HTTP_ReadTimeout"`
	WriteTimeout time.Duration `json:"HTTP_WriteTimeout"`
	ApiPrefix    string        `json:"HTTP_ApiPrefix"`
	BaseURL      string        `json:"HTTP_BaseURL"`
}

type DatabaseConfig struct {
	DriverName      string `json:"DATABASE_DriverName"`
	Master          string `json:"DATABASE_Master"`
	Slaves          string `json:"DATABASE_Slaves"`
	MaxOpenConns    int    `json:"DATABASE_MaxOpenConns"`
	MaxIdleConns    int    `json:"DATABASE_MaxIdleConns"`
	ConnMaxLifetime int    `json:"DATABASE_ConnMaxLifetime"`
	User            string `json:"DATABASE_User"`
	Password        string `json:"DATABASE_Password"`
	Network         string `json:"DATABASE_Network"`
	DBName          string `json:"DATABASE_DBName"`
	Address         string `json:"DATABASE_Address"`
}

type RedisConfig struct {
	Address  string `json:"REDIS_Address"`
	Password string `json:"REDIS_Password"`
	DB       int    `json:"REDIS_DB"`
}

type FonnteConfig struct {
	BaseURL           string `json:"FONNTE_BaseURL"`
	OverSeasAPIKey    string `json:"FONNTE_OverSeasAPIKey"`
	DomesticAPIKey    string `json:"FONNTE_DomesticAPIKey"`
	OTPOverSeasAPIKey string `json:"FONNTE_OTPOverSeasAPIKey"`
	OTPDomesticAPIKey string `json:"FONNTE_OTPDomesticAPIKey"`
	Fallback          string `json:"FONNTE_Fallback"`
}

type StarsenderConfig struct {
	BaseURL           string `json:"STARSENDER_BaseURL"`
	OverSeasAPIKey    string `json:"STARSENDER_OverSeasAPIKey"`
	DomesticAPIKey    string `json:"STARSENDER_DomesticAPIKey"`
	OTPOverSeasAPIKey string `json:"STARSENDER_OTPOverSeasAPIKey"`
	OTPDomesticAPIKey string `json:"STARSENDER_OTPDomesticAPIKey"`
	Fallback          string `json:"STARSENDER_Fallback"`
}

type PrivyConfig struct {
	BaseURL     string `json:"PRIVY_BaseURL"`
	Username    string `json:"PRIVY_Username"`
	Password    string `json:"PRIVY_Password"`
	MerchantKey string `json:"PRIVY_MerchantKey"`
}

type S3Config struct {
	AccessKeyID     string `json:"S3_AccessKeyID"`
	SecretAccessKey string `json:"S3_SecretAccessKey"`
	SessionToken    string `json:"S3_SessionToken"`
	Region          string `json:"S3_Region"`
	Bucket          string `json:"S3_Bucket"`
	EndpointUrl     string `json:"S3_EndpointUrl"`
}

type Bsi struct {
	Url   string `json:"BSI_URL"`
	Token string `json:"BSI_Token"`
}

type SESConfig struct {
	Region      string `json:"SES_Region"`
	AccessKey   string `json:"SES_AccessKey"`
	SecretKey   string `json:"SES_SecretKey"`
	Session     string `json:"SES_Session"`
	FromAddress string `json:"SES_FromAddress"`
	FromName    string `json:"SES_FromName"`
}
type SMTPConfig struct {
	SMTPServer         string `json:"SMTP_SMTPServer"`
	SMTPPort           string `json:"SMTP_SMTPPort"`
	Identity           string `json:"SMTP_Identity"`
	Username           string `json:"SMTP_Username"`
	Password           string `json:"SMTP_Password"`
	FromAddress        string `json:"SMTP_FromAddress"`
	FromName           string `json:"SMTP_FromName"`
	AuthType           string `json:"SMTP_AuthType"`
	UseTLS             bool   `json:"SMTP_UseTLS"`
	InsecureSkipVerify bool   `json:"SMTP_InsecureSkipVerify"`
}
type TelegramConfig struct {
	BaseURL       string `json:"TELEGRAM_BaseURL"`
	ApiPrefix     string `json:"TELEGRAM_ApiPrefix"`
	DefaultChatID string `json:"TELEGRAM_DefaultChatID"`
	Token         string `json:"TELEGRAM_Token"`
}
type DiscordConfig struct {
	BaseURL          string `json:"DISCORD_BaseURL"`
	ApiPrefix        string `json:"DISCORD_ApiPrefix"`
	DefaultWebhookID string `json:"DISCORD_DefaultWebhookID"`
	Token            string `json:"DISCORD_Token"`
}

type SecretManagerConfig struct {
	Provider   string `json:"SECRET_Provider"`
	Region     string `json:"SECRET_Region"`
	AccessKey  string `json:"SECRET_AccessKey"`
	SecretKey  string `json:"SECRET_SecretKey"`
	SecretName string `json:"SECRET_SecretName"`
}

type GotenbergConfig struct {
	Url string
}

type MediaConfig struct {
	RootPath string
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func GetConfig() (*Config, error) {
	var config Config
	env := goDotEnvVariable("ENVIRONMENT")
	if env == "" {
		env = LOCAL_ENV
	}
	if env == LOCAL_ENV {
		viper.SetConfigName("env-local")
		viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/appname/") // path to look for the config file in
		viper.AddConfigPath("$HOME/.appname")
		viper.AddConfigPath("/opt/")
		viper.AddConfigPath("./config/") // optionally look for config in the working directory
		err := viper.ReadInConfig()      // Find and read the config file
		if err != nil {                  // Handle errors reading the config file
			app_log.Fatalf("Error read config file: %s", err)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			app_log.Fatalf("Error unmarshall config struct: %s", err)
		}
	}
	if env != LOCAL_ENV {
		config.loadFromSecretManager()
	}
	return &config, nil
}

func (c *Config) loadFromSecretManager() {
	app_log.Info("load config from secret manager...")
	var provider, region, accessKey, secretKey, secretName string
	provider, region, accessKey, secretKey, secretName =
		goDotEnvVariable("SECRET_PROVIDER"),
		goDotEnvVariable("SECRET_REGION"),
		goDotEnvVariable("SECRET_ACCESS_KEY_ID"),
		goDotEnvVariable("SECRET_SECRET_ACCESS_KEY"),
		goDotEnvVariable("SECRET_SECRET_NAME")

	if provider == "" || region == "" || accessKey == "" || secretKey == "" || secretName == "" {
		app_log.Fatalf("Failed to load configuration: invalidsecret manager for non local env")
	}
	secret := &SecretManagerConfig{
		Provider:   provider,
		Region:     region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		SecretName: secretName,
	}
	c.getKeyFromSecretManager(secret)

}

func (c *Config) getKeyFromSecretManager(secret *SecretManagerConfig) {
	app_log.Info("get key from secret manager...")
	staticCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(secret.AccessKey, secret.SecretKey, ""))

	// Load the default AWS configuration with static credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(staticCreds),
		config.WithRegion(secret.Region), // Replace with your region
	)
	if err != nil {
		app_log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	// Replace with the ARN or name of your secret
	secretName := secret.SecretName

	// Retrieve the secret value
	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		app_log.Fatalf("Failed to retrieve secret: %v", err)
	}

	// Check the secret value
	if result.SecretString == nil {
		app_log.Fatal("Secret value is binary, which is not supported by this example.")
	}

	secretByte := []byte(*result.SecretString)
	//parsing Http config
	err = json.Unmarshal(secretByte, &c.Http)
	if err != nil {
		app_log.Fatalf("Error parsing secret Http: %v", err)
	}

	//parsing Database config
	err = json.Unmarshal(secretByte, &c.Database)
	if err != nil {
		app_log.Fatalf("Error parsing secret Database: %v", err)
	}

	//parsing Redis config
	err = json.Unmarshal(secretByte, &c.Redis)
	if err != nil {
		app_log.Fatalf("Error parsing secret Redis: %v", err)
	}

	//parsing Fonnte config
	err = json.Unmarshal(secretByte, &c.Fonnte)
	if err != nil {
		app_log.Fatalf("Error parsing secret Fonnte: %v", err)
	}

	//parsing Starsender config
	err = json.Unmarshal(secretByte, &c.Starsender)
	if err != nil {
		app_log.Fatalf("Error parsing secret Starsender: %v", err)
	}

	//parsing S3 config
	err = json.Unmarshal(secretByte, &c.S3)
	if err != nil {
		app_log.Fatalf("Error parsing secret S3: %v", err)
	}

	//parsing Bsi config
	err = json.Unmarshal(secretByte, &c.Bsi)
	if err != nil {
		app_log.Fatalf("Error parsing secret Bsi: %v", err)
	}

	//parsing SES config
	err = json.Unmarshal(secretByte, &c.SES)
	if err != nil {
		app_log.Fatalf("Error parsing secret SES: %v", err)
	}

	//parsing SMTP config
	err = json.Unmarshal(secretByte, &c.SMTP)
	if err != nil {
		app_log.Fatalf("Error parsing secret SMTP: %v", err)
	}

	//parsing Telegram config
	err = json.Unmarshal(secretByte, &c.Telegram)
	if err != nil {
		app_log.Fatalf("Error parsing secret Telegram: %v", err)
	}

	//parsing Discord config
	err = json.Unmarshal(secretByte, &c.Discord)
	if err != nil {
		app_log.Fatalf("Error parsing secret Discord: %v", err)
	}

	//parsing Secret config
	err = json.Unmarshal(secretByte, &c.Secret)
	if err != nil {
		app_log.Fatalf("Error parsing secret Secret: %v", err)
	}
}
