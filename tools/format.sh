bazel run @go_sdk//:bin/gofmt -- -s -w -l .
tools/format_jsonnet.sh
yapf --in-place --recursive --parallel .
bazel run //:gazelle
bazel run //tools/ctflark
bazel run //:buildifier
