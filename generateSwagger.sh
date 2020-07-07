#!/bin/bash

# go get -u github.com/swaggo/swag/cmd/swag

echo "generate swagger"
$GOPATH/bin/swag init --parseVendor --generalInfo cmd/template-engine/template-engine.go --output ./web/src/openapi/
rm ./web/src/openapi/swagger.json
rm ./web/src/openapi/docs.go