const path = require("node:path");
const serve_handler = require('serve-handler');
const { DURIAN_API_PREFIX } = require("./env");

const PUBLIC = "public";
const STATIC_DIR = path
  // .relative(process.cwd(), path.join(path.dirname(process.argv[1]), PUBLIC))
  .relative(process.cwd(), path.join(path.dirname(__dirname), PUBLIC))
  .replace(/\\/g, "/");

  module.exports = async (app) => {
    app.use((req, res) => {
      serve_handler(req, res, {
        public: STATIC_DIR,
        rewrites: [
          { source: '/', destination: '/index.html' },
          { source: '/favicon.ico', destination: '/icon.webp'},
          // { source: `${DURIAN_API_PREFIX}/sio/(.+)`, destination: '/socket.io/$1'},
        ],
        unlisted: [
          "/**/.*",
        ],
        directoryListing: true,
        etag: true,
        cleanUrls: false,
      })
    })
  }
