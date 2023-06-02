#!/bin/sh

# https://openapi-generator.tech/docs/installation
# download the jar from https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/6.3.0/openapi-generator-cli-6.3.0.jar
# script has to be run in the same directory as the jar

SCRIPTPATH=`realpath "$0"`
SCRIPTDIR=`dirname $SCRIPTPATH`

ls $SCRIPTDIR/../internal/api/models/*.go
rm -I $SCRIPTDIR/../internal/api/models/*.go
ls $SCRIPTDIR/../client/*
rm -I $SCRIPTDIR/../client/*
java -jar openapi-generator-cli-6.3.0.jar generate -g go-gin-server -i $SCRIPTDIR/../api/schema.yaml -c $SCRIPTDIR/../api/config.yaml --global-property models,test -o $SCRIPTDIR/../
java -jar openapi-generator-cli-6.3.0.jar generate -g html -i $SCRIPTDIR/../api/schema.yaml --global-property docs -o $SCRIPTDIR/../docs/
java -jar openapi-generator-cli-6.3.0.jar generate -g go -i $SCRIPTDIR/../api/schema.yaml --global-property client -o $SCRIPTDIR/../client/
java -jar openapi-generator-cli-6.3.0.jar generate -g bash -i $SCRIPTDIR/../api/schema.yaml --global-property client -o $SCRIPTDIR/../client/
