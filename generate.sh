#!/bin/bash


set -e 

protoc  --go_out=paths=source_relative:. any.proto