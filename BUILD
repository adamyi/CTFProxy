load("@bazel_gazelle//:def.bzl", "gazelle")
load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")
load("@io_bazel_rules_docker//container:container.bzl", "container_bundle")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "container_push")
load("@io_bazel_rules_go//go:def.bzl", "TOOLS_NOGO", "nogo")
load("@rules_python//python:defs.bzl", "py_library")
load("//:config.bzl", "CTF_DOMAIN")

# exports_files(["tsconfig.json"])

# gazelle:exclude dist
# gazelle:exclude node_modules
# gazelle:prefix github.com/adamyi/CTFProxy
gazelle(name = "gazelle")

buildifier(
    name = "buildifier",
    exclude_patterns = [
        "./dist/*",
        "./node_modules/*",
    ],
    lint_mode = "fix",
    lint_warnings = ["all"],
)

buildifier(
    name = "buildifier_check",
    exclude_patterns = [
        "./dist/*",
        "./node_modules/*",
    ],
    lint_mode = "warn",
    lint_warnings = ["all"],
    mode = "check",
)

nogo(
    name = "nogo",
    config = "//tools:nogoconfig.json",
    visibility = ["//visibility:public"],
    deps = TOOLS_NOGO,
)

container_bundle(
    name = "all_containers",
    images = {
        "gcr.io/ctfproxy/elk/elasticsearch:latest": "//infra/elk:elasticsearch",  #ctflark: keep
        "gcr.io/ctfproxy/elk/kibana:latest": "//infra/elk:kibana",  #ctflark:keep
        "gcr.io/ctfproxy/infra/ctfd:latest": "//infra/ctfd:image",
        "gcr.io/ctfproxy/infra/ctfproxy:latest": "//infra/ctfproxy:image",
        "gcr.io/ctfproxy/infra/dns:latest": "//infra/dns:image",
        "gcr.io/ctfproxy/infra/flaganizer:latest": "//infra/flaganizer:image",
        "gcr.io/ctfproxy/infra/gaia:latest": "//infra/gaia:image",
        "gcr.io/ctfproxy/infra/isodb:latest": "//infra/isodb:image",
        "gcr.io/ctfproxy/infra/requestz:latest": "//infra/requestz:image",
        "gcr.io/ctfproxy/infra/whoami:latest": "//infra/whoami:image",
        "gcr.io/ctfproxy/infra/xssbot:latest": "//infra/xssbot:image",
    },
)

container_push(
    name = "all_containers_push",
    bundle = ":all_containers",
    format = "Docker",
)

genrule(
    name = "ctf_domain_py",
    srcs = [],
    outs = ["ctf_domain.py"],
    cmd = "echo CTF_DOMAIN=\\\"" + CTF_DOMAIN + "\\\"> \"$@\"",
)

py_library(
    name = "python_ctf_domain",
    srcs = [":ctf_domain.py"],
    visibility = ["//visibility:public"],
)
