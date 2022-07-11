package werf

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mayukh42/logx/logx"
	"github.com/mayukh42/werf/awslib"
	"github.com/mayukh42/werf/config"
	"github.com/mayukh42/werf/lib"
)

type Quay struct {
	// ship worker unloads cargo and delivers to warehouse
	Id     int
	Name   string
	done   chan bool
	Log    *logx.Logger
	Queue  *awslib.QueueSvc
	Config *config.Config
}

func NewQuay(cfg *config.Config) *Quay {
	logger := lib.NewLogger(cfg.Log)
	qs := awslib.NewQueueSvc(cfg)
	return &Quay{
		Name:   "main",
		Log:    logger,
		Queue:  qs,
		Config: cfg,
	}
}

func (q *Quay) Run(name string, id int) {
	// listen forever
	q.done = lib.Terminator(q.Name)

	for {
		url, err := q.Queue.GetURL()
		// log.Printf("url error: [%v]", err)
		if err != nil {
			q.Log.Errorf("could not get queue url ", err.Error())
			q.done <- true
		}

		payload, err := q.Queue.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(url),
		})

		if err != nil {
			q.Log.Errorf("error in receiving messages from queue: ", err)

		} else {
			for _, m := range payload.Messages {
				m_ := m
				body := &ShipManifest{}
				err := json.Unmarshal([]byte(*m_.Body), body)
				if err != nil {
					// incorrect request format
					q.Log.Errorf("incorrect request format: ", err)
					continue
				}

				go func() {
					// process
					err := q.Process(body)
					if err != nil {
						// requeue w/ new visibility
						q.Queue.SQS.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
							QueueUrl:          aws.String(url),
							ReceiptHandle:     m_.ReceiptHandle,
							VisibilityTimeout: aws.Int64(q.Config.AWS.SQS.VisibilityTimeout),
						})
					} else {
						// delete from queue
						q.Queue.SQS.DeleteMessage(&sqs.DeleteMessageInput{
							QueueUrl:      aws.String(url),
							ReceiptHandle: m_.ReceiptHandle,
						})
					}
				}()
			}
		}

		q.Next()
	}
}

func (q *Quay) Next() {
	select {
	case <-q.done:
		q.Log.Infof("terminating process")
		q.Close()
		os.Exit(1)
	default:
		time.Sleep(time.Duration(q.Config.AWS.SQS.PollInterval) * time.Millisecond)
	}
}

func (q *Quay) Close() {
	if q.Log != nil {
		q.Log.Close()
	}

	if q.Queue.Log != nil {
		q.Queue.Log.Close()
	}
}
