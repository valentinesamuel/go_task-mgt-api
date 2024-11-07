#!/bin/bash

rm -rf docs
swag init --dir ./,./internal/task --parseDependency --parseInternal
