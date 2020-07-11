"""
a list of challenge
author: adamyi
"""

load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_library")

def challenges_list(name, deps, visibility):
    """a list of challenges

    Args:
      name: name
      deps: dependencies
      visibility: visibility
    """
    cmd = "./$(location //tools:generatechallengeslist.py)"
    for src in deps:
        cmd += " " + src
    cmd += "> \"$@\""

    native.genrule(
        name = name + "_file",
        outs = ["challenges.libsonnet"],
        srcs = deps,
        tools = ["//tools:generatechallengeslist.py"],
        cmd = cmd,
    )

    jsonnet_library(
        name = name,
        srcs = [
            ":" + name + "_file",
        ],
        deps = deps,
        visibility = visibility,
    )
