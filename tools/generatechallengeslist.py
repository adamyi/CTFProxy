#!/usr/bin/env python

import sys
import string

print("# THIS FILE IS AUTO-GENERATED")

chals = []

for f in sys.argv[1:]:
  p = f.split(":")[0][2:]
  cn = p.split("/")[-1].replace('-', '_')
  chals.append(cn)
  print("local %s = import '%s/challenge.libsonnet';" % (cn, p))

print("[%s]" % string.join(chals, ", "))
