local SUBNET_CTFPROXY_OFFSET = 2;
local SUBNET_SERVICE_OFFSET = 3;
local SUBNET_DNS_OFFSET = 4;

local utils = import 'infra/jsonnet/utils.libsonnet';
local tmpservices = import 'infra/jsonnet/services.libsonnet';

local image(service) = if 'image' in service then
  service.image
else if service.category == 'infra' then
  '%s/infra/%s:latest' % [std.extVar('container_registry'), service.name]
else
  '%s/challenges/%s/%s:latest' % [std.extVar('container_registry'), service.category, service.name];


local services = [service for service in tmpservices if ((std.extVar('cluster') == 'all') || ('clustertype' in service && service.clustertype == std.extVar('cluster')))];

// NOTES(adamyi@): they mess up with uberproxy... disable it for now
local searchdomains = [];

local defaultdomain = '.' + std.extVar('ctf_domain');

local tservices = {
  [service.name]: {
    image: image(service),
    networks: {
      ['beyondcorp_' + service.name]: {
        aliases: [
          service.name + if 'domain' in service then service.domain else defaultdomain,
        ],
        ipv4_address: utils.subnetToAddress(service.subnet, SUBNET_SERVICE_OFFSET),
      },
    },
    dns: utils.subnetToAddress(service.subnet, SUBNET_DNS_OFFSET),
    dns_search: searchdomains,
    ports: if 'ports' in service then service.ports else [],
  } + if 'others' in service then service.others else {}
  for service in services
};

local networks = {
  ['beyondcorp_' + service.name]: {
    ipam: {
      driver: 'default',
      config: [
        {
          subnet: utils.subnetToAddress(service.subnet, 0) + '/29',
        },
      ],
    },
  }
  for service in services
};

{
  version: '2',
  services: {
    dns: {
      image: std.extVar('container_registry') + '/infra/dns:latest',
      networks: {
        ['beyondcorp_' + service.name]: {
          ipv4_address: utils.subnetToAddress(service.subnet, SUBNET_DNS_OFFSET),
        }
        for service in services
      },
    },
    ctfproxy: {
      image: std.extVar('container_registry') + '/infra/ctfproxy:latest',
      networks: {
        ['beyondcorp_' + service.name]: {
          ipv4_address: utils.subnetToAddress(service.subnet, SUBNET_CTFPROXY_OFFSET),
        }
        for service in services
      },
      ports: [
        '80:80',
        '443:443',
      ],
      dns_search: searchdomains,
      environment: [
        'CTFPROXY_CLUSTER=' + std.extVar('cluster'),
      ],
    },
  } + tservices,
  networks: networks,
}
