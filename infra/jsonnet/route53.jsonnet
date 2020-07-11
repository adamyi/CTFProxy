local TTL = 86400;

local services = import 'infra/jsonnet/services.libsonnet';

local sns = [service.name for service in services] + std.flattenArrays([if 'alias' in service then service.alias else [] for service in services]) +
            ['cli-relay', 'login', 'ctfproxyz', 'kubernetes-dashboard'];

local s2d(sn) = std.strReplace(sn, '-dot-', '.') + '.' + std.extVar('ctf_domain') + '.';

local domains = [s2d(sn) for sn in sns];

local upsertRecord(recordSet) = {
  Action: 'UPSERT',
  ResourceRecordSet: recordSet,
};

// NOTES: you can also add none-challenge records manually below
// to be automated and version-controlled
{
  Changes: [
    upsertRecord({
      Name: domain,
      Type: 'CNAME',
      TTL: TTL,
      ResourceRecords: [
        {
          Value: s2d('master-dot-prod'),
        },
      ],
    })
    for domain in domains
  ],
}
