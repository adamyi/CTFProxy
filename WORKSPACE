workspace(
    name = "ctfproxy",
    managed_directories = {"@npm": ["node_modules"]},
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "a8d6b1b354d371a646d2f7927319974e0f9e52f73a2452d2b3877118169eb6bb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(nogo = "@//:nogo")

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "a0e79f5876a1552ae8000882e4189941688f359a80b2bc1d7e3a51cab6257ba1",
    strip_prefix = "buildtools-3.0.0",
    url = "https://github.com/bazelbuild/buildtools/archive/3.0.0.tar.gz",
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "3cee35dcfbebd224f500f77210f99986a0d303ccf9ce006e3fe0dcf14d50689d",
    strip_prefix = "rules_docker-a012c6d7ab8d58417f4b805a61cbd397248c9f21",
    urls = ["https://github.com/adamyi/rules_docker/archive/a012c6d7ab8d58417f4b805a61cbd397248c9f21.tar.gz"],
)

# local_repository(
#     name = "io_bazel_rules_docker",
#     path = "../rules_docker",
# )

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//repositories:pip_repositories.bzl", container_pip_deps = "pip_deps")

container_pip_deps()

http_archive(
    name = "bazel_gazelle",
    sha256 = "bfd86b3cbe855d6c16c6fce60d76bd51f5c8dbc9cfcaef7a2bb5c1aafd0710e8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.0/bazel-gazelle-v0.21.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.0/bazel-gazelle-v0.21.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

http_archive(
    name = "com_google_protobuf",
    sha256 = "6f3e1a448af71b2e98f1f38d25dcc07bd7c008eea03fec4b6c9a0ea2bfc1778f",
    strip_prefix = "protobuf-3.12.0",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v3.12.0/protobuf-all-3.12.0.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "rules_python",
    sha256 = "b5668cde8bb6e3515057ef465a35ad712214962f0b3a314e551204266c7be90c",
    strip_prefix = "rules_python-0.0.2",
    url = "https://github.com/bazelbuild/rules_python/releases/download/0.0.2/rules_python-0.0.2.tar.gz",
)

RULES_NODEJS_VERSION = "1.7.0"

RULES_NODEJS_SHA256 = "84abf7ac4234a70924628baa9a73a5a5cbad944c4358cf9abdb4aab29c9a5b77"

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = RULES_NODEJS_SHA256,
    url = "https://github.com/bazelbuild/rules_nodejs/releases/download/%s/rules_nodejs-%s.tar.gz" % (RULES_NODEJS_VERSION, RULES_NODEJS_VERSION),
)

http_archive(
    name = "io_bazel_rules_closure",
    sha256 = "9161f3b719008b223846b0df63c7674c6e2d67c81e052a9864f90736505c35f3",
    strip_prefix = "rules_closure-62746bdd1087c1198a81143e7d8ef3d144a43c0f",
    urls = [
        "https://github.com/bazelbuild/rules_closure/archive/62746bdd1087c1198a81143e7d8ef3d144a43c0f.tar.gz",
    ],
)

load("@io_bazel_rules_closure//closure:repositories.bzl", "rules_closure_dependencies", "rules_closure_toolchains")

rules_closure_dependencies()

rules_closure_toolchains()

http_archive(
    name = "com_github_grpc_grpc_web",
    sha256 = "283fdea0ff1539f47315700fc557837d539cbaa2b6e5e5dcb0b7280d3f054029",
    strip_prefix = "grpc-web-6b99a37519fd5de2d46f7fd4d1d293504e15161f",
    urls = [
        "https://github.com/grpc/grpc-web/archive/6b99a37519fd5de2d46f7fd4d1d293504e15161f.tar.gz",
    ],
)

