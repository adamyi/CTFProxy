for i in $(find -name \*.jsonnet -or -name \*.libsonnet); do
  bazel run --ui_event_filters=-INFO --noshow_progress @jsonnet_go//cmd/jsonnetfmt -- $PWD/$i > format_check_tmp.jsonnet
  if cmp -s $i format_check_tmp.jsonnet; then
    echo OK $i
  else
    echo NOT OK $i
    rm format_check_tmp.jsonnet
    exit 1
  fi
done
rm format_check_tmp.jsonnet
