{
  services: [
    {
      name: 'xssbot',
      replicas: 1,
      category: 'infra',
      clustertype: 'master',
      access: |||
        def checkAccess():
          if user.endswith("@services." + corpDomain):
            grantAccess()
        checkAccess()
      |||,
    },
  ],
}
