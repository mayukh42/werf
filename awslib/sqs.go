package awslib

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mayukh42/logx/logx"
	"github.com/mayukh42/werf/config"
	"github.com/mayukh42/werf/lib"
)

type QueueSvc struct {
	url    string
	SQS    *sqs.SQS
	Log    *logx.Logger
	Config *config.SQSCfg
}

func NewQueueSvc(cfg *config.Config) *QueueSvc {
	awsCfg := GetAWSConfig(&AWSConfigInput{
		Dev:                 cfg.AWS.Dev,
		Region:              cfg.AWS.Region,
		Host:                cfg.AWS.Host,
		Port:                cfg.AWS.Port,
		DisableRestCleaning: true,
		S3ForcePathStyle:    true,
	})

	sess := NewSession(awsCfg)
	svc := sqs.New(sess)

	qs := &QueueSvc{
		SQS:    svc,
		Config: cfg.AWS.SQS,
	}

	var sqsLog *logx.Logger
	if cfg.AWS.SQS.Log != nil {
		// configure dedicated log for sqs
		sqsLog = lib.NewLogger(cfg.AWS.SQS.Log)
	} else {
		// default log config
		sqsLog = lib.NewLogger(&config.LogCfg{
			Location: cfg.Log.Location,
			Service:  "sqs",
			Level:    cfg.Log.Level,
			File:     "sqs.log",
		})
	}
	qs.Log = sqsLog

	return qs
}

/** GetURL()
 * get aws host url for the sqs queue. if error code is 400, create new one
 */
func (q *QueueSvc) GetURL() (string, error) {
	if q.url != "" {
		return q.url, nil
	}

	url, err := q.SQS.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(q.Config.Name),
	})

	if err != nil {
		q.Log.Errorf("could not retrieve sqs queue url ", err.Error())
		if strings.Contains(err.Error(), "AWS.SimpleQueueService.NonExistentQueue") {
			return q.CreateQueue()
		}
		return "", fmt.Errorf("could not retrieve sqs queue url %v", err)
	}

	q.url = *url.QueueUrl
	q.Log.Infof("queue-url: ", q.url)

	return q.url, nil
}

func (q *QueueSvc) CreateQueue() (string, error) {
	cqo, err := q.SQS.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(q.Config.Name),
		Attributes: map[string]*string{
			"MessageRetentionPeriod": aws.String(fmt.Sprintf("%d", q.Config.Retention)),
			"VisibilityTimeout":      aws.String(fmt.Sprintf("%d", q.Config.VisibilityTimeout)),
		},
	})
	if err != nil {
		return "", err
	}

	url := *cqo.QueueUrl
	q.Log.Infof("created non-existent queue ", url)
	return url, nil
}
