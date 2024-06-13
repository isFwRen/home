const fs = require("node:fs/promises");
const path = require("node:path");
const qs = require("node:querystring");
const XLSX = require("xlsx");
const { chunk } = require("lodash");
const klaw = require("klaw");
const logger = require("./logger");
const Db = require("./db");
const seed = require("./seed");
const common = require("./common");
const parse_ndjson = require("./parse_ndjson");
const sleep = require("./sleep");

const {
  check_api_token,
  check_proCode,
  check_name,
  INVALID_CHARS,
  INVALID_STR,
} = require("./check");
const {
  DURIAN_TEMP,
  DURIAN_DIR,
  DURIAN_API_PREFIX,
  DURIAN_BULK,
  DURIAN_PAGE_SIZE_MAX,
} = require("./env");

// 清理临时文件
// @TODO: 定时删除旧文件，因为服务节点不一定经常重启
(async (e) => {
  try {
    await fs.rm(DURIAN_TEMP, { recursive: true, force: true });
  } catch (error) {}
  await fs.mkdir(DURIAN_TEMP, { recursive: true });
  for await (const file of klaw(DURIAN_TEMP)) {
    try {
      const relative_path = path.relative(DURIAN_TEMP, file.path);
      if (relative_path === "") continue;
      fs.rm(file.path, { recursive: true, force: true });
    } catch (error) {}
  }
})();

const upload = require("multer")({
  dest: DURIAN_TEMP,
  // 文件名补丁，解决中文乱码问题
  fileFilter: (req, file, cb) => {
    file.originalname = path
      .basename(Buffer.from(file.originalname, "latin1").toString("utf8"))
      .replace(INVALID_CHARS, "-")
      .replace(INVALID_STR, "-");
    cb(null, true);
  },
});

let notify;
module.exports = (app) => {
  notify = (m) => app._io.emit("notify", m);
  // 后端维护功能
  let api_intro = `## 导入常量表
- PUT  ${DURIAN_API_PREFIX}/import/:proCode { files:[...] }
- POST ${DURIAN_API_PREFIX}/import/:proCode { files:[...] }
`;
  app.put(
    `${DURIAN_API_PREFIX}/import/:proCode`,
    check_api_token,
    check_proCode,
    upload.array("files"),
    Db.waitReady,
    import_files
  );
  app.post(
    `${DURIAN_API_PREFIX}/import/:proCode`,
    check_api_token,
    check_proCode,
    upload.array("files"),
    Db.waitReady,
    import_files
  );

  api_intro += `## 添加多个常量
- PUT  ${DURIAN_API_PREFIX}/insert/:proCode/:name { items:[...] }
- POST ${DURIAN_API_PREFIX}/insert { proCode, name, items:[...] }
`;
  app.put(
    `${DURIAN_API_PREFIX}/insert/:proCode/:name`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    insert_docs
  );
  app.post(
    `${DURIAN_API_PREFIX}/insert`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    insert_docs
  );

  api_intro += `## 修改常量
- PATCH ${DURIAN_API_PREFIX}/edit/:proCode/:name/:id { item }
- POST  ${DURIAN_API_PREFIX}/edit { proCode, name, id, item }
`;
  app.patch(
    `${DURIAN_API_PREFIX}/edit/:proCode/:name/:id`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    edit_doc
  );
  app.post(
    `${DURIAN_API_PREFIX}/edit`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    edit_doc
  );

  api_intro += `## 删除多个常量
- DELETE ${DURIAN_API_PREFIX}/del-lines/:proCode/:name { ids:[...] }
- POST   ${DURIAN_API_PREFIX}/del-lines { proCode, name, ids:[...] }
`;
  app.delete(
    `${DURIAN_API_PREFIX}/del-lines/:proCode/:name`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    remove_docs
  );
  app.post(
    `${DURIAN_API_PREFIX}/del-lines`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    remove_docs
  );

  api_intro += `## 删除常量
- DELETE ${DURIAN_API_PREFIX}/del-docs/:proCode/:name { ids:[...] }
- POST   ${DURIAN_API_PREFIX}/del-docs { proCode, name, ids:[...] }
`;
  app.delete(
    `${DURIAN_API_PREFIX}/del-docs/:proCode/:name`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    del_docs
  );
  app.post(
    `${DURIAN_API_PREFIX}/del-docs`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    del_docs
  );

  api_intro += `## 删除常量表
- DELETE ${DURIAN_API_PREFIX}/del-tables/:proCode { name }
- POST   ${DURIAN_API_PREFIX}/del-tables { proCode, name }
`;
  app.delete(
    `${DURIAN_API_PREFIX}/del-tables/:proCode`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    remove_tables
  );
  app.post(
    `${DURIAN_API_PREFIX}/del-tables`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    remove_tables
  );

  api_intro += `## 导出一个常量表，下载格式为 xlsx
- GET  ${DURIAN_API_PREFIX}/export/:proCode/:name
- POST ${DURIAN_API_PREFIX}/export { proCode, name }
`;
  app.get(
    `${DURIAN_API_PREFIX}/export/:proCode/:name`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    export_file
  );
  app.post(
    `${DURIAN_API_PREFIX}/export`,
    check_api_token,
    check_proCode,
    check_name,
    Db.waitReady,
    export_file
  );

  api_intro += `## 发布新版数据表，使客户端自动更新
- PATCH ${DURIAN_API_PREFIX}/release/:proCode { name: [...] }
- POST  ${DURIAN_API_PREFIX}/release { proCode, name: [...] }
- socket.on("notify", ({proCode, name, action: "update|drop", info}) => {})
`;
  app.patch(
    `${DURIAN_API_PREFIX}/release/:proCode`,
    check_api_token,
    check_proCode,
    Db.waitReady,
    release
  );
  app.patch(
    `${DURIAN_API_PREFIX}/release/:proCode/:name`,
    check_api_token,
    check_proCode,
    Db.waitReady,
    release
  );
  app.post(
    `${DURIAN_API_PREFIX}/release`,
    check_api_token,
    check_proCode,
    Db.waitReady,
    release
  );

  api_intro += `## 下载 常量表名.zip 和 .torrent 文件，用于自动同步的降级处理
- GET ${DURIAN_API_PREFIX}/release/:proCode/:name`;
  app.get(
    `${DURIAN_API_PREFIX}/release/:proCode/:name`,
    check_proCode,
    check_name,
    http_download_table
  );
  app.get(`${DURIAN_API_PREFIX}/release/:hash`, http_download_table);
  logger.document(api_intro);
};

async function insert_docs(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
  };
  logger.debug(
    "新增数据",
    `${base_data.proCode}.${base_data.name}`,
    req.body.items.length
  );
  if (
    !base_data.proCode ||
    !base_data.name ||
    !req.body.items ||
    typeof req.body.items !== "object" ||
    req.body.items.length === 0
  )
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
    });
  const proj_db = await Db(base_data.proCode, { create: true });
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
    });
  const coll = proj_db.table(base_data.name);
  // logger.debug("insertEx", base_data, req.body.items);
  const items = await coll.insertEx(req.body.items);
  // logger.debug("items", items);
  await on_update(coll);
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    items,
  });
}

async function on_update(coll) {
  await coll._db.table("META").updateEx(
    {
      name: coll._name,
    },
    {
      $set: {
        updatedAt: Date.now(),
        released: 0,
      },
    },
    {
      upsert: true,
    }
  );
}

async function edit_doc(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
    id: req.params.id || req.body.id,
  };
  logger.debug("edit", base_data);
  if (!base_data.proCode || !base_data.name || !base_data.id || !req.body.item)
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
      ...base_data,
      item: req.body.item,
    });
  const proj_db = await Db(base_data.proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
    });
  const coll = proj_db.table(base_data.name);
  // logger.debug('updateEx', { _id: base_data.id }, req.body.item)
  const count = await coll.updateEx({ _id: base_data.id }, req.body.item);
  if (count === 0)
    return res.status(404).json({
      status: 404,
      msg: "数据不存在",
      ...base_data,
    });
  await on_update(coll);
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    ...base_data,
  });
}

async function remove_docs(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
    ids: req.body.ids,
  };
  logger.debug("remove", base_data);
  if (
    !base_data.proCode ||
    !base_data.name ||
    !req.body.ids ||
    req.body.ids.length === 0
  )
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
      ...base_data,
      total: 0,
    });
  const proj_db = await Db(base_data.proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
      total: 0,
    });
  const coll = proj_db.table(base_data.name);
  const total = await coll.removeEx(
    {
      _id: {
        $in: req.body.ids,
      },
    },
    {
      multi: true,
    }
  );
  if (total === 0)
    return res.status(404).json({
      status: 404,
      msg: "数据不存在",
      ...base_data,
      total,
    });
  await on_update(coll);
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    ...base_data,
    total,
  });
}

async function del_docs(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
    item: req.body.item,
  };
  logger.debug("del", base_data);
  if (!base_data.proCode || !base_data.name || !req.body.item)
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
      ...base_data,
      total: 0,
    });
  const proj_db = await Db(base_data.proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
      total: 0,
    });
  const coll = proj_db.table(base_data.name);
  const total = await coll.removeEx(
    {
      ...base_data.item,
    },
    {
      multi: true,
    }
  );
  if (total === 0)
    return res.status(404).json({
      status: 404,
      msg: "数据不存在",
      ...base_data,
      total,
    });
  await on_update(coll);
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    ...base_data,
    total,
  });
}

async function remove_tables(req, res) {
  const droped = [];
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
  };
  base_data.name = common.stringlist(base_data.name);
  logger.log("remove tables", base_data);
  if (!base_data.proCode || !base_data.name || base_data.name.length === 0)
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
      ...base_data,
      total: droped.length,
      droped,
    });
  const proj_db = await Db(base_data.proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
      total: droped.length,
      droped,
    });
  for (let coll_name of base_data.name) {
    coll_name = coll_name.toUpperCase();
    try {
      // 停止分享
      await seed.stopShare(proj_db, coll_name);
      const info = await proj_db.getTableMeta(coll_name);
      if (await proj_db.dropCollectionEx(coll_name)) {
        droped.push(coll_name);
        notify({
          action: "drop",
          proCode: proj_db._name,
          name: info.name,
          info,
        });
      }
    } catch (error) {
      logger.log("drop table", `${base_data.proCode}.${coll_name}`, error);
    }
  }
  if (droped.length === 0)
    return res.status(404).json({
      status: 404,
      msg: "数据表不存在",
      ...base_data,
      total: droped.length,
      droped,
    });
  // 如果数据库已经空了，要不要删库？
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    ...base_data,
    total: droped.length,
    droped,
  });
}

async function import_files(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
  };
  logger.debug("import to", base_data.proCode, "files:", req.files.length);
  const proj_db = await Db(base_data.proCode, { create: true });
  const tables = [];
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
      total: tables.length,
      tables,
    });
  for (const file of req.files || []) {
    await (async (file) => {
      logger.debug(
        "import",
        base_data.proCode,
        file.originalname,
        "bytes:",
        file.size
      );
      const is_xls = /\.(xlsx?|csv)$/i.test(file.originalname);
      const { fields, docs } = is_xls
        ? await parse_excel(file.path)
        : await parse_ndjson(file.path);
      // 对于无类型的 xls/csv, 把所有的值转换为字符串类型, ndjson则不做转换
      if (is_xls) {
        for (const doc of docs) {
          for (const k of fields) {
            const v = doc[k];
            if (v === undefined || v === null) continue;
            if (typeof v !== "string") doc[k] = `${v}`;
          }
        }
      }
      // 删除上传的文件
      fs.rm(file.path, { force: true }).catch((e) => null);
      const coll_name = file.originalname
        .replace(/\.+[^\.]*$/, "")
        .toUpperCase();
      logger.debug(
        "上传结果",
        coll_name,
        "数据项:",
        docs.length,
        "字段:",
        fields
      );
      // 加锁
      if (await proj_db.tableExists(coll_name))
        await proj_db.table(coll_name).archive();
      // 新表建索引
      const idx = fields.map((f) => [f, 1]);
      const coll = proj_db.table(coll_name);
      coll.importing = true;
      await coll.createIndex(idx);
      logger.debug("更新META", `${base_data.proCode}.${coll_name}`, fields);
      await proj_db.table("META").updateEx(
        {
          name: coll._name,
        },
        {
          $set: {
            name: coll._name,
            header: fields,
            total: docs.length,
            updatedAt: Date.now(),
            released: 0,
          },
        },
        {
          upsert: true,
        }
      );
      // 导入数据
      // logger.debug("create", coll_name);
      const bulks = chunk(docs, DURIAN_BULK);
      for (const bulk of bulks) {
        await coll.insertEx(bulk);
        await sleep(1);
      }
      coll.importing = false;
      // 解锁
      tables.push({
        name: coll_name,
        header: fields,
        total: docs.length,
      });
      logger.debug("file done", file.originalname);
    })(file);
  }
  // 整理 meta
  await proj_db.table("META").compactCollectionEx();
  logger.debug("import done", base_data.proCode);
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    ...base_data,
    total: tables.length,
    tables,
  });
}

async function export_file(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name,
  };
  if (!base_data.proCode || !base_data.name) return res.status(500).end();
  const proj_db = await Db(base_data.proCode);
  if (!proj_db) return res.status(404).end();
  const info = await proj_db.getTableMeta(base_data.name);
  if (!info || !info.header) return res.status(404).end();
  const data = [info.header];
  const coll = proj_db.table(base_data.name);
  let skip = 0;
  const limit = DURIAN_PAGE_SIZE_MAX;
  for (;;) {
    const docs = await coll
      .findEx(
        {},
        {
          skip,
          limit,
          sort: { _id: 1 },
        }
      )
      .toArrayEx();
    skip += limit;
    if (!docs || !docs.length) break;
    const lines = docs.map((doc) => info.header.map((field) => doc[field]));
    data.push(...lines);
    if (docs.length < limit) break;
  }
  logger.debug(base_data.name, data.length);
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(
    workbook,
    XLSX.utils.aoa_to_sheet(data),
    "Sheet1",
    true
  );
  const buf = XLSX.write(workbook, {
    bookType: "xlsx",
    compression: true,
    type: "buffer",
  });
  res.set(
    "Content-Disposition",
    `inline; filename="${qs.escape(base_data.name)}.xlsx"`
  );
  res.set("Content-Length", buf.length);
  res.set(
    "Content-Type",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
  );
  res.send(buf);
}

