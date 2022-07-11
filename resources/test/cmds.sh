# send-message 
awsl sqs send-message \
    --queue-url http://localhost:4566/00000000000/werf-test \
    --message-body file://ship_manifest.json

# purge-queue
awsl sqs purge-queue \
    --queue-url http://localhost:4566/00000000000/werf-test
