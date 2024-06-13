const express = require("express");
const chalk = require("chalk");
const app = express();
const common = require("./common");
const logger = require("./logger");
const cors = require("cors");
const {
  DURIAN_MODE,
  DURIAN_HOST,
  DURIAN_PORT,
  DURIAN_API,
  DURIAN_API_PREFIX,
} = require("./env");

app.enable("trust proxy");
app.set("x-powered-by", false);

app.use((req, res, next) => {
  if (req.url.indexOf(DURIAN_API_PREFIX) !== 0) return next();
  logger.debug(
    req.method,
    req.url,
  );
  // res._startAt = Date.now();
  // if (!res._end) {
  //   res._end = res.end;
  //   res.end = (...args) => {
  //     res._endAt = Date.now();
  //     const spent = res._endAt - res._startAt;
  //     const bytes = parseInt(res.get("Content-Length") || "0");
  //     logger.debug(
  //       req.method,
  //       req.url,
  //       res.statusCode,
  //       bytes,
  //       "bytes",
  //       spent,
  //       "ms"
  //     );
  //     res._end(...args);
  //   };
  // }
  next();
});

app._server = require("node:http").createServer(app, {
  allowEIO3: true, // 兼容旧版 socket.io
  pingInterval: 50000,
  pingTimeout: 40000,
  transports: ["websocket", "polling"],
  // path: DURIAN_API_PREFIX, // 此参数有问题，会造成404无法连接
});
app._io = new (require("socket.io").Server)(app._server, { cors: true });

// for parsing application/json
app.use(express.json({ limit: "200mb" }));
// for parsing application/x-www-form-urlencoded
app.use(express.urlencoded({ extended: true, limit: "200mb" }));

// cors
app.use(cors());

// static files
setTimeout((e) => {
  require("./static")(app);
}, 100);

require("./query")(app);

process.on("unhandledRejection", (reason, p) => {
  logger.error(reason, "unhandled rejection at Promise", p);
});
process.on("uncaughtException", (err) => {
  logger.error("uncaughtException", err.stack);
});

const hostname = DURIAN_MODE === "BACKEND" ? DURIAN_HOST : "127.0.0.1";
logger.info(
  `durian ${
    DURIAN_MODE === "BACKEND" ? "backend" : "client"
  } listening on ${chalk.green(`http://${hostname}:${DURIAN_PORT}`)}`
);
app._online = false;
app._server.listen(DURIAN_PORT, hostname, () =>
  logger.debug("press ctrl-c to exit")
);

const socket = (app._socket = require("socket.io-client").io(
  new URL(DURIAN_API).origin,
  {
    ackTimeout: 90000,
  }
));
socket.only = (...args) => socket.off(...args).on(...args);

if (DURIAN_MODE === "BACKEND") {
  socket.only("notify", (m) =>
    logger.info("推送更新", `${m.proCode}.${m.name}`, m.action)
  );
} else {
  socket.only("disconnect", (err) => {
    logger.warn(`离线`);
    app._online = false;
  });
  socket.only("connect", () => {
    logger.info(`在线`);
    app._online = true;
  });
  socket.only("notify", (m) =>
    logger.info("收到更新", `${m.proCode}.${m.name}`, m.action)
  );
}

module.exports = app;
