#!/bin/bash

rm -rf docs
swag init --dir ./cmd/api,./internal/task --parseDependency --parseInternal
