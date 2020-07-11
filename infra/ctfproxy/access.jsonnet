local services = import 'infra/jsonnet/services.libsonnet';
local utils = import 'infra/jsonnet/utils.libsonnet';

local defaultaccess = 'denyAccess()';

{
  'kubernetes-dashboard': |||
    def checkAccess():
      if ("k8s-admin@groups." + corpDomain) in groups:
        grantAccess()
    checkAccess()
  |||,
}
{
  [service.name]: if 'access' in service then service.access else defaultaccess
  for service in services
}
