# Basic Makefile for Golang project
BUILD_VERSION="1.0.0"
BUILD_DATE=`date +"%Y/%m/%d-%H:%M:%S"`

build:
	echo "build 20210925 server"
	go build -ldflags "-X main.BuildVersion=${BUILD_VERSION} -X main.BuildDate=${BUILD_DATE}" -o study20210925
	chmod +x study20210925

clean:
	echo "clean 20210925 server"
	rm -rf release
	rm -rf study20210925

release: build
	rm -rf release
	mkdir release
	echo "release 20210925 server"
	cp study20210925 release
	echo "build 20210925 server container"
	docker build -t meavy/study20210925:v1.0.0 -f Dockerfile .	

push:
	echo "push meavy/study20210925:v1.0.0"
	docker push meavy/study20210925:v1.0.0
