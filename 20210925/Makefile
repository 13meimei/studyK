# Basic Makefile for Golang project
BUILD_VERSION="1.0.0"
BUILD_DATE=`date +"%Y/%m/%d-%H:%M:%S"`

build:
	echo "build 20210925 server"
	go build -ldflags "-X main.BuildVersion=${BUILD_VERSION} -X main.BuildDate=${BUILD_DATE}" -o study20210925

clean:
	rm -rf release
	rm -rf study20210925

release: build
	rm -rf release
	mkdir release
	cp study20210925 release
