#!/bin/bash


set -e 

protoc  --go_out=paths=source_relative:. any.proto

protoc  --go_out=paths=source_relative:. ./any_test/any_test.proto
