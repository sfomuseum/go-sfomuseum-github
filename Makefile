GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/monitor-ratelimit cmd/monitor-ratelimit/main.go

lambda-monitor:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f monitor-ratelimit.zip; then rm -f monitor-ratelimit.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/monitor-ratelimit/main.go
	zip monitor-ratelimit.zip bootstrap
	rm -f bootstrap
