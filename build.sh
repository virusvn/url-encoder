#!/bin/bash
env GOOS=linux GOARCH=386 go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`" -o ./build/urlencoder.linux.x386 -v ./
env GOOS=windows GOARCH=386 go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`" -o ./build/urlencoder.window.x386.exe -v ./
go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`" -o ./build/urlencoder.macos -v ./
