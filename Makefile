VERSION=0.3

.PHONY: all
all: release

.PHONY: preparation
preparation:
	mkdir -p bin

bin/rostictl-${VERSION}.linux.arm: preparation
	env GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o bin/rostictl-${VERSION}.linux.arm

bin/rostictl-${VERSION}.linux.arm64: preparation
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/rostictl-${VERSION}.linux.arm64
	
bin/rostictl-${VERSION}.linux.i386: preparation
	env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o bin/rostictl-${VERSION}.linux.i386
	
bin/rostictl-${VERSION}.linux.amd64: preparation
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/rostictl-${VERSION}.linux.amd64

bin/rostictl-${VERSION}.darwin.amd64: preparation
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/rostictl-${VERSION}.darwin.amd64
	
bin/rostictl-${VERSION}.windows.i386: preparation
	env GOOS=windows GOARCH=386 go build -o bin/rostictl-${VERSION}.windows.i386
	
bin/rostictl-${VERSION}.windows.amd64: preparation
	env GOOS=windows GOARCH=amd64 go build -o bin/rostictl-${VERSION}.windows.amd64

.PHONY: release
release: bin/rostictl-${VERSION}.linux.arm bin/rostictl-${VERSION}.linux.arm64 bin/rostictl-${VERSION}.linux.i386 bin/rostictl-${VERSION}.linux.amd64 bin/rostictl-${VERSION}.darwin.amd64 bin/rostictl-${VERSION}.windows.i386 bin/rostictl-${VERSION}.windows.amd64
	
.PHONY: clean
clean: preparation
	rm bin/*
