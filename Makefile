build_date = -X github.com/pconcepcion/telegram_dice_bot/cmd.BuildDate=`date '+%Y-%m-%dT%H:%M:%S%:z'`
commit_hash = -X github.com/pconcepcion/telegram_dice_bot/cmd.CommitHash=`git rev-parse HEAD`
version = -X github.com/pconcepcion/telegram_dice_bot/cmd.Version=`head -1 VERSION`

build:
	go build -ldflags "$(build_date) $(commit_hash) $(version)" -tags sqlite_omit_load_extension  -o telegram_dice_bot

build-static:
	go build -ldflags "-extldflags=-static $(build_date) $(commit_hash) $(version)" -tags sqlite_omit_load_extension  -o telegram_dice_bot_static

install:
	go install

goformat:
	go fmt ./...

test:
	go test -v ./...

test-race:
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
	docker-compose build --parallel
