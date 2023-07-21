FROM ubuntu:16.04 as builer

MAINTAINER Meavy

RUN mkdir /cmd/httpServer

COPY study20210925 /cmd/httpServer

WORKDIR /cmd/httpServer

ENV VERSION 1.0.0

EXPOSE 23000

ENTRYPOINT [./study20210925]
