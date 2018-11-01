#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w"
zip -r9 aws-s3-bucket-purger-1.0.0-x86_64-apple-darwin.zip aws-s3-bucket-purger
GOOS=darwin GOARCH=386 go build -ldflags="-s -w" 
zip -r9 aws-s3-bucket-purger-1.0.0-i386-apple-darwin.zip aws-s3-bucket-purger
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" 
zip -r9 aws-s3-bucket-purger-1.0.0-x86_64-pc-windows.zip aws-s3-bucket-purger.exe
GOOS=windows GOARCH=386 go build -ldflags="-s -w" 
zip -r9 aws-s3-bucket-purger-1.0.0-i386-pc-windows.zip aws-s3-bucket-purger.exe
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" 
zip -r9 aws-s3-bucket-purger-1.0.0-x86_64-unknow-linux.zip aws-s3-bucket-purger
GOOS=linux GOARCH=386 go build -ldflags="-s -w" 
zip -r9 aws-s3-bucket-purger-1.0.0-i386-unknown-linux.zip aws-s3-bucket-purger