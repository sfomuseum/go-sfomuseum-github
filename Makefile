GOMOD=vendor

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/monitor-ratelimit cmd/monitor-ratelimit/main.go

lambda-monitor:
	if test -f main; then rm -f main; fi
	if test -f monitor-ratelimit.zip; then rm -f monitor-ratelimit.zip; fi
	GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -o main cmd/monitor-ratelimit/main.go
	zip monitor-ratelimit.zip main
	rm -f main
