local challenges = import 'challenges/challenges.libsonnet';
local infras = import 'infra/challenges.libsonnet';
local utils = import 'infra/jsonnet/utils.libsonnet';

local combined = infras + challenges;

local services = std.flattenArrays([utils.extractServices(chal) for chal in combined]);

local transform(service, id) = service {
  subnet:: [100, 100, std.floor(id / 32), (id % 32) * 8],
};

[transform(services[i], i) for i in std.range(0, std.length(services) - 1)]
