#!/bin/bash

if [[ $1 = "--watch" ]]
then
  # https://github.com/golang/go/issues/23449
  # Install looper if not working, you might also need readline. Check the docs of looper
  # go get -u github.com/nathany/looper
  looper
elif [[ $1 = "--integration" ]]
then
  go test ./... -tags=integration
elif [[ $1 = "help" ]]
then
  echo "Usage: Runs unit test, integration tests and unit tests in watch mode";
  echo "";
  echo "Unit Test watch mode:   ./run-tests.sh --watch";
  echo "Integration Test:       ./run-tests.sh --integration";
  echo "Unit Test:              ./run-tests.sh";
else
  go test ./...
fi