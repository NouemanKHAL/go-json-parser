# Coding Challenge #2: Build Your Own JSON Parser

This is my solution to the second problem in the [John Crickett's Coding Challenges](https://codingchallenges.fyi/challenges/challenge-json-parser/).


## Setup

1. Clone the repo
1. Run the tool using of the following approaches:

    ```shell
    # Run the tool automatically using the go command
    $ go run . [file] [path]

    # Build a binary and run it manually
    $ go build -o gojson
    $ ./gojson [file] [path]

    # Install the binary in your environment, and run it:
    $ go install
    $ gojson [file] [path]
    ```
1. Done!


## Usage
```shell
Usage:
	gojson [FILE] [PATH]

Example:
    gojson tests/step1/valid.json .key
    cat file.json | gojson .foo.bar[0]
```


## Examples

```shell
$ gojson tests/step4/valid2.json .key-o
{"inner key": "inner value"}

$ gojson tests/step4/valid.json .key
"value"

$ gojson tests/step3/invalid.json
cannot parse value, got token 'False' at line 3:9
exit status 1
```