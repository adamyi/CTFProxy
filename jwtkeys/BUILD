load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

exports_files(glob(["jwt*"]))

go_library(
    name = "go_default_library",
    srcs = ["generate.go"],
    importpath = "github.com/adamyi/CTFProxy/jwtkeys",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "jwtkeys",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