load("@io_bazel_rules_docker//container:pull.bzl", "container_pull")
load("@io_bazel_rules_docker//go:image.bzl", _go_image_repos = "repositories")
load("@io_bazel_rules_docker//java:image.bzl", _java_image_repos = "repositories")
load("@io_bazel_rules_docker//nodejs:image.bzl", _nodejs_image_repos = "repositories")
load("@io_bazel_rules_docker//python:image.bzl", _py_image_repos = "repositories")

_java_image_repos()

_go_image_repos()

_py_image_repos()

_nodejs_image_repos()

git_repository(
    name = "jsonnet_go",
    commit = "8a0084e643955269558e206eb9f4d90e2e569781",
    init_submodules = True,
    remote = "https://github.com/google/go-jsonnet",
)

http_archive(
    name = "io_bazel_rules_jsonnet",
    sha256 = "7f51f859035cd98bcf4f70dedaeaca47fe9fbae6b199882c516d67df416505da",
    strip_prefix = "rules_jsonnet-0.3.0",
    urls = ["https://github.com/bazelbuild/rules_jsonnet/archive/0.3.0.tar.gz"],
)

load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_repositories")

jsonnet_repositories()

load("@jsonnet_go//bazel:repositories.bzl", "jsonnet_go_repositories")

jsonnet_go_repositories()

load("@jsonnet_go//bazel:deps.bzl", "jsonnet_go_dependencies")

jsonnet_go_dependencies()

http_archive(
    name = "base_images_docker",
    sha256 = "be6043d38aa7fad421910311aecec865beb060eb56d8c3eb5af62b2805e9379c",
    strip_prefix = "base-images-docker-3eed1bbda3e500f72b36745c9d74385d82ca1b19",
    urls = ["https://github.com/GoogleCloudPlatform/base-images-docker/archive/3eed1bbda3e500f72b36745c9d74385d82ca1b19.tar.gz"],
)

container_pull(
    name = "tomcat9",
    registry = "index.docker.io",
    repository = "tomcat",
    tag = "9.0.21-jdk8",
)

container_pull(
    name = "wordpress",
    registry = "index.docker.io",
    repository = "wordpress",
    tag = "php7.4-apache",
)

container_pull(
    name = "tomcat-jython",
    digest = "sha256:27526ffde703e09cdf8cbb3cb781c169ac48f2e2ba3a6fbe3238c9fff9b80fc7",
    registry = "index.docker.io",
    repository = "adamyi/tomcat-jython",
)

container_pull(
    name = "python-with-latex",
    digest = "sha256:db92134dd530dd3b666a5b420029886c01faf72a7ea366726c6eeac45ae4ed64",
    registry = "index.docker.io",
    repository = "adamyi/python-with-latex",
)

container_pull(
    name = "nginx-php-fpm-with-imagick",
    digest = "sha256:2ac175a4b6faff45ca12325de3ed3899c8acba2d38955fab5b8b877c8cb7c6d5",
    registry = "index.docker.io",
    repository = "adamyi/nginx-php-fpm-with-imagick",
)

container_pull(
    name = "ubuntu1804-with-32bit-libc",
    digest = "sha256:3225563499e60d3bacd4db8f05920ae5d86635372d1c77024fd73d6db9d04cca",
    registry = "index.docker.io",
    repository = "adamyi/ubuntu1804-with-32bit-libc",
)

container_pull(
    name = "ubuntu1804-with-zbar",
    digest = "sha256:cc47d8fc8309178954287c6419f3f39aa3741b6c540c351bd5d71c3662b9d6ba",
    registry = "index.docker.io",
    repository = "adamyi/ubuntu1804-with-zbar",
)

container_pull(
    name = "chrome-base-without-chrome",
    digest = "sha256:b5c86894a56352eb4f91c462d7cb95b5475b3e3735d4faea9893cfaca668c467",
    registry = "index.docker.io",
    repository = "adamyi/chrome-base-without-chrome",
)

container_pull(
    name = "nginx-php-fpm",
    digest = "sha256:2e9718f4bdca05f577cb8cf046327cb9232e4fd817fe32f470db0a65660a6e46",
    registry = "index.docker.io",
    repository = "richarvey/nginx-php-fpm",
)

