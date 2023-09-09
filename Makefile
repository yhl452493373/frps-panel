export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

build: frps-panel
	cp ./config/frps-panel.ini ./bin/frps-panel.ini
	cp -r ./assets/ ./bin/assets/

frps-panel:
	rm -rf ./bin
	go build -o ./bin/frps-panel ./cmd/frps-panel
