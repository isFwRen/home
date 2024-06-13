const fs = require("node:fs/promises");
const path = require("node:path");
const util = require("node:util");
const exec = util.promisify(require("node:child_process").exec);
const got = require("got");
const logger = require("./logger");
const sleep = require("./sleep");
const { DURIAN_MODE } = require("./env");

const PLATFORM = require("node:os").platform();

(async function () {
  // 清理上次残留的进程
  for (const name of ["aria2c", "ipfs"]) {
    try {
      await exec(
        PLATFORM === "win32"
          ? `taskkill /f /im "${name}.exe"`
          : `killall -9 "${name}"`
      );
    } catch (error) {}
  }
  // 添加可执行文件路径
  process.env.PATH = [
    process.env.PATH,
    // path.dirname(process.argv0),
    "assets/lib",
    "resources/assets/lib",
  ].join(PLATFORM === "win32" ? ";" : ":");
  // 检测所需的工具是否存在
  const tools =
    DURIAN_MODE === "BACKEND"
      ? {
          // ipfs: {
          //   test: "ipfs --version",
          //   needed: true,
          // },
        }
      : {
          aria2: {
            test: "aria2c --version",
            needed: true,
          },
          // ipfs: {
          //   test: "ipfs --version",
          // },
        };
  for (const name of Object.keys(tools)) {
    const tool = tools[name];
    try {
      const { stdout, stderr } = await exec(tool.test);
      logger.debug(`${stdout}`.trim());
    } catch (error) {
      if (tool.needed) {
        logger.error(`未找到下载工具 ${name}，建议重启进程`);
      } else {
        logger.warn(`未找到下载工具 ${name}`);
      }
    }
  }
})();

async function download({
  url,
  pathname,
  attempts = 30,
  timeout = 5000,
  method = "GET",
  headers,
}) {
  if (pathname === undefined || pathname === null || pathname === "")
    pathname = url.split(/[\/\\]+/).pop();
  // logger.debug(`download ${url} => ${pathname}`);
  const dirname = path.dirname(pathname).replace(/\\+/, "/");
  const filename = path.basename(pathname);
  if (dirname !== ".") await fs.mkdir(dirname, { recursive: true });

  const { stdout, stderr } = await exec(
    `aria2c -o "${filename}" -d "${dirname}" "${url}" -t ${Math.ceil(
      timeout / 1000
    )} --file-allocation=none -R -c -s 4 -j 4 -x 4 -k 1M -m 0 --retry-wait=5 --check-certificate=false --max-overall-upload-limit=30K -q --summary-interval=0 --log-level=error --console-log-level=error --follow-torrent=false -U "durian/1.0"`
  );
  if (stdout || stderr)
    logger.warn("aria2c", url, `-o ${filename} =>`, { stdout, stderr });
}

// 旧的从头顺序下载，不如多块同时下载的速度
// const gotResume = require("got-resume");
// await gotResume.toFile(pathname, url, {
//   attempts,
//   timeout,
//   got: { method },
//   pre: async (transfer) => {
//     transfer.gotOptions.headers["user-agent"] = "durian/1.0";
//     if (!headers) return;
//     for (const k of Object.keys(headers)) {
//       transfer.gotOptions.headers[k] = headers[k];
//     }
//   },
// });

async function fetch(options) {
  return await got(options);
}

module.exports = {
  fetch,
  download,
};
