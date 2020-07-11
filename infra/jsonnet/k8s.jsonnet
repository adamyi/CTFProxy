local services = import 'infra/jsonnet/services.libsonnet';
local utils = import 'infra/jsonnet/utils.libsonnet';

local image(service) = if 'image' in service then
  service.image
else if service.category == 'infra' then
  '%s/infra/%s:latest' % [std.extVar('container_registry'), service.name]
else
  '%s/challenges/%s/%s:latest' % [std.extVar('container_registry'), service.category, service.name];

local kDeployment(service) = {
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: {
    name: service.name,
    labels: {
      app: service.name,
    },
  },
  spec: {
    replicas: if 'replicas' in service then service.replicas else 1,
    selector: {
      matchLabels: {
        app: service.name,
      },
    },
    template: {
      metadata: {
        labels: {
          app: service.name,
          zerotrust: 'ctfproxy',
        },
      },
      spec: {
        automountServiceAccountToken: false,
        enableServiceLinks: false,
        volumes: if 'persistent' in service then [{
          name: 'data',
          persistentVolumeClaim: {
            claimName: service.name + '-pv-claim',
          },
        }] else [],
        containers: [
          {
            name: service.name,
            image: image(service),
            volumeMounts: if 'persistent' in service then [{ name: 'data', mountPath: '/data' }] else [],
            ports: [
              {
                containerPort: 80,
              },
            ],
            livenessProbe: {
              failureThreshold: 3,
              httpGet: {
                path: if 'health' in service then service.health else '/healthz',
                port: 80,
                scheme: 'HTTP',
                httpHeaders: [
                  {
                    name: 'Host',
                    value: std.strReplace(service.name, '-dot-', '.') + '.' + std.extVar('ctf_domain'),
                  },
                  {  // temporary hack
                    name: 'X-Cluster-Health-Check',
                    value: 'lol',
                  },
                ],
              },
              initialDelaySeconds: if 'startTime' in service then service.startTime else 30,
              periodSeconds: 10,
              successThreshold: 1,
              timeoutSeconds: 5,
            },
          },
        ],
      },
    },
  },
};

local kService(service) = {
  apiVersion: 'v1',
  kind: 'Service',
  metadata: {
    name: service.name,
    labels: {
      app: service.name,
    },
  },
  spec: {
    ports: [
      {
        port: 80,
      },
    ],
    selector: {
      app: service.name,
    },
  },
};

local kPV(service) = {
  kind: 'PersistentVolume',
  apiVersion: 'v1',
  metadata: {
    name: service.name + '-pv',
    labels: {
      app: service.name,
    },
  },
  spec: {
    accessModes: [
      'ReadWriteOnce',
    ],
    capacity: {
      storage: service.persistent,
    },
    hostPath: {
      path: '/data/' + service.name,
      type: 'DirectoryOrCreate',
    },
  },
};

local kPVC(service) = {
  kind: 'PersistentVolumeClaim',
  apiVersion: 'v1',
  metadata: {
    name: service.name + '-pv-claim',
  },
  spec: {
    accessModes: [
      'ReadWriteOnce',
    ],
    resources: {
      requests: {
        storage: service.persistent,
      },
    },
    selector: {
      matchLabels: {
        app: service.name,
      },
    },
  },
};


[kDeployment(service) for service in services] + [kService(service) for service in services] + [kPV(service) for service in services if 'persistent' in service] + [kPVC(service) for service in services if 'persistent' in service]
