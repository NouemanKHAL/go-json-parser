#!/usr/bin/env bash

function assert() {
    local have="$1", want="$2", test_name="$3"
    local verdict="SKIPPED"
    if [ "$have" = "$want" ]; then
        verdict="PASSED"
    else
        verdict="FAILED"
    fi
    printf "\t%-30s: %10s\n" $test_name $verdict
}

go build -o gojson .

echo -e "Testing tests/step1"

./gojson tests/step1/invalid.json > /dev/null 2>&1
assert $? 1 "StepOneInvalid"

./gojson tests/step1/valid.json >/dev/null 2>&1
assert $? 0 "StepOneValid"

echo "Testing tests/step2"

./gojson tests/step2/invalid.json > /dev/null 2>&1
assert $? 1 "StepTwoInvalidOne"

./gojson tests/step2/invalid2.json >/dev/null 2>&1
assert $? 1 "StepTwoInvalidtwo"

./gojson tests/step2/valid.json > /dev/null 2>&1
assert $? 0 "StepTwoValidOne"

./gojson tests/step2/valid2.json >/dev/null 2>&1
assert $? 0 "StepTwoValidtwo"


echo "Testing tests/step3"


./gojson tests/step3/invalid.json > /dev/null 2>&1
assert $? 1 "StepThreeInvalid"

./gojson tests/step3/valid.json >/dev/null 2>&1
assert $? 0 "StepThreeValid"


echo "Testing tests/step4"


./gojson tests/step4/invalid.json > /dev/null 2>&1
assert $? 1 "StepFourInvalid"

./gojson tests/step4/valid.json > /dev/null 2>&1
assert $? 0 "StepFourValidOne"

./gojson tests/step4/valid2.json >/dev/null 2>&1
assert $? 0 "StepFourValidtwo"

