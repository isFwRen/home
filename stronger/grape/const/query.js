const logger = require("./logger");
const Db = require("./db");
const convert_regex = require("./convert_regex");
const { check_proCode, check_name } = require("./check");
const {
  DURIAN_API_PREFIX,
  DURIAN_PAGE_SIZE,
  DURIAN_PAGE_SIZE_MAX,
} = require("./env");

// 获取项目清单
async function listDatabases(req, res) {
  // logger.debug("list databases");
  res.status(200).json({
    status: 200,
    msg: "操作成功",
    list: await Db.listDatabases(),
  });
}

// 查询某项目的常量清单
async function list_collections(req, res) {
  const proCode = req.params.proCode || req.body.proCode;
  // logger.debug("list_collections", proCode);
  const proj_db = await Db(proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      proCode,
    });
  res.status(200).json({
    status: 200,
    msg: "操作成功",
    proCode,
    list: await proj_db.listTables(),
  });
}

// 常量查询
async function db_query(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
  };
  if (!req.body || !base_data.proCode || !base_data.name)
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
      ...base_data,
    });
  const proj_db = await Db(base_data.proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
    });
  if (!(await proj_db.tableExists(base_data.name)))
    return res.status(404).json({
      status: 404,
      msg: "常量表不存在",
      ...base_data,
    });
  if (req.body.pageIndex === undefined) req.body.pageIndex = 1;
  if (req.body.pageSize === undefined) req.body.pageSize = DURIAN_PAGE_SIZE;
  if (req.body.pageSize > DURIAN_PAGE_SIZE_MAX)
    req.body.pageSize = DURIAN_PAGE_SIZE_MAX;
  try {
    convert_regex(req.body.queryNames);
    const coll = proj_db.table(base_data.name);
    const total = await coll.countEx(req.body.queryNames);
    const result = {
      status: 200,
      msg: "操作成功",
      ...base_data,
      list: [],
      total,
      pageIndex: req.body.pageIndex,
      pageSize: req.body.pageSize,
    };
    if (req.body.sort) result.sort = req.body.sort;
    if (total === 0) return res.status(200).json(result);
    const opt = {
      skip: (req.body.pageIndex - 1) * req.body.pageSize,
      limit: req.body.pageSize,
    };
    if (req.body.respNames) opt.fields = req.body.respNames;
    if (req.body.sort) opt.sort = req.body.sort;
    result.list = await coll.findEx(req.body.queryNames, opt).toArrayEx();
    return res.status(200).json(result);
  } catch (err) {
    return res.status(500).json({
      status: 500,
      msg: err.message || err,
      ...base_data,
    });
  }
}

module.exports = (app) => {
  let api_intro = `# API 说明
## 获取所有项目:
- GET  ${DURIAN_API_PREFIX}/info-list
## 获取某项目的常量表清单:
- GET  ${DURIAN_API_PREFIX}/info-list/:proCode
- POST ${DURIAN_API_PREFIX}/info-list { proCode }
`;
  app.get(`${DURIAN_API_PREFIX}/info-list`, Db.waitReady, listDatabases);
  app.get(`${DURIAN_API_PREFIX}/info-list/`, Db.waitReady, listDatabases);
  app.get(
    `${DURIAN_API_PREFIX}/info-list/:proCode`,
    check_proCode,
    Db.waitReady,
    list_collections
  );
  app.post(
    `${DURIAN_API_PREFIX}/info-list`,
    check_proCode,
    Db.waitReady,
    list_collections
  );

  api_intro += `## 分页查询
- POST ${DURIAN_API_PREFIX}/page/:proCode/:name { queryNames }
- POST ${DURIAN_API_PREFIX}/page { proCode, name, queryNames }`;
  app.post(
    `${DURIAN_API_PREFIX}/page/:proCode/:name`,
    check_proCode,
    check_name,
    Db.waitReady,
    db_query
  );
  app.post(
    `${DURIAN_API_PREFIX}/page`,
    check_proCode,
    check_name,
    Db.waitReady,
    db_query
  );
  logger.document(api_intro);
};
