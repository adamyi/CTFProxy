load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_library", "jsonnet_to_json")
load("//:config.bzl", "CONTAINER_REGISTRY", "CTF_DOMAIN")

exports_files(["cli-static-sffe.jsonnet"])

jsonnet_library(
    name = "utils",
    srcs = [
        "utils.libsonnet",
    ],
    visibility = ["//visibility:public"],
)

jsonnet_library(
    name = "services",
    srcs = [
        "services.libsonnet",
    ],
    visibility = ["//visibility:public"],
    deps = [
        ":utils",
        "//challenges:challenges_jsonnet",
        "//infra:infra_jsonnet",
    ],
)

jsonnet_to_json(
    name = "k8s",
    src = "k8s.jsonnet",
    outs = ["k8s.yaml"],
    ext_strs = {
        "container_registry": CONTAINER_REGISTRY,
        "ctf_domain": CTF_DOMAIN,
    },
    yaml_stream = True,
    deps = [
        ":services",
        ":utils",
    ],
)

jsonnet_to_json(
    name = "route53",
    src = "route53.jsonnet",
    outs = ["route53.json"],
    ext_strs = {
        "ctf_domain": CTF_DOMAIN,
    },
    deps = [
        ":services",
        ":utils",
    ],
)

jsonnet_to_json(
    name = "all-docker-compose",
    src = "docker-compose.jsonnet",
    outs = ["all-docker-compose.json"],
    ext_strs = {
        "cluster": "all",
        "container_registry": CONTAINER_REGISTRY,
        "ctf_domain": CTF_DOMAIN,
    },
    deps = [
        ":services",
        ":utils",
    ],
)

jsonnet_to_json(
    name = "cluster-master-docker-compose",
    src = "docker-compose.jsonnet",
    outs = ["cluster-master-docker-compose.json"],
    ext_strs = {
        "cluster": "master",
        "container_registry": CONTAINER_REGISTRY,
        "ctf_domain": CTF_DOMAIN,
    },
    deps = [
        ":services",
        ":utils",
    ],
)

jsonnet_to_json(
    name = "cluster-team-docker-compose",
    src = "docker-compose.jsonnet",
    outs = ["cluster-team-docker-compose.json"],
    ext_strs = {
        "cluster": "team",
        "container_registry": CONTAINER_REGISTRY,
        "ctf_domain": CTF_DOMAIN,
    },
    deps = [
        ":services",
        ":utils",
    ],
)
