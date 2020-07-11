"""
register all config of a ctf challenge
author: adamyi
"""

load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_library", "jsonnet_to_json")

def ctf_challenge():
    jsonnet_library(
        name = "challenge",
        srcs = [
            "challenge.libsonnet",
        ],
        visibility = ["//challenges:__pkg__", "//infra:__pkg__"],
    )
    jsonnet_to_json(
        name = "clisffe",
        src = "//infra/jsonnet:cli-static-sffe.jsonnet",
        outs = ["clisffe.json"],
        tla_code_files = {
            "challenge.libsonnet": "challenge",
        },
    )
