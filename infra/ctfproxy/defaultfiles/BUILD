load("@io_bazel_rules_go//extras:embed_data.bzl", "go_embed_data")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_embed_data(
    name = "defaultfiles_data",
    srcs = glob(
        ["*"],
        exclude = [
            "BUILD",
            "README.md",
        ],
    ),
    flatten = True,
    package = "defaultfiles",
    visibility = ["//visibility:public"],
)

# keep
go_library(
    name = "go_default_library",
    srcs = [":defaultfiles_data"],
    importpath = "github.com/adamyi/CTFProxy/infra/ctfproxy/defaultfiles",
    visibility = ["//visibility:public"],
)
