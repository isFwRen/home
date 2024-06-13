const fs = require("node:fs/promises");
const logger = require("./logger");

// ndjson 文件解析
async function parse_ndjson(pathname) {
  const fields = [],
    docs = [];
  try {
    var file = await fs.open(pathname, "r");
    let i = 0;
    for await (let line of file.readLines()) {
      i += 1;
      try {
        const doc = JSON.parse(line);
        if (typeof doc !== "object" || Array.isArray(doc))
          throw new Error("行内容错误");
        docs.push(doc);
        Object.keys(doc).map((k) => {
          if (fields.indexOf(k) < 0) fields.push(k);
        });
      } catch (error) {
        logger.error(`ndjson行${i}解析失败 ${line.trim()}`);
      }
    }
  } finally {
    await file?.close();
  }
  return {
    fields,
    docs,
  };
}

module.exports = parse_ndjson;
