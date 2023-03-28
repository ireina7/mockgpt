#!/usr/bin/env sh

# this file is used by .github/workflows/integration-test.yml

mockgpt --stub=example/simple/stub examples/hello_world/hello.proto &

# wait for generated files to be available and gripmock is up
sleep 2

go run examples/hello_world/client/*.go