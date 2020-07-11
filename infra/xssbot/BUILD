load("@io_bazel_rules_docker//container:container.bzl", "container_image")
load("@io_bazel_rules_docker//nodejs:image.bzl", "nodejs_image")
load("//:config.bzl", "CTF_DOMAIN")
load("//tools:challenge.bzl", "ctf_challenge")

ctf_challenge()

container_image(
    name = "image",
    base = ":nodejs_image",
    env = {
        "CTFDOMAIN": CTF_DOMAIN,
        "DEBUG": "puppeteer-cluster:*",
        "MAXCONCURRENTY": "5",
        "PORT": "80",
        "TASKTIMEOUT": "90000",
    },
    files = [
        "//jwtkeys:jwtRS256.key.pub",
    ],
    tars = [
        "@chromium",
    ],
    visibility = ["//visibility:public"],
)

nodejs_image(
    name = "nodejs_image",
    base = "@chrome-base-without-chrome//image",
    data = [":server.js"],
    entry_point = "server.js",
    node_modules = "@npm//:node_modules",
)
