export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

build: frps-multiuser
	cp ./config/frps-multiuser.ini ./bin/frps-multiuser.ini
	cp -r ./assets/ ./bin/assets/

frps-multiuser:
	rm -rf ./bin
	go build -o ./bin/frps-multiuser ./cmd/frps-multiuser
