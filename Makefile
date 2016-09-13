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

dev-deps: deps
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck

deps:
#	go get -u github.com/pconcepcion/dice
	go get -u gopkg.in/telegram-bot-api.v4
	go get -u github.com/Sirupsen/logrus
	go get -u github.com/x-cray/logrus-prefixed-formatter
	go get -u github.com/pkg/errors


