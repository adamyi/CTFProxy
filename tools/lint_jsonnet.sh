for i in $(find -name \*.jsonnet -or -name \*.libsonnet); do
  bazel run --ui_event_filters=-INFO --noshow_progress @jsonnet_go//linter/jsonnet-lint $PWD/$i
done
