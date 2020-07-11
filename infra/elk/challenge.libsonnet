{
  services: [
    {
      name: 'elasticsearch',
      category: 'infra',
      image: 'gcr.io/ctfproxy/elk/elasticsearch:latest',
      clustertype: 'master',
      access: |||
        def checkAccess():
          if user == "kibana@services." + corpDomain:
            grantAccess()
        checkAccess()
      |||,
      health: '/_cat/health',
      startTime: 60,
      persistent: '250G',
    },
    {
      name: 'kibana',
      category: 'infra',
      image: 'gcr.io/ctfproxy/elk/kibana:latest',
      clustertype: 'master',
      access: |||
        def checkAccess():
          if ("kibana-user@groups." + corpDomain) in groups:
            grantAccess()
        checkAccess()
      |||,
      health: '/api/status',
      startTime: 60,
    },
  ],
}
