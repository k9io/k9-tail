#!/bin/bash

env CGO_ENABLED=0 go build -ldflags "-s"

