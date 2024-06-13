const { DURIAN_API_TOKEN } = require("./env");
const INVALID_CHARS = (exports.INVALID_CHARS = /[\/\\:\*\?"<>\|\x00-\x19]/g);
const INVALID_STR = (exports.INVALID_STR =
  /^(aux|con|prn|nul|com\d+|lpt\d+)(\.[^\.]+)?$/i);

exports.check_proCode = (req, res, next) => {
  if (!req.params) req.params = {};
  const proCode = req.params.proCode || req.body.proCode;
  if (
    [undefined, null, ""].indexOf(proCode) >= 0 ||
    INVALID_CHARS.test(proCode) ||
    INVALID_STR.test(proCode)
  )
    return res.status(404).json({
      status: 500,
      msg: "参数错误",
      ...req.params,
      ...(req.body || {}),
    });
  next();
};

exports.check_name = (req, res, next) => {
  if (!req.params) req.params = {};
  const name = req.params.name || req.body.name;
  if (
    [undefined, null, ""].indexOf(name) >= 0 ||
    INVALID_CHARS.test(name) ||
    INVALID_STR.test(name)
  )
    return res.status(404).json({
      status: 500,
      msg: "参数错误",
      ...req.params,
      ...(req.body || {}),
    });
  next();
};

TOKEN_FORMAT = /^Bearer\s+(\S+)$/i;
exports.check_api_token = (req, res, next) => {
  if (!DURIAN_API_TOKEN) return next();
  try {
    const auth_header = req.get('Authorization');
    const token = req.query.api_token || TOKEN_FORMAT.exec(auth_header)[1];
    Buffer.from(token, "base64");
    if (DURIAN_API_TOKEN === token) return next();
  } catch (err) {
    return res.status(401).end();
  }
  return res.status(401).end();
}
