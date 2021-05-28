{
  services: [
    {
      name: 'ctfd',
      replicas: 1,
      category: 'infra',
      clustertype: 'master',
      persistent: '100M',
      access: |||
        def checkAccess():
          # NOTES: admin apis don't have specific separate endpoint so we only protect admin frontend here
          # we'll rely on CTFd's own permission settings for fine-grained api access control
          if path.startswith("/admin") and ("ctfd-admin@groups." + corpDomain) not in groups:
            denyAccess()
          grantAccess()
        checkAccess()
      |||,
    },
  ],
}
