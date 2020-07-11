load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "chromium",
    srcs = glob(["chrome-linux/**/*"]),
    mode = "0755",
    strip_prefix = ".",
    visibility = ["//visibility:public"],
)
