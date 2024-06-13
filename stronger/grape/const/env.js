const os = require("node:os");
const path = require("node:path");
const common = require("./common");

// 服务地址与端口
exports.DURIAN_HOST = process.env.DURIAN_HOST || "127.0.0.1";
exports.DURIAN_PORT = parseInt(process.env.DURIAN_PORT || "30080");
// 服务端的API密钥，base64字符串，内部业务使用，不设置就无法访问
exports.DURIAN_API_TOKEN = process.env.DURIAN_API_TOKEN;

// 存储路径名，服务器版本是 .duriand，客户端版本是 .durian
const DIRNAME = process.env.DURIAN_MODE === "BACKEND" ? ".duriand" : ".durian";
// 用于处理上传文件的临时路径，默认在临时目录
exports.DURIAN_TEMP =
  process.env.DURIAN_TEMP ||
  path.join(os.tmpdir(), DIRNAME).replace(/\\/g, "/");

// 数据存储路径，默认在用户home目录
const DURIAN_DIR = (exports.DURIAN_DIR =
  process.env.DURIAN_DIR ||
  path.join(os.homedir(), DIRNAME).replace(/\\/g, "/"));
process.env.IPFS_PATH = `${DURIAN_DIR}/.ipfs`;

// 数据库文件的 AES256 密钥，32字符长，默认不加密
exports.DURIAN_DB_SECRET = process.env.DURIAN_DB_SECRET;

// 同一个数据表的最大备份数
exports.DURIAN_DB_BACKUP_MAX = parseInt(
  process.env.DURIAN_DB_BACKUP_MAX || "3"
);

// 最多在内存中同时加载2个数据库
exports.DURIAN_DB_LOAD_MAX = parseInt(
  process.env.DURIAN_DB_LOAD_MAX ||
    (process.env.DURIAN_MODE === "BACKEND" ? "6" : "3")
);
exports.DURIAN_DB_TABLE_LOAD_MAX = parseInt(
  process.env.DURIAN_DB_TABLE_LOAD_MAX ||
    (process.env.DURIAN_MODE === "BACKEND" ? "32" : "4")
);

// 导入文件时批处理数据个数
exports.DURIAN_BULK = parseInt(process.env.DURIAN_BULK || "5120");

// API前缀
exports.DURIAN_API_PREFIX = process.env.DURIAN_API_PREFIX || "/sys-const";
// 完整的API路径，正式版本应该用 https
exports.DURIAN_API =
  process.env.DURIAN_API ||
  `http://${exports.DURIAN_HOST}:${exports.DURIAN_PORT}${exports.DURIAN_API_PREFIX}`;

// 分页设置
exports.DURIAN_PAGE_SIZE = parseInt(process.env.DURIAN_PAGE_SIZE || "10");
exports.DURIAN_PAGE_SIZE_MAX = parseInt(
  process.env.DURIAN_PAGE_SIZE_MAX || "100"
);

// 客户端更新设置
// 完成更新后的动作，默认无操作，可以设置成 quit 或 restart
exports.DURIAN_ON_UPDATED = process.env.DURIAN_ON_UPDATED;
// 下载失败后的动作，默认 retry，也可以设置成 quit 或 restart
exports.DURIAN_ON_DOWNLOAD_FAILED =
  process.env.DURIAN_ON_DOWNLOAD_FAILED || "retry";
// 通过BT下载的百分比，默认是0, 仅使用HTTP下载
exports.DURIAN_PEER_DOWNLOAD_PROBABILITY = parseInt(
  process.env.DURIAN_PEER_DOWNLOAD_PROBABILITY || "0"
);
// 通过BT下载文件时的超时值，单位为秒，默认300秒，超时后降级为HTTP下载
exports.DURIAN_PEER_DOWNLOAD_TIMEOUT = parseInt(
  process.env.DURIAN_PEER_DOWNLOAD_TIMEOUT || "300"
);
// 通过HTTP下载文件时的超时值，单位为秒，默认30秒
exports.DURIAN_HTTP_DOWNLOAD_TIMEOUT = parseInt(
  process.env.DURIAN_HTTP_DOWNLOAD_TIMEOUT || "30"
);
// HTTP下载的最大并发数量，默认为4
exports.DURIAN_HTTP_DOWNLOAD_MAX = parseInt(
  process.env.DURIAN_HTTP_DOWNLOAD_MAX || "6"
);

// BT设置
// 服务器端的 tracker 端口，用于管理分布式数据，包括TCP和UDP，不能和api接口共用
exports.DURIAN_TRACKER_PORT = parseInt(
  process.env.DURIAN_TRACKER_PORT || exports.DURIAN_PORT + 1
);
// 服务器端的 DHT 端口
exports.DURIAN_DHT_PORT = parseInt(
  process.env.DURIAN_DHT_PORT || exports.DURIAN_PORT + 2
);
// DHT的启动节点
exports.DURIAN_DHT_BOOTSTRAP = common.stringlist(
  process.env.DURIAN_DHT_BOOTSTRAP ||
    `${exports.DURIAN_HOST}:${exports.DURIAN_DHT_PORT},router.bittorrent.com:6881,router.utorrent.com:6881,dht.transmissionbt.com:6881`
);
// 数据传输端口
exports.DURIAN_TORRENT_PORT = parseInt(
  process.env.DURIAN_TORRENT_PORT || exports.DURIAN_PORT + 3
);

// 下载及上传限速
exports.DURIAN_DOWNLOAD_LIMIT = parseInt(
  process.env.DURIAN_DOWNLOAD_LIMIT || "-1"
);
exports.DURIAN_UPLOAD_LIMIT = parseInt(process.env.DURIAN_UPLOAD_LIMIT || "-1");

const our_tracker = `${exports.DURIAN_HOST.replace(
  /^0\.0\.0\.0$/i,
  "127.0.0.1"
)}:${exports.DURIAN_TRACKER_PORT}`;

// 用于发布种子的 tracker 列表, 用逗号或换行隔开。可参考:
// https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best.txt
// https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_https.txt
exports.DURIAN_TRACKERS =
  process.env.DURIAN_TRACKERS ||
  common.stringlist(`udp://${our_tracker}/announce`);

// 屏幕日志时间格式 YYYY-MM-DD HH:mm:ssZ
exports.DURIAN_LOG_TIME = process.env.DURIAN_LOG_TIME || "MM-DD HH:mm:ss";

// 以下参数仅内部使用

// 后端模式还是客户端模式
exports.DURIAN_MODE = process.env.DURIAN_MODE;

// 产品模式
if (!process.env.NODE_ENV) {
  const prog = path.basename(process.argv[0], ".exe").toLowerCase();
  process.env.NODE_ENV = prog !== "node" ? "production" : "development";
}
exports.NODE_ENV = process.env.NODE_ENV;
