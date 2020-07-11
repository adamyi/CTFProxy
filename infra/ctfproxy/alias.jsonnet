local services = import 'infra/jsonnet/services.libsonnet';
local utils = import 'infra/jsonnet/utils.libsonnet';

local aliases = std.flattenArrays([
  if 'alias' in service then [
    {
      alias: alias,
      origin: service.name,
    }
    for alias in service.alias
  ] else []
  for service in services
]);

{
  [alias.alias]: alias.origin
  for alias in aliases
}
