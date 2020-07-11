local challenges = import 'challenges/challenges.libsonnet';
local infras = import 'infra/challenges.libsonnet';
local utils = import 'infra/jsonnet/utils.libsonnet';

local combined = challenges + infras;

/*Id string
DisplayName string
Category string
Points int
Type string
Flag string
Prefix string
Owner string*/

local parseFlag(flag, chal) = {
  Id: if 'Id' in flag then flag.Id else error 'Id must be set',
  DisplayName: if 'DisplayName' in flag then flag.DisplayName else chal.services[0].name,
  Category: if 'Category' in flag then flag.Category else chal.services[0].category,
  Points: if 'Points' in flag then flag.Points else error 'Points must be set',
  Type: if 'Type' in flag && (flag.Type == 'fixed' || flag.Type == 'dynamic') then flag.Type else error 'Invalid type',
  Flag: if 'Flag' in flag then flag.Flag else error 'Flag must be set',
  Prefix: if 'Prefix' in flag then flag.Prefix else 'FLAG',
  Owner: if 'Owner' in flag then flag.Owner else chal.services[0].name,
};

local extractFlags(chal) = if 'flags' in chal then [parseFlag(flag, chal) for flag in chal.flags] else [];

std.flattenArrays([extractFlags(chal) for chal in combined])
