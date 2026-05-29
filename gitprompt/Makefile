OUTDIR := ../shell

.PHONY: build build-all clean

build:
	go build -ldflags="-s -w" -o $(OUTDIR)/gitprompt .

build-all:
	GOOS=linux  GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTDIR)/gitprompt-linux-amd64 .
	GOOS=linux  GOARCH=arm64 go build -ldflags="-s -w" -o $(OUTDIR)/gitprompt-linux-arm64 .
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(OUTDIR)/gitprompt-windows-arm64.exe .

update-deps:
	go get -u ./...
	go mod tidy

clean:
	rm -f $(OUTDIR)/gitprompt \
	      $(OUTDIR)/gitprompt-linux-amd64 \
	      $(OUTDIR)/gitprompt-linux-arm64 \
	      $(OUTDIR)/gitprompt-windows-arm64.exe
