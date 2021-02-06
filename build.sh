#!/bin/bash
export GO111MODULE=on
mkdir -p output/bin
go build -o output/bin/experiment
