# werf
a playground with containers

## Werf: Simple Logistics Simulator

Ships arrive at a quay in the wharf. The cargo is represented by a manifest. Each item has a price, weight, size category, and expiry date. 

A warehouse stores the cargo. Items above a certain price, or expiring within 7 days are considered hot, i.e. should be moved out of the wharf asap. Such hot items are kept separately for easy access, and send notifications to the harbormaster on a daily basis. 

Items that have expired are removed from the warehouse every month. That day no ships are allowed to reach the quays, and the werf shuts down.

The wharf contains several quays, and a sentry outpost routes the ships to the next available quay. 

The name "werf" is a predecessor of the English word "Wharf".

### Technical Overview
Werf uses a container having Golang sqs worker that connects to Localstack, both as docker containers. DynamoDB is used for the warehouse.

TBA
