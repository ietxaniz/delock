#!/bin/bash
docker build -t delock-test .
docker run --rm delock-test go run ./examples/simple
docker rmi delock-test
