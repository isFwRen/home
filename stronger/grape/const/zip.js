const fs = require("node:fs/promises");
const {createWriteStream} = require("node:fs");
const path = require("node:path");
const util = require("node:util");
const archiver = require("archiver");
const StreamZip = require("node-stream-zip");
const logger = require("./logger");

function _zip_files(files, zipfile, callback) {
  const stream_output = createWriteStream(zipfile);
  stream_output.on("close", () => callback());
  const archive = archiver("zip", {
    zlib: { level: 9 },
  });
  archive.on("zip", (err) => logger.warn("zip", err));
  archive.on("error", (err) => {
    logger.error("zip failed:", err);
    stream_output.close(() => {
      fs.rm(zipfile);
    });
    callback(err);
  });
  archive.pipe(stream_output);
  files.map((f) => archive.file(f, { name: path.basename(f) }));
  archive.finalize();
}

// 解压缩zip文件
async function unzip_file(file, selected, save_pathname) {
  // logger.debug('unzip', file)
  const zip_stream = new StreamZip.async({ file });
  await fs.mkdir(path.dirname(save_pathname), { recursive: true });
  await zip_stream.extract(selected, save_pathname);
  await zip_stream.close();
  // logger.debug('unzip ok', save_pathname);
  return save_pathname;
}

exports.compress = util.promisify(_zip_files);
exports.decompress = unzip_file;
