const fs = require("fs");
const express = require("express");
// const puppeteer = require('puppeteer');
const { Cluster } = require("puppeteer-cluster");
const jose = require("jose");

const PORT = process.env.PORT || 8080;
const TASKTIMEOUT = process.env.TASKTIMEOUT || 5000;
const NETIDLETIMEOUT = process.env.NETIDLETIMEOUT || 2000;
const NEWREQIDLETIMEOUT = process.env.NEWREQIDLETIMEOUT || 2000;
const MAXCONCURRENTY = process.env.MAXCONCURRENCY || 2;
const CTFDOMAIN = "." + process.env.CTFDOMAIN;
const app = express();

var publicKEY = jose.JWK.asKey(fs.readFileSync("/jwt.pub", "utf8"));

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// idle when there's no traffic in timeout, or no new request in reqtimeout
function waitForNetworkIdle(
  page,
  timeout,
  reqtimeout,
  maxInflightRequests = 0
) {
  page.on("request", onRequestStarted);
  page.on("requestfinished", onRequestFinished);
  page.on("requestfailed", onRequestFinished);

  let inflight = 0;
  let fulfill;
  let promise = new Promise((x) => (fulfill = x));
  let timeoutId = setTimeout(onTimeoutDone, timeout);
  let reqtimeoutId = setTimeout(onTimeoutDone, reqtimeout);
  return promise;

  function onTimeoutDone() {
    console.log("network idled or no new requests for a while");
    page.removeListener("request", onRequestStarted);
    page.removeListener("requestfinished", onRequestFinished);
    page.removeListener("requestfailed", onRequestFinished);
    fulfill();
  }

  function onRequestStarted() {
    clearTimeout(reqtimeoutId);
    reqtimeoutId = setTimeout(onTimeoutDone, reqtimeout);
    ++inflight;
    if (inflight > maxInflightRequests) clearTimeout(timeoutId);
  }

  function onRequestFinished() {
    if (inflight === 0) return;
    --inflight;
    if (inflight === maxInflightRequests)
      timeoutId = setTimeout(onTimeoutDone, timeout);
  }
}

(async () => {
  const cluster = await Cluster.launch({
    concurrency: Cluster.CONCURRENCY_CONTEXT,
    timeout: TASKTIMEOUT,
    maxConcurrency: 5,
    puppeteerOptions: {
      // headless : false,
      // slowMo : 250,
      executablePath: "/chrome-linux/chrome",
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
    },
  });
  await cluster.task(async ({ page, data }) => {
    console.log(data);
    page.on("console", (msg) => console.log("PAGE LOG:", msg.text()));
    await page.setRequestInterception(true);
    page.on("request", (request) => {
      console.log("Requesting " + request.url());
      const headers = request.headers();
      const hostname = new URL(request.url()).hostname;
      if (hostname.endsWith(CTFDOMAIN)) {
        headers["X-CTFProxy-SubAcc"] = data.subacc;
      }
      request.continue({
        headers,
      });
    });
    await page.setExtraHTTPHeaders({ "X-Powered-By": "CTFProxy/xssbot" });
    await Promise.all([
      page.goto(data.url),
      waitForNetworkIdle(page, NETIDLETIMEOUT, NEWREQIDLETIMEOUT, 0),
    ]);
    // await sleep(1500);
    console.log("done");
  });

  // setup server
  app.get("/healthz", (req, res) => res.send("ok"));
  app.get("/", async function (req, res) {
    console.log("incoming request");
    let token = req.headers["x-ctfproxy-jwt"];
    console.log(token);
    var djwt;
    if (token) {
      try {
        djwt = jose.JWT.verify(token, publicKEY);
      } catch (err) {
        console.log("token invalid");
        return res.json({ success: false, message: "Token is not valid" });
      }
    } else {
      console.log("auth token not supplied");
      return res.json({
        success: false,
        message: "Auth token is not supplied",
      });
    }

    if (!req.query.url) {
      console.log("no url");
      return res.json({ success: false, message: "url invalid" });
    }
    console.log(req.query.url);
    try {
      userparts = djwt["username"].split("@")[0].split("+");
      cluster.queue({
        url: req.query.url,
        subacc: userparts[userparts.length - 1],
      });
    } catch (err) {
      console.log(err.message);
      return res.json({ success: false, message: err.message });
    }
    console.log("queued");

    return res.json({ success: true, message: "queued" });
  });

  app.listen(PORT, function () {
    console.log("xssbot listening on port " + PORT);
  });
})();
