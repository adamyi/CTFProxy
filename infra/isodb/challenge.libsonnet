{
  services: [
    {
      name: 'isodb',
      replicas: 1,  // we need distributed lock to enable replicas, ceebs
      category: 'infra',
      clustertype: 'master',
      access: |||
        def checkAccess():
          # only allow other services to access me, i.e. no direct student access
          if user.endswith("@services." + corpDomain):
            grantAccess()
        checkAccess()
      |||,
    },
  ],
}
