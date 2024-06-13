const util = require("node:util");

module.exports = util.promisify((ms, cb) => {
  setTimeout(cb, ms);
});
