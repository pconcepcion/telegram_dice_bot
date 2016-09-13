build:
	go build -ldflags "-X github.com/pconcepcion/telegram_dice_bot/cmd.BuildDate=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X github.com/pconcepcion/telegram_dice_bot/cmd.CommitHash=`git rev-parse HEAD`" -o telegram_dice_bot

install:
	go install

test:
	go test -v -race ./...

lint:
	go vet ./...
	golint ./...
	gocyclo -over 10 .
	errcheck ./...

clean:
	go clean

deps: dev-deps

dev-deps:
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck
	go get github.com/x-cray/logrus-prefixed-formatter
	go get github.com/pconcepcion/dice
	go get gopkg.in/telegram-bot-api.v4
	go get github.com/x-cray/logrus-prefixed-formatter
	go get github.com/pkg/errors
	go get github.com/Sirupsen/logrus

