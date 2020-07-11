set -e

if [[ $(bazel run @go_sdk//:bin/gofmt --  -s -d .) ]]; then
  echo "ERROR: go code not formatted, please run \`bazel run @go_sdk//:bin/gofmt -- -s -w -l .\`"
  exit 1
fi
tools/format_jsonnet_check.sh ||  (echo "ERROR: jsonnet files not formatted, please run \`tools/format_jsonnet.sh\`" >&2; exit 1)
yapf --recursive --parallel --quiet . || (echo "ERROR: py files not formatted, please run \`yapf --in-place --recursive --parallel .\`" >&2; exit 1)
bazel run //:gazelle -- --mode=diff || (echo "ERROR: Bazel files out-of-date, please run \`bazel run //:gazelle\`" >&2; exit 1)
bazel run //tools/ctflark -- -mode=check || (echo "ERROR: Bazel files out-of-date, please run \`bazel run //tools/ctflark\`" >&2; exit 1)
bazel run //:buildifier_check || (echo "ERROR: Bazel files not formatted, please run \`bazel run //:buildifier\`" >&2; exit 1)
bazel build //infra/jsonnet:route53 || (echo "ERROR: Bazel build route53 failed" >&2; exit 1)
bazel build //infra/jsonnet:all-docker-compose || (echo "ERROR: Bazel build docker-compose_all failed" >&2; exit 1)
bazel build //infra/jsonnet:cluster-master-docker-compose || (echo "ERROR: Bazel build docker-compose_master failed" >&2; exit 1)
bazel build //infra/jsonnet:cluster-team-docker-compose || (echo "ERROR: Bazel build docker-compose_team failed" >&2; exit 1)
bazel build //infra/jsonnet:k8s || (echo "ERROR: Bazel build k8s yaml failed" >&2; exit 1)
bazel build //:all_containers || (echo "ERROR: Bazel build all_containers failed" >&2; exit 1)
