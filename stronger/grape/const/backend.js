process.env.DURIAN_MODE = "BACKEND";
const { DURIAN_API_TOKEN, NODE_ENV } = require("./env");
const logger = require("./logger");

if (!DURIAN_API_TOKEN) {
  logger.error("未设置 base64 格式的环境变量: DURIAN_API_TOKEN, 存在严重的安全问题");
  if (NODE_ENV === 'production')
    return process.exit(-1);
}
if (DURIAN_API_TOKEN) {
  try {
    Buffer.from(DURIAN_API_TOKEN, "base64");
  } catch (err) {
    logger.error(
      "环境变量不符合 base64 格式:",
      `DURIAN_API_TOKEN=${DURIAN_API_TOKEN}`
    );
    process.exit(-2);
  }
}  

const app = require("./app");
require("./edit")(app);
