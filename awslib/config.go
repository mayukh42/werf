package awslib

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWSConfigInput struct {
	Dev                 bool
	Region              string
	Host                string
	Port                int
	DisableRestCleaning bool
	S3ForcePathStyle    bool
}

func (a *AWSConfigInput) isProd() bool {
	return !a.Dev
}

func GetAWSConfig(cfg *AWSConfigInput) *aws.Config {
	awsCfg := &aws.Config{
		Region:                         aws.String(cfg.Region),
		DisableRestProtocolURICleaning: aws.Bool(cfg.DisableRestCleaning),
		S3ForcePathStyle:               aws.Bool(cfg.S3ForcePathStyle),
	}

	if !cfg.isProd() {
		endpoint := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
		log.Printf("[dev] aws endpoint: %s", endpoint)
		awsCfg.Endpoint = aws.String(endpoint)
	}

	return awsCfg
}

func NewSession(cfg *aws.Config) *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: *cfg,
	}))
}
