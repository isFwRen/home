const ULID = require("ulid");
const { ExpressPeerServer } = require("peer");
const logger = require("./logger");
const {
  DURIAN_API_PREFIX,
  DURIAN_CLUSTER_NAME,
  DURIAN_CLUSTER_TIMEOUT,
  DURIAN_CLUSTER_CONCURRENT_LIMIT,
} = require("./env");

module.exports = (app) => {
  const cluster_api = ExpressPeerServer(app._server, {
    key: DURIAN_CLUSTER_NAME,
    alive_timeout: 1000 * DURIAN_CLUSTER_TIMEOUT,
    concurrent_limit: DURIAN_CLUSTER_CONCURRENT_LIMIT,
    allow_discovery: true,
    proxied: true,
    debug: true,
    generateClientId: () => ULID.ulid(),
  });
  cluster_api.set("x-powered-by", false);
  cluster_api.on("connection", (client) => {
    logger.debug("peer join", client);
  });
  cluster_api.on("disconnect", (client) => {
    logger.debug("peer exit", client);
  });

  let api_intro = `## 取得集群新节点的ID
  - GET ${DURIAN_API_PREFIX}/cluster/durian/id
  ## 取得集群的全部在线节点
  - GET ${DURIAN_API_PREFIX}/cluster/durian/peers
  `;
  [
    `${DURIAN_API_PREFIX}/cluster`,
    `${DURIAN_API_PREFIX}/cluster/durian`,
    `${DURIAN_API_PREFIX}/cluster/durian/`,
  ].map((u) => app.get(u, api_description));
  app.use(`${DURIAN_API_PREFIX}/cluster`, cluster_api);
  logger.document(api_intro);
  return cluster_api;
};

function api_description(req, res) {
  res.json({
    name: "cluster API",
    description: "cluster API to broker connections between durian clients.",
    website: "https://www.i-confluence.com",
  });
}
