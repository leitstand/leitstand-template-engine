#!/bin/bash

# go get -u github.com/swaggo/swag/cmd/swag
# go get github.com/rakyll/statik
# go get github.com/tdewolff/minify/cmd/minify

echo "prepare statik folder"
mkdir -p ./statik
(cd web || exit && sh buildRollup.sh)
#echo "minify"
#$GOPATH/bin/minify -a -r ./statik/
#$GOPATH/bin/minify -o ./statik/openapi/client.json ./statik/openapi/client.json
echo "generate statik"
"$GOPATH/bin/statik" -src=./statik -dest pkg
echo "remove statik folder"
rm -r ./statik/