container_pull(
    name = "alpine_linux_amd64",
    registry = "index.docker.io",
    repository = "library/alpine",
    tag = "3.8",
)

container_pull(
    name = "ubuntu1804",
    digest = "sha256:3942b604b2f23e9b08aa6f3c51dc19efa2b570ae93ce8aaabf94e02111ddcca9",
    registry = "gcr.io",
    repository = "cloud-marketplace/google/ubuntu1804",
)

container_pull(
    name = "python2-base",
    digest = "sha256:938d21703d929295337f5aafd038a8d93172e11e1746f6e87f02eb53e32bcea0",
    registry = "index.docker.io",
    repository = "python",
)

container_pull(
    name = "python2.7.17",
    digest = "sha256:9b7d62026be68c2e91c17fb4e0499454e41ebf498ef345f9ad6e100a67e4b697",
    registry = "index.docker.io",
    repository = "python",
)

container_pull(
    name = "python3-base",
    digest = "sha256:d182a775e372d82d92b59ff5debeabcb699964fe163320eb9fc5ebb971c51ec3",
    registry = "index.docker.io",
    repository = "python",
)

container_pull(
    name = "elasticsearch",
    digest = "sha256:59342c577e2b7082b819654d119f42514ddf47f0699c8b54dc1f0150250ce7aa",  # 7.6.2
    registry = "docker.elastic.co",
    repository = "elasticsearch/elasticsearch",
)

container_pull(
    name = "kibana",
    digest = "sha256:e8f3743e404462709663422056db2d5076a7a6bd6024f64aea1599b3014c63be",  # 7.6.2
    registry = "docker.elastic.co",
    repository = "kibana/kibana",
)

go_repository(
    name = "com_github_dgrijalva_jwt_go",
    importpath = "github.com/dgrijalva/jwt-go",
    tag = "v3.2.0",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    tag = "v1.7.3",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    tag = "v1.4.0",
)

go_repository(
    name = "com_github_tuotoo_qrcode",
    commit = "ac9c44189bf2",
    importpath = "github.com/tuotoo/qrcode",
)

go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    tag = "v1.1.1",
)

go_repository(
    name = "com_github_skip2_go_qrcode",
    commit = "dc11ecdae0a9",
    importpath = "github.com/skip2/go-qrcode",
)

go_repository(
    name = "com_github_syndtr_goleveldb",
    commit = "02440ea7a285",
    importpath = "github.com/syndtr/goleveldb",
)

go_repository(
    name = "com_github_mattn_go_sqlite3",
    importpath = "github.com/mattn/go-sqlite3",
    tag = "v1.10.0",
)

go_repository(
    name = "com_github_go_sql_driver_mysql",
    importpath = "github.com/go-sql-driver/mysql",
    tag = "v1.4.1",
)

go_repository(
    name = "com_github_joewalnes_websocketd",
    importpath = "github.com/joewalnes/websocketd",
    tag = "v0.3.1",
)

go_repository(
    name = "com_github_miekg_dns",
    importpath = "github.com/miekg/dns",
    tag = "v1.1.22",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:cg5LA/zNPRzIXIWSCxQW10Rvpy94aQh3LT/ShoCpkHw=",
    version = "v0.0.0-20200510223506-06a226fb4e37",
)

go_repository(
    name = "org_golang_x_net",
    commit = "aa69164e4478",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_sync",
    commit = "112230192c58",
    importpath = "golang.org/x/sync",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "2837fb4f24fe",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    tag = "v0.3.2",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "2ca718005c18",
    importpath = "golang.org/x/tools",
)

go_repository(
    name = "com_github_jhillyerd_enmime",
    importpath = "github.com/jhillyerd/enmime",
    tag = "v0.6.0",
)

go_repository(
    name = "com_github_ssor_bom",
    commit = "6386211fdfcf",
    importpath = "github.com/ssor/bom",
)

