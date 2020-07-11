local extractComponent(chal, component) = if component in chal then chal[component] else [];
{
  extractServices(chal):: extractComponent(chal, 'services'),
  extractEmails(chal):: extractComponent(chal, 'emails'),
  extractFlags(chal):: extractComponent(chal, 'flags'),
  extractFiles(chal):: extractComponent(chal, 'staticfiles'),
  extractCLIFiles(chal):: extractComponent(chal, 'clistaticfiles'),
  subnetToAddress(subnet, no):: subnet[0] + '.' + subnet[1] + '.' + subnet[2] + '.' + (subnet[3] + no),
}
