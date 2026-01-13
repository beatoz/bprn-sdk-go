#!/bin/bash

echo "go build"
go build ./...

echo "go vet"
go vet ./...
