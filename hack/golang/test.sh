#!/bin/bash

set -e

CGO_ENABLED=1

perform_tests() {
  echo "Performing tests"
  go test -v -covermode=atomic -coverprofile=artifacts/coverage.txt -p 4 ./...
}

perform_short_tests() {
  echo "Performing short tests"
  go test -v -short ./...
}

case $1 in
  --short)
  perform_short_tests
  ;;
  *)
  perform_tests
  ;;
esac

exit 0
