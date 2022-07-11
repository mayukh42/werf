FROM golang:1.18

# set user and group
# USER rufus:rufus

WORKDIR /opt/docker/werf

# create log dir and set permissions
RUN mkdir -p /var/log/dev

# pre-copy/cache go.mod for pre-downloading dependencies 
# and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./
RUN env GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/werf cmd/werf.go

CMD ["werf", "-config=./resources/dev/config.yml"]
