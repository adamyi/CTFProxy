for i in $(find -name \*.jsonnet -or -name \*.libsonnet); do
  bazel run --experimental_ui_limit_console_output=1 @jsonnet_go//linter/jsonnet-lint $PWD/$i
done
