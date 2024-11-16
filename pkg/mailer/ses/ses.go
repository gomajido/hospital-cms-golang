package ses

import (
	"context"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
)

func InitSES(ctx context.Context, cfg *config.SESConfig) *sesv2.Client {
	defaultConfig, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(cfg.Region), awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, cfg.Session)))
	if err != nil {
		app_log.Fatalf("Error: %v", err)
	}
	client := sesv2.NewFromConfig(defaultConfig)
	return client
}