async function release(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.body.proCode,
    name: req.params.name || req.body.name || [],
  };
  base_data.name = common.stringlist(base_data.name);
  let items = [];
  if (!base_data.proCode || base_data.name.length === 0)
    return res.status(500).json({
      status: 500,
      msg: "参数错误",
      ...base_data,
      total: items.length,
      items,
    });
  const proj_db = await Db(base_data.proCode);
  if (!proj_db)
    return res.status(404).json({
      status: 404,
      msg: "项目不存在",
      ...base_data,
      total: items.length,
      items,
    });
  for (let coll_name of base_data.name) {
    coll_name = coll_name.toUpperCase();
    if (!(await proj_db.tableExists(coll_name))) continue;
    const coll = proj_db.table(coll_name);
    const torrent = await seed.release(coll);
    items.push({
      name: coll_name,
      hash: torrent.infoHash,
      magnet: torrent.magnetURI,
    });
    // 推送通知
    notify({
      action: "update",
      proCode: coll._db._name,
      name: coll._name,
      info: await coll._db.getTableMeta(coll._name),
    });
  }
  if (items.length === 0)
    return res.status(404).json({
      status: 404,
      msg: "数据表不存在",
      ...base_data,
      total: items.length,
      items,
    });
  return res.status(200).json({
    status: 200,
    msg: "操作成功",
    ...base_data,
    total: items.length,
    items,
  });
}

async function http_download_table(req, res) {
  const base_data = {
    proCode: req.params.proCode || req.query.proCode,
    name: req.params.name || req.query.name,
    hash: req.params.hash || req.query.hash,
  };
  if (!base_data.hash && (!base_data.proCode || !base_data.name))
    return res.status(500).end();
  let pathname = "";
  if (base_data.hash) {
    const hashInfo = seed.hash(base_data.hash);
    if (!hashInfo) return res.status(404).end();
    pathname = `${hashInfo.path}/${hashInfo.filename}.torrent`;
  }
  const proj_db = await Db(base_data.proCode);
  if (!proj_db) return res.status(404).end();
  if (/\.(torrent|zip)$/i.test(base_data.name))
    pathname = `${DURIAN_DIR}/seed/${proj_db._name}/${base_data.name}`;
  else
    pathname = `${DURIAN_DIR}/seed/${proj_db._name}/${base_data.name}.zip.torrent`;
  try {
    await fs.access(pathname);
  } catch (err) {
    return res.status(404).end();
  }
  res.status(200).sendFile(pathname, {
    headers: {
      "Content-Disposition": `inline; filename="${qs.escape(
        path.basename(pathname)
      )}"`,
      "x-timestamp": Date.now(),
    },
  });
}

// excel 文件解析
async function parse_excel(pathname) {
  // logger.debug("parse", pathname);
  const workbook = XLSX.readFile(pathname);
  // logger.debug("workbook.SheetNames", workbook.SheetNames);
  const sheet = workbook.Sheets[workbook.SheetNames[0]];
  // logger.debug('sheet', sheet)
  const fields = get_sheet_fields(sheet);
  const docs = XLSX.utils.sheet_to_json(sheet) || [];
  return {
    fields,
    docs,
  };
}

function get_sheet_fields(sheet) {
  const fields = [];
  if (!sheet) return fields;
  let blanks = 0;
  let i = 0;
  while (blanks <= 3) {
    // A1, B1, C1, ...
    const col =
      i
        .toString(26)
        .split("")
        .map((o) => BASE26[o])
        .join("") + "1";
    i += 1;
    const field = sheet[col];
    if (!field) {
      blanks += 1;
      continue;
    }
    const value = field.w || field.v;
    if (!value) {
      blanks += 1;
      continue;
    }
    blanks = 0;
    fields.push(value);
  }
  return fields;
}

const BASE26 = {
  0: "A",
  1: "B",
  2: "C",
  3: "D",
  4: "E",
  5: "F",
  6: "G",
  7: "H",
  8: "I",
  9: "J",
  a: "K",
  b: "L",
  c: "M",
  d: "N",
  e: "O",
  f: "P",
  g: "Q",
  h: "R",
  i: "S",
  j: "T",
  k: "U",
  l: "V",
  m: "W",
  n: "X",
  o: "Y",
  p: "Z",
};
