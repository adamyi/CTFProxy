{
  services: [
    {
      name: 'xssbot',
      replicas: 4,
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
