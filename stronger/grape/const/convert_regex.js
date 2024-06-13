const { isRegExp, uniq } = require("lodash");

function convert_regex(obj) {
  if ([undefined, null].indexOf(obj) >= 0) return obj;
  if (isRegExp(obj)) return obj;
  if (Array.isArray(obj)) {
    obj.map((o) => convert_regex(o));
    return obj;
  }
  if (typeof obj !== "object") return obj;
  Object.keys(obj).map((k) => {
    if (k !== "$regex") return convert_regex(obj[k]);
    const str = obj[k];
    if (typeof str !== "string") return str;
    // 格式 /regex/options
    const parts = /^\/(.+)\/([dgimsuvy]*)$/.exec(str);
    if (!parts || parts.length < 3) return str;
    let options = [];
    if (obj.$options) options.push(obj.$options.split(""));
    if (parts[2]) options.push(parts[2].split(""));
    const $options = uniq(options).join("");
    obj[k] = new RegExp(parts[1], $options);
  });
  return obj;
}

module.exports = convert_regex;