go_repository(
    name = "net_starlark_go",
    importpath = "go.starlark.net",
    sum = "h1:S2s+dYPyDg/vF7KbcRIB2831xVimJoR4zebfoVBzn7Q=",
    version = "v0.0.0-20200330013621-be5394c419b6",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    sum = "h1:EpMNVUorLiZIELdMZbCYX/ByTFCdoYopYAGxaGVz9ms=",
    version = "v0.57.0",
)

go_repository(
    name = "com_google_cloud_go_storage",
    importpath = "cloud.google.com/go/storage",
    sum = "h1:86K1Gel7BQ9/WmNWn7dTKMvTLFzwtBe5FNqYbi9X35g=",
    version = "v1.8.0",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:TzXSXBo42m9gQenoE3b9BGiEpg5IG2JkU5FkPIawgtw=",
    version = "v0.0.0-20200107190931-bf48bf16ab8d",
)

go_repository(
    name = "org_golang_google_api",
    importpath = "google.golang.org/api",
    sum = "h1:cG03eaksBzhfSIk7JRGctfp3lanklcOM/mTGvow7BbQ=",
    version = "v0.24.0",
)

go_repository(
    name = "com_github_googleapis_gax_go_v2",
    importpath = "github.com/googleapis/gax-go/v2",
    sum = "h1:sjZBwGj9Jlw33ImPtvFviGYvseOtDM7hkSKB7+Tv3SM=",
    version = "v2.0.5",
)

go_repository(
    name = "io_opencensus_go",
    importpath = "go.opencensus.io",
    sum = "h1:8sGtKOrtQqkN1bp2AtX+misvLIlOmsEsNd+9NIcPEm8=",
    version = "v0.22.3",
)

go_repository(
    name = "com_github_golang_groupcache",
    importpath = "github.com/golang/groupcache",
    sum = "h1:1r7pUrabqp18hOBcwBwiTsbnFeTZHV9eER/QT5JVZxY=",
    version = "v0.0.0-20200121045136-8c9f03a8e57e",
)

go_repository(
    name = "org_golang_google_grpc",
    importpath = "google.golang.org/grpc",
    sum = "h1:M5a8xTlYTxwMn5ZFkwhRabsygDY5G8TYLyQDBxJNAxE=",
    version = "v1.30.0",
)

go_repository(
    name = "com_github_adamyi_hotconfig",
    commit = "59069be03b90",
    importpath = "github.com/adamyi/hotconfig",
)

load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

yarn_install(
    name = "npm",
    package_json = "//:package.json",
    yarn_lock = "//:yarn.lock",
)

load("@rules_python//python:pip.bzl", "pip_import")

http_archive(
    name = "chromium",
    build_file = "@//third_party:chromium.BUILD",
    sha256 = "10ae4e05d9f01a8b646dd2ccc2ac1135e597c472abe5be71552aae7d8a35e2ac",
    url = "https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Linux_x64%2F650583%2Fchrome-linux.zip?generation=1555131417316559&alt=media",
)

http_archive(
    name = "ctfd",
    build_file = "@//third_party/ctfd:ctfd.BUILD",
    sha256 = "31a0046dba807a66b637ff82e53ed99159b27c6fc5ac5d9ed78b59e68fb43600",
    strip_prefix = "CTFd-b3fd3832d5ab7484abe52dcc84e15a93374b617b",
    url = "https://github.com/secedu/CTFd/archive/b3fd3832d5ab7484abe52dcc84e15a93374b617b.tar.gz",
)

# Local copy of CTFd, for local debugging purposes
#new_local_repository(
#    name = "ctfd",
#    build_file = "@//third_party/ctfd:ctfd.BUILD",
#    path = "../CTFd",
#)

pip_import(
    name = "ctfd_pip",
    requirements = "@ctfd//:requirements.txt",
)

load(
    "@ctfd_pip//:requirements.bzl",
    _ctfd_install = "pip_install",
)

