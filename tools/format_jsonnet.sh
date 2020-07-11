for i in $(find -name \*.jsonnet -or -name \*.libsonnet); do
  echo $i
  bazel run --experimental_ui_limit_console_output=1 @jsonnet_go//cmd/jsonnetfmt -- -i $PWD/$i
done
