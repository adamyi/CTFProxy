# CTFProxy

Your ultimate CTF infrastructure, with a BeyondCorp-like Zero-Trust Network and simple infrastructure-as-code configuration.

This was used for [UNSW's Web Application Security course (COMP6443 and COMP6843)](https://webcms3.cse.unsw.edu.au/COMP6443/20T2/), for which we are running 50+ challenges with 100+ containers. Only a single command is used to build all the containers (written in Golang, Python, PHP, Jsonnet, etc.). We also have CI server setup to perform auto linting, compiling, pushing, and deploying. We have a test environment running docker-compose and a prod environment running Kubernetes. All configurations live in a single monorepo and do not require manual (undeterministic) work.

Players are SSO-authenticated with either mTLS certificate or username/password. This authentication info is then passed as JWT to internal services. Inter-service communication is also authenticated via IP address and signed with JWT tokens. All challenges and infra services have a single configuration file called `challenge.libsonnet` that defines its domain names, service images, replicas, health check, flags, and access control policies. Access control policies are defined in [Starlark](https://github.com/bazelbuild/starlark), a dialect of Python to allow really fine-grained access control.

This is a fork of our original [Geegle3](https://github.com/adamyi/Geegle3)

## Infra Microservices

- ctfd: CTFd scoreboard (it doesn't require any manual setup for challenges. It's calling the flaganizer microservice)
- ctfproxy: TLS (HTTPS) termination, SSO, Access Control, Domain Routing, Centralized Logging
- dns: this is used only for `docker` deployment (not used with `k8s`) to help with inter-service communication
- elk: elasticsearch stack for ctfproxy logging
- flaganizer: flag signing and verification service (supports unique flag per player)
- gaia: internal user authentication and groups microservice
- isodb: isolated database as a service
- requestz: a simple network debugging tool
- whoami: shows you who you are
- xssbot: headless chrome as a service (puppeteer queue), integrated with ctfproxy auth

We are, by design, using an external MySQL database out of the cluster so that we won't risk losing the data in a container.

It's a TODO to add PostgreSQL support to isodb, so we can use the same isodb RPC interface to connect to different backends
for various SQL injection challenges.

## Installation

### Configuration

First fork this repo!

You need to configure your domain name and container registry in [config.bzl](config.bzl).

You also need to add a few commandline arguments to infra services as configuration. You can grep for `FIXME` to find them

```
$ grep -r FIXME infra
infra/ctfproxy/BUILD.bazel:        # FIXME: configure ssl_cert, ssl_key, mtls_ca here
infra/ctfproxy/BUILD.bazel:        # FIXME: configure cert_gcs_bucket gcp_service_account gcp_project here
infra/flaganizer/BUILD.bazel:        # FIXME: configure flag_key here
infra/gaia/BUILD.bazel:        # FIXME: configure dbaddr here
infra/isodb/BUILD.bazel:        # FIXME: configure dbaddr here
infra/ctfd/BUILD.bazel:        # FIXME: configure DATABASE_URL and SECRET_KEY here
```

### Challenge Creation

TODO(adamyi): better documentation

Every challenge lives under its own directory under `challenges/` with a single configuration file `challenges.libsonnet`

An example configuration looks like this:

```
{
  services: [
    {
      name: 'welcome',
      replicas: 3,
      category: 'week0',
      clustertype: 'master',
      access: 'grantAccess()',
    },
  ],
  flags: [
    {
      Id: 'welcome_flag_dynamic',  // globally unique flag id among entire monorepo
      DisplayName: 'Welcome!',  // a display name of the flag in submission history visible to players on CTFd
      Category: 'misc',  // a challenge category visible to players on CTFd
      Points: 1,  // points associated with the flag
      Type: 'dynamic',  // unique flag per user
      Flag: 'WELCOME_TO_CTFPROXY',  // the flag will look like FLAG{WELCOME_TO_CTFPROXY.XXX.YYY}
    },
    {
      Id: 'welcome_flag_fixed',  // globally unique flag id among entire monorepo
      DisplayName: 'Welcome!',  // a display name of the flag in submission history visible to players on CTFd
      Category: 'misc',  // a challenge category visible to players on CTFd
      Points: 1,  // points associated with the flag
      Type: 'fixed',  // fixed static flag
      Flag: 'WELCOME2CTFPROXY',  // the flag will look like FLAG{WELCOME2CTFPROXY}
    },
  ],
}
```

This is all you need for your challenge - no need to configure CTFd manually. The challenge will be added lazily on CTFd.

The structure of your challenges are going to look very similar to the infra services under [infra](infra) directory. In `BUILD.bazel` file, you need to define a `:image` target using [rules_docker](https://github.com/bazelbuild/rules_docker).

While you can still use Dockerfile (there's a `dockerfile_image` rule in bazel rules_docker), try using `container_image` instead. Unlike Dockerfile, this ensures that the building process is actually deterministic/reproduceable. Use all existing infra/challenges BUILD files as reference, also learn more about it here: https://github.com/bazelbuild/rules_docker/tree/master/testing/examples.

After creating your challenge, run `bazel run //tools/ctflark` to update `challenges_list` and `container_bundle` (this is automatic, you just need to run it)

User authentication info is passed to challenges in `X-CTFProxy-JWT` header. Response headers starting with `X-CTFProxy-I-` are only shown to staff, not to players. They also appear in ELK logs. Use this for better logging/debugging.

Please have a `/healthz` endpoint (this path is also configurable) on your challenge returning a 200. This is for load balancing and health checking purpose. Your instance will be killed and respawned if health check fails in k8s.

### Domain Routing

Domain routing is handled automatically by CTFProxy. All you need is set up the `challenge.libsonnet` file. Specifically, your domain will be `{service_name}.{ctf_domain}`. Multi-level subdomains are also supported. E.g., if my service is named `dev-dot-blog`, and my ctf_domain is configured as `adamyi.com`. The service will live on `https://dev.blog.adamyi.com`.

It's recommended that you set up a wildcard record to CTFProxy so you don't need to worry about the rest (you get a nice error page if the domain doesn't exist). If you prefer to only have records for existing domains, we have our infra to set up AWS Route53 configuration. Just run `bazel build //infra/jsonnet:route53`. Please send a pull request if you added integration with other cloud providers.

TLS (HTTPS) is terminated by CTFProxy. You don't need to set up anything manually. It automatically provisions your certificate using ACME protocol with Let's Encrypt. If you wish, you can also supply your own certificates by providing filepath to ctfproxy's commandline argument. Edit `infra/ctfproxy/BUILD.bazel` for this. Auto certs are propagated with GCS (Google Cloud Storage) among replicas.

### Access Control

As mentioned above, you are using Starlark for access control policy. A more complex example would be:

```
      access: |||
        def checkAccess():
          if method != "GET" and path != "/wp-login.php":
            denyAccess()
          if path.startswith("/wp-admin/plugin") or path.startswith("/wp-admin/option") or path.startswith("/wp-admin/theme"):
            denyAccess()
          if path in ["/wp-admin/customize.php", "/wp-admin/widgets.php", "/wp-admin/nav-menus.php", "/wp-admin/site-health.php", "/wp-admin/tools.php", "/wp-admin/import.php", "/wp-admin/export.php", "/wp-admin/erase-personal-data.php", "/wp-admin/export-personal-data.php", "/wp-admin/update-core.php"]:
            denyAccess()
          if path == "/wp-admin/users.php" and len(query) > 0:
            denyAccess()
          # only staff has access for now, because it's not yet released
          if ("staff@groups." + corpDomain) in groups:
            grantAccess()
        checkAccess()
      |||,
```

In this example, we made a shared Wordpress instance's admin portal safe in a CTF environment.

Specifically, you can call the following functions:

- denyAccess() - deny access
- grantAccess() - grant access provided that the user has logged in
- openAccess() - grant access no matter if the user is logged in

You also have access to the following variables:

- host: request hostname
- method: request method
- path: request path
- rawpath: request path (url-encoded raw path)
- query: raw query string
- ip: remote IP address
- user: SSO username
- corpDomain: your configured CTF domain
- groups: SSO groups that user is in
- timestamp: current unix milli epoch timestamp

### Building

Please build using Linux AMD64. Cuz it's hard to set up cross-compiling for C programs on mac, ceebs.

Build only:

```
bazel build //:all_containers
```

Build and tag locally (so that you can use docker-compose to boot them up):

```
bazel run //:all_containers
```

### Deploying (kubernetes, alpha)

We support both kubernetes and docker-compose!

While the recommended practice is to use k8s, docker-compose has not been deprecated yet. This means that we shouldn't introduce any breaking changes that rely solely on k8s api. Everything should be able to run on just docker.

The approach we adopted for COMP6443 is that we use docker-compose for local dev & testing, because it's simple to set up. Use k8s cluster for prod high-availability deployment.

```
bazel build //infra/jsonnet:k8s
kubectl apply -f dist/bin/infra/k8s.yaml
```

**Rationale for Switching from Docker-Compose to K8S**

Our manual approach for managing IP addresses works great, but we can only add service to the end of list so it gets a new ip range. Otherwise, we're changing ip addresses for existing docker networks and this requires a full reboot with downtime.

k8s gives us better service discovery and can maintain zero-downtime for adding and updating services. k8s also gives us rolling updates, replicas, scaling, etc. Also, for isolated challenges, we no longer need one server per student in k8s infra (this is not yet supported with k8s deployment and is a TODO).

The downside of k8s is that it's such a pain to set up a working cluster... But after it's set up, everything just works (as long as you know what you're doing).

### Deploying (docker-compose, legacy stable)

#### Main Server (Shared Server)

```
bazel build //infra/jsonnet:cluster-master-docker-compose
docker-compose -f dist/bin/infra/jsonnet/cluster-master-docker-compose.json up -d
```

#### Team Server (Separate Isolated Server)

```
bazel build //infra/jsonnet:cluster-team-docker-compose
docker-compose -f dist/bin/infra/jsonnet/cluster-team-docker-compose.json up -d
```

#### Test Server (All-in-one Server)

```
bazel build //infra/jsonnet:all-docker-compose
docker-compose -f dist/bin/infra/jsonnet/all-docker-compose.json up -d
```

## Coding Style

- Run `bazel run //tools/ctflark` to update `challenges_list` and `container_bundle` for challenges.
- For golang, we use the go standard style. Please use `bazel run @go_sdk//:bin/gofmt -- -s -d .` to format your code. We use `nogo` to lint your code at compile-time. Use `bazel run //:gazelle` to generate BUILD files.
- For python, we use yapf style, which is like PEP8, but uses 2-space indentation instead of 4. Please use `yapf --in-place --recursive --parallel .` to format your code. Use `pylint` as a linter (not yet enabled due to failed linting). Please ensure code is python2-compatible.
- For jsonnet, please use `tools/format_jsonnet.sh` to format your code.
- For bazel, use `bazel run //:buildifier` to format BUILD files.

There's a nice `tools/format.sh` to run all commands above.

You can run `tools/PRESUBMIT.sh` to check if your code meets all style requirements and compiles correctly.

## Contributing

You're welcome to contribute! Just send a pull request.

## Author

[Adam Yi](https://github.com/adamyi)

## LICENSE

Open-sourced with love, under [Apache 2.0 License](LICENSE).