_ctfd_install()

pip_import(
    name = "ctfproxy_pip",
    requirements = "//:requirements.txt",
)

load(
    "@ctfproxy_pip//:requirements.bzl",
    _ctfproxy_install = "pip_install",
)

_ctfproxy_install()

go_repository(
    name = "com_github_bmatcuk_doublestar",
    importpath = "github.com/bmatcuk/doublestar",
    sum = "h1:1jLE2y0VpSrOn/QR9G4f2RmrCtkM3AuATcWradjHUvM=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_flosch_pongo2",
    importpath = "github.com/flosch/pongo2",
    sum = "h1:GY1+t5Dr9OKADM64SYnQjw/w99HMYvQ0A8/JoUkxVmc=",
    version = "v0.0.0-20190707114632-bbf5a6c351f4",
)

go_repository(
    name = "com_github_gholt_blackfridaytext",
    importpath = "github.com/gholt/blackfridaytext",
    sum = "h1:aWeuOsmyHzAuZvekBl4pnJgJCtYLnc7X5JlCQocUros=",
    version = "v0.0.0-20190816214545-16f7b9b9742e",
)

go_repository(
    name = "com_github_gholt_brimtext",
    importpath = "github.com/gholt/brimtext",
    sum = "h1:OfEy3A+F4fmU2ZgBd6fBJ03gR6Kw5euUbs5tpGXD/6U=",
    version = "v0.0.0-20190811231012-1fbdf4665642",
)

go_repository(
    name = "com_github_go_check_check",
    importpath = "github.com/go-check/check",
    sum = "h1:0gkP6mzaMqkmpcJYCFOLkIBwI7xFExG03bbkOkCvUPI=",
    version = "v0.0.0-20180628173108-788fd7840127",
)

go_repository(
    name = "com_github_gosimple_slug",
    importpath = "github.com/gosimple/slug",
    sum = "h1:BlCZq+BMGn+riOZuRKnm60Fe7+jX9ck6TzzkN1r8TW8=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_juju_errors",
    importpath = "github.com/juju/errors",
    sum = "h1:rhqTjzJlm7EbkELJDKMTU7udov+Se0xZkWmugr6zGok=",
    version = "v0.0.0-20181118221551-089d3ea4e4d5",
)

go_repository(
    name = "com_github_juju_loggo",
    importpath = "github.com/juju/loggo",
    sum = "h1:MK144iBQF9hTSwBW/9eJm034bVoG30IshVm688T2hi8=",
    version = "v0.0.0-20180524022052-584905176618",
)

go_repository(
    name = "com_github_juju_testing",
    importpath = "github.com/juju/testing",
    sum = "h1:WQM1NildKThwdP7qWrNAFGzp4ijNLw8RlgENkaI4MJs=",
    version = "v0.0.0-20180920084828-472a3e8b2073",
)

go_repository(
    name = "com_github_kr_pretty",
    importpath = "github.com/kr/pretty",
    sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_kr_pty",
    importpath = "github.com/kr/pty",
    sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_kr_text",
    importpath = "github.com/kr/text",
    sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_mattn_goveralls",
    importpath = "github.com/mattn/goveralls",
    sum = "h1:7eJB6EqsPhRVxvwEXGnqdO2sJI0PTsrWoTMXEk9/OQc=",
    version = "v0.0.2",
)

go_repository(
    name = "com_github_rainycape_unidecode",
    importpath = "github.com/rainycape/unidecode",
    sum = "h1:ta7tUOvsPHVHGom5hKW5VXNc2xZIkfCKP8iaqOyYtUQ=",
    version = "v0.0.0-20150907023854-cb7f23ec59be",
)

go_repository(
    name = "com_github_russross_blackfriday",
    importpath = "github.com/russross/blackfriday",
    sum = "h1:hgS5QyP981zzGr3UYaoHb5+fpgK1lHleAOq5znvfJxU=",
    version = "v0.0.0-20171011182219-6d1ef893fcb0",
)

