build_date = -X github.com/pconcepcion/telegram_dice_bot/cmd.BuildDate=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
commit_hash = -X github.com/pconcepcion/telegram_dice_bot/cmd.CommitHash=`git rev-parse HEAD`

build:
	go build -ldflags "$(build_date) $(commit_hash)" -tags sqlite_omit_load_extension  -o telegram_dice_bot

build-static:
	go build -ldflags "-extldflags=-static $(build_date) $(commit_hash)" -tags sqlite_omit_load_extension  -o telegram_dice_bot_static

install:
	go install

goformat:
	go fmt ./...

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
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck

docker: build-static
	docker build -t telegram_dice_bot:0.1 .
docker-compose: build-static
	docker-compose build 
