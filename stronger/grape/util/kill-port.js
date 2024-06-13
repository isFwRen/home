const child_process = require("child_process");
const log = require("electron-log");
module.exports = function (port, cb) {
  if (process.platform === "win32") {
    return child_process.exec(
      `netstat -nao | findstr :${port}`,
      (err, stdout, stderr) => {
        if (stdout == "") {
          // 没有找到这个端口运行的服务直接返回就行
          return cb("", "", "");
        }
        if (err || stderr) {
          return cb(err, stdout, stderr);
        }
        const lines = stdout.split("\n");
        var pids = lines.reduce((acc, line) => {
          const match = line.match(/(\d*)\w*(\n|$)/gm);
          return match && match[0] && !acc.includes(match[0])
            ? acc.concat(match[0])
            : acc;
        }, []);
        pids = pids.filter((item) => item !== "0");
        log.info("pids", pids);
        if (pids.length == 0) {
          return cb("", "", "");
        }
        return child_process.exec(
          `TaskKill /F /PID ${pids.join(" /PID ")}`,
          (err, stdout, stderr) => {
            return cb(err, stdout, stderr);
          }
        );
      }
    );
  }

  return child_process.exec(
    `lsof -i:${port} | grep :${port} | awk '{print $2}' | xargs kill -9`,
    (err, stdout, stderr) => {
      return cb(err, stdout, stderr);
    }
  );
};
