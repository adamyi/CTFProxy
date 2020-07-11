{
  services: [
    {
      name: 'xssbot',
      replicas: 3,
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
