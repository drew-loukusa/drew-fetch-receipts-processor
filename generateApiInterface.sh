#!/bin/sh
openapi-generator-cli generate \
  -g go-server -i api.yml \
	-o server \
	--additional-properties=outputAsLibrary=true,onlyInterfaces=true,sourceFolder=openapi