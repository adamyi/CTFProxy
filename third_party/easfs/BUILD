load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

container_image(
    name = "image",
    base = ":easfs-image",
    tars = [
        ":templates",
        ":static",
    ],
    visibility = ["//visibility:public"],
)

go_image(
    name = "easfs-image",
    args = [
        "-site",
        "/site",
    ],
    embed = [":go_default_library"],
)

pkg_tar(
    name = "templates",
    srcs = glob(["templates/**/*"]),
    mode = "0755",
    strip_prefix = ".",
)

pkg_tar(
    name = "static",
    srcs = glob(["static/**/*"]),
    mode = "0755",
    strip_prefix = ".",
)

go_library(
    name = "go_default_library",
    srcs = [
        "BookParser.go",
        "FooterParser.go",
        "MDParser.go",
        "ProjectParser.go",
        "YAMLParser.go",
        "error.go",
        "helper.go",
        "load.go",
        "main.go",
        "page.go",
        "redirector.go",
    ],
    importpath = "github.com/adamyi/CTFProxy/third_party/easfs",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_flosch_pongo2//:go_default_library",
        "@com_github_gholt_blackfridaytext//:go_default_library",
        "@com_github_gosimple_slug//:go_default_library",
        "@com_github_nytimes_gziphandler//:go_default_library",
        "@com_github_shirou_gopsutil//load:go_default_library",
        "@in_gopkg_russross_blackfriday_v2//:go_default_library",
        "@in_gopkg_yaml_v2//:go_default_library",
    ],
)

go_binary(
    name = "easfs",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
