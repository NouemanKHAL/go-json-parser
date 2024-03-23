#!/usr/bin/env bash

function assert() {
    local have="$1", want="$2", test_name="$3"
    local verdict="SKIPPED"
    if [ "$have" = "$want" ]; then
        verdict="PASSED"
    else
        verdict="FAILED"
    fi
    printf "%-30s: %10s\n" $test_name $verdict
}

go build -o gojson main.go

./gojson tests/step1/invalid.json > /dev/null 2>&1
assert $? 1 "TestInvalidEmptyJSON"

./gojson tests/step1/valid.json >/dev/null 2>&1
assert $? 0 "TestValidEmptyObjectJSON"
