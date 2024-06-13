const { DURIAN_DB_SECRET } = require("./env");

// 请确保密钥的长度足够用于简单加密（位异或），这里是 128K 字节
const key = Buffer.alloc(1024 * 128, DURIAN_DB_SECRET);

function encrypt(s) {
  const txt = Buffer.from(s);
  const len = txt.length;
  for (let i = 0; i < len; ++i) {
    txt[i] = txt[i] ^ key[i];
  }
  return txt.toString("base64").replace(/=+$/, "");
}

function decrypt(s) {
  const raw = Buffer.from(s, "base64");
  const len = raw.length;
  for (let i = 0; i < len; ++i) {
    raw[i] = raw[i] ^ key[i];
  }
  return raw.toString();
}

module.exports = {
  encrypt,
  decrypt,
};