go_repository(
    name = "com_github_shurcool_sanitized_anchor_name",
    importpath = "github.com/shurcooL/sanitized_anchor_name",
    sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
    version = "v1.0.0",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
    version = "v0.0.0-20161208181325-20d25e280405",
)

go_repository(
    name = "in_gopkg_mgo_v2",
    importpath = "gopkg.in/mgo.v2",
    sum = "h1:xcEWjVhvbDy+nHP67nPDDpbYrY+ILlfndk4bRioVHaU=",
    version = "v2.0.0-20180705113604-9856a29383ce",
)

go_repository(
    name = "in_gopkg_russross_blackfriday_v2",
    importpath = "gopkg.in/russross/blackfriday.v2",
    sum = "h1:+FlnIV8DSQnT7NZ43hcVKcdJdzZoeCmJj4Ql8gq5keA=",
    version = "v2.0.0",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:ZCJp+EgiOT7lHqUV2J862kp8Qj64Jo6az82+3Td9dZw=",
    version = "v2.2.2",
)

go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
    version = "v0.0.0-20160126235308-23def4e6c14b",
)

go_repository(
    name = "com_github_shirou_gopsutil",
    importpath = "github.com/shirou/gopsutil",
    sum = "h1:1eaJvGomDnH74/5cF4CTmTbLHAriGFsTZppLXDX93OM=",
    version = "v2.18.12+incompatible",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:ZDRjVQ15GmhC3fiQ8ni8+OwkZQO4DARzQgrnXU1Liz8=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_nytimes_gziphandler",
    importpath = "github.com/NYTimes/gziphandler",
    sum = "h1:ZUDjpQae29j0ryrS0u/B8HZfJBtBQHjqw2rQ2cqUQ3I=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:TivCn/peBQ7UY8ooIcPgZFpTNSz0Q2U6UrFlUfqbe0Q=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_bdwilliams_go_jsonify",
    importpath = "github.com/bdwilliams/go-jsonify",
    sum = "h1:PnDia1dKFSh7KDMoe9mRkSEFAIX3yN36Kc+zf/hLQVA=",
    version = "v0.0.0-20141020182238-48749139e742",
)

go_repository(
    name = "com_github_ulule_limiter_v3",
    importpath = "github.com/ulule/limiter/v3",
    sum = "h1:QRAebbswjlezHIfiSQgM8+jMxaz/zsrxGRuiUJ43MHo=",
    version = "v3.5.0",
)

go_repository(
    name = "com_github_improbable_eng_grpc_web",
    importpath = "github.com/improbable-eng/grpc-web",
    sum = "h1:7XqtaBWaOCH0cVGKHyvhtcuo6fgW32Y10yRKrDHFHOc=",
    version = "v0.13.0",
)

go_repository(
    name = "com_github_desertbit_timer",
    importpath = "github.com/desertbit/timer",
    sum = "h1:U5y3Y5UE0w7amNe7Z5G/twsBW0KEalRQXZzf8ufSh9I=",
    version = "v0.0.0-20180107155436-c41aec40b27f",
)

go_repository(
    name = "com_github_rs_cors",
    importpath = "github.com/rs/cors",
    sum = "h1:+88SsELBHx5r+hZ8TCkggzSstaWNbDvThkVK8H6f9ik=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_rpcxio_gomemcached",
    importpath = "github.com/rpcxio/gomemcached",
    sum = "h1:/jGV5cu7zNStDo+A5r3/CHIBH2ENgvw59fM2wgzFtQE=",
    version = "v0.0.0-20200223142310-2dc6f77e072e",
)

go_repository(
    name = "com_github_yeka_zip",
    importpath = "github.com/yeka/zip",
    sum = "h1:OJYP70YMddlmGq//EPLj8Vw2uJXmrA+cGSPhXTDpn2E=",
    version = "v0.0.0-20180914125537-d046722c6feb",
)
