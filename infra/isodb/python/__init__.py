import requests
from ctf_domain import CTF_DOMAIN

ISODB_ENDPOINT = "https://isodb." + CTF_DOMAIN


class IsoDbError(Exception):
  pass


def initDB(version, sql):
  # wait for isodb to start
  ok = 500
  while ok >= 400:
    try:
      ok = requests.get(ISODB_ENDPOINT + "/healthz").status_code
    except:
      pass
  # isodb is up
  resp = requests.post(
      ISODB_ENDPOINT + "/api/init?version=" + version, data=sql)
  if resp.status_code >= 400:
    raise IsoDbError(resp.text)
  print("isodb %s initiated" % version)


def query(params, instance):
  resp = requests.post(
      ISODB_ENDPOINT + "/api/sql",
      json=params,
      headers={"X-CTFProxy-SubAcc": instance})
  if resp.status_code >= 400:
    raise IsoDbError(resp.text)
  return resp.json()


def queryWithJWT(params, jwt):
  resp = requests.post(
      ISODB_ENDPOINT + "/api/sql",
      json=params,
      headers={"X-CTFProxy-SubAcc-JWT": jwt})
  if resp.status_code >= 400:
    raise IsoDbError(resp.text)
  return resp.json()


def queryMaster(params):
  return query(params, "master")
