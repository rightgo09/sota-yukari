#!/usr/bin/env bash

GOOS=linux  GOARCH=386   go build -o bin/sota-yukari_i386
GOOS=darwin GOARCH=amd64 go build -o bin/sota-yukari_macosx64

chmod +x bin/sota-yukari_*
