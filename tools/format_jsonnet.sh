for i in $(find -name \*.jsonnet -or -name \*.libsonnet); do
  echo $i
  bazel run --ui_event_filters=-INFO --noshow_progress @jsonnet_go//cmd/jsonnetfmt -- -i $PWD/$i
done
