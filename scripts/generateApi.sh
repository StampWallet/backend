#!/bin/sh

# https://openapi-generator.tech/docs/installation
# download from https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/6.3.0/openapi-generator-cli-6.3.0.jar

SCRIPTPATH=`realpath "$0"`
SCRIPTDIR=`dirname $SCRIPTPATH`

ls $SCRIPTDIR/../internal/api/models/*.go
rm -I $SCRIPTDIR/../internal/api/models/*.go
java -jar openapi-generator-cli-6.3.0.jar generate -g go-gin-server -i $SCRIPTDIR/../api/schema.yaml -c $SCRIPTDIR/../api/config.yaml --global-property models,test -o $SCRIPTDIR/../
java -jar openapi-generator-cli-6.3.0.jar generate -g html2 -i $SCRIPTDIR/../api/schema.yaml --global-property docs -o $SCRIPTDIR/../docs/
