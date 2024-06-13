const fs = require("node:fs/promises");
const { statSync, rmSync } = require("node:fs");
const util = require("node:util");
const Table = require("nedb/lib/datastore");
const Cursor = require("nedb/lib/cursor");
const Async = require("async");
const lodash = require("lodash");
const picomatch = require("picomatch");
// const gen_id = require("ulid").monotonicFactory();
const crypt = require("./crypt");
const logger = require("./logger");
const sleep = require("./sleep");

const {
  DURIAN_DIR,
  DURIAN_DB_SECRET,
  DURIAN_DB_TABLE_LOAD_MAX,
  DURIAN_MODE,
} = require("./env");
const { uniq } = require("lodash");

const recent_colls = [];

async function load(dbname, options = {}) {
  const proj_dir = `${DURIAN_DIR}/db/${dbname}`;
  await fs.mkdir(proj_dir, { recursive: true });
  dbname = dbname.toUpperCase();
  const self = {
    _name: dbname,
    _dir: proj_dir,
  };
  self.listTables = async () => {
    const inited = self._colls;
    const coll_meta = self.table("META");
    const docs = await coll_meta
      .findEx(
        {},
        {
          sort: { name: 1 },
        }
      )
      .toArrayEx();
    if (!inited) {
      docs.map((doc) => self.table(doc.name));
    }
    return docs;
  };

  self.tableExists = self.getTableMeta = async (coll_name) => {
    return await self
      .table("META")
      .findOneEx({ name: coll_name.toUpperCase() });
  };

  // 校验文件
  self._check_file = function (coll_name) {
    const meta_file = `${proj_dir}/${coll_name}`;
    try {
      var stats = statSync(meta_file);
    } catch (error) {}
    if (stats && stats.isFile() && stats.size < 64) {
      if (coll_name === "META") {
        if (DURIAN_MODE !== "BACKEND") {
          logger.warn("删除破损的文件", meta_file, stats);
          try {
            rmSync(meta_file, { force: true });
          } catch (error) {}
          logger.error("请立即重新运行客户端");
        }
      }
    }
    // 文件不存在 META数量改为 0
    // if (!stats || !stats.isFile()) {
    //   if (coll_name !== "META") {
    //     logger.warn("数据表文件不存在", `${self._name}.${coll_name}`);
    //     const coll_meta = self.table("META");
    //     coll_meta.allowEdit = true;
    //     coll_meta.updateEx(
    //       { name: coll_name },
    //       {
    //         $set: { total: 0 },
    //       }
    //     );
    //   }
    // }
  };

  // 加载期间禁止更新
  self._protect_when_loading = function (coll) {
    coll._loading = true;
    if (DURIAN_MODE === "BACKEND" || coll._name === "META") return;
    coll.persistence.persistCachedDatabase_org =
      coll.persistence.persistCachedDatabase;
    coll.persistence.persistCachedDatabase = (cb) => {
      if (coll.allowEdit) return coll.persistence.persistCachedDatabase_org(cb);
      if (!coll._loading)
        logger.error("persistCachedDatabase", `${self._name}.${coll._name}`);
      if (cb) cb();
    };
    coll.persistence.persistNewState_org = coll.persistence.persistNewState;
    coll.persistence.persistNewState = (newDocs, cb) => {
      if (
        coll.allowEdit ||
        !newDocs ||
        !newDocs[0] ||
        newDocs[0].$$indexCreated
      )
        return coll.persistence.persistNewState_org(newDocs, cb);
      logger.error("persistNewState", `${self._name}.${coll._name}`, newDocs);
      if (cb) cb();
    };
  };

  self.table = function (coll_name) {
    coll_name = coll_name.toUpperCase();
    if (!self._colls) self._colls = {};
    if (self._colls[coll_name]) return self._colls[coll_name];

    // 处理最近加载的数据表
    if (coll_name !== "META") {
      lodash.remove(
        recent_colls,
        (o) => o.coll === coll_name && o.db._name === self._name
      );
      recent_colls.push({
        db: self,
        coll: coll_name,
      });
      // 只保留最近的3个数据表
      while (recent_colls.length > DURIAN_DB_TABLE_LOAD_MAX) {
        const o = recent_colls.shift();
        if (!o.db._colls || !o.db._colls[o.coll]) continue;
        // o.db._colls[o.coll]._unloaded = true;
        delete o.db._colls[o.coll];
      }
    }

    // logger.debug("初始化数据表", `${self._name}.${coll_name}`);
    self._check_file(coll_name);
    const table_options = {
      filename: `${proj_dir}/${coll_name}`,
      autoload: true,
      onload: function (error) {
        if (!error) return;
        if (error.code !== "ENOENT") {
          logger.warn("初始化数据表", error);
        }
      },
    };
    if (DURIAN_DB_SECRET && DURIAN_DB_SECRET !== "") {
      table_options.afterSerialization = crypt.encrypt;
      table_options.beforeDeserialization = crypt.decrypt;
    }
    const coll = (self._colls[coll_name] = new Table(table_options));
    coll._db = self;
    coll._name = coll_name;
    coll._load_at = Date.now();
    self._protect_when_loading(coll);
    coll.checkIndex();
    return coll;
  };

  self.dropCollectionEx = async function (coll_name) {
    if (!(await self.tableExists(coll_name))) return false;
    await self.table(coll_name).dropEx();
    return true;
  };

  self.archive = async function (coll_name) {
    // logger.debug("存档", `${self._name}/${coll_name}`);
    const bak_name = await self._archive(coll_name);
    await self._clean(`${coll_name}.*.old`, {
      reserved: [coll_name],
    });
  };

  self._archive = async function (coll_name) {
    const src_path = `${self._dir}/${coll_name}`;
    if (self._colls) {
      var coll = self._colls[coll_name];
      if (coll) {
        coll._droped = true;
        delete self._colls[coll_name];
      }
    }
    const bak_name = `${coll_name}.${Date.now()}.old`;
    const dst_path = `${self._dir}/${bak_name}`;
    try {
      await fs.copyFile(src_path, dst_path);
    } catch (err) {
      logger.warn(`存档失败 ${bak_name}`);
    }
    logger.debug(`删除数据表文件 ${self._name}.${coll_name}`);
    try {
      await fs.rm(src_path, { force: true });
    } catch (error) {
      logger.warn(`删除数据表文件失败 ${self._name}.${coll_name}`);
    }
    return bak_name;
  };

  // 删除陈旧的备份
  self._clean = async function (old_files, opts = {}) {
    if (!opts.reserved) opts.reserved = [];
    const backups = [];
    const isMatch = picomatch(old_files);
    for (const basename of await fs.readdir(self._dir)) {
      if (opts.reserved.indexOf(basename) >= 0) continue;
      if (isMatch(basename)) backups.push(basename);
    }
    // 保留最多2个
    const old_backups = backups.slice(0, -2);
    if (old_backups.length === 0) return
    logger.debug(`删除 ${self._name} 旧备份`, old_backups);
    for (const basename of old_backups) {
      try {
        await fs.rm(`${self._dir}/${basename}`, { force: true });
      } catch (err) {}
    }
  };

  // 删库，存档所有的表，META排在最后
  self.dropEx = async function () {
    const coll_names = await fs.readdir(self._dir);
    for (const coll_name of coll_names) {
      if (coll_name === "META" || REGEXP_RESERVED.test(coll_name)) continue;
      await self.archive(coll_name);
    }
    await self.archive("META");
    return;
  };

  // 预加载数据表
  self.preload = async function () {
    for (const info of await self.listTables()) {
      const coll = self.table(info.name);
      if (!coll) continue;
      await sleep(10);
      // logger.debug(
      //   `加载 ${self._name}.${info.name}`
      // );
    }
  };
  await self.preload();
  await sleep(100);
  return self;
}

const REGEXP_RESERVED = /(\.old|\.new|\.bak|\.tmp|~)$/i;

// 创建好用的函数
function init_table() {
  if (Cursor.prototype.toArrayEx) return;
  // Table.prototype.createNewId = function () {
  //   var tentativeId = gen_id();
  //   if (this.indexes._id.getMatching(tentativeId).length > 0) {
  //     tentativeId = this.createNewId();
  //   }
  //   return tentativeId;
  // };
  Table.prototype.findEx = function (filter, options = {}) {
    // 字段选择
    let projection = options.fields || options.projection;
    if (Array.isArray(projection)) {
      const o = {};
      projection.map((field) => (o[field] = 1));
      projection = o;
    }
    // 查询
    let cursor = projection ? this.find(filter, projection) : this.find(filter);
    // 排序
    // tingodb: [["_id", "ascending"]]
    // nedb: { firstField: 1, secondField: -1 }
    if (options.sort) {
      let sort = {};
      if (!Array.isArray(options.sort)) {
        Object.keys(options.sort).map(function (k) {
          sort[k] = options.sort[k];
        });
      } else {
        for (const o of options.sort) {
          if (Array.isArray(o)) sort[o[0]] = o[1];
          else {
            Object.keys(o).map(function (k) {
              sort[k] = o[k];
            });
          }
        }
      }
      cursor = cursor.sort(sort);
    }
    // 分页
    if (options.skip) cursor = cursor.skip(options.skip);
    if (options.limit) cursor = cursor.limit(options.limit);
    return cursor;
  };
  Table.prototype.findOneEx = async function (filter, options = {}) {
    const docs = await this.findEx(filter, {
      ...options,
      limit: 1,
    }).toArrayEx();
    if (!docs) return;
    return docs[0];
  };

  Table.prototype.createIndex = function (idx, options, cb) {
    const self = this;
    if (typeof options === "function") {
      cb = options;
      options = {};
    } else {
      if (options === undefined) {
        options = {};
        cb = (e) => null;
      }
    }
    if (cb === undefined) cb = (e) => null;
    if (!Array.isArray(idx)) idx = [idx];
    let fields = [];
    idx.map(function (o) {
      if (Array.isArray(o)) fields.push(o[0]);
      else fields.push(Object.keys(o));
    });
    fields = uniq(fields);
    if (options.unique) {
      const fieldName = fields[0];
      if (self.indexes[fieldName]) return cb();
      self.ensureIndex({ fieldName, unique: true }, cb);
    } else {
      Async.eachSeries(
        fields,
        function (fieldName, next) {
          if (self.indexes[fieldName]) return next();
          self.ensureIndex({ fieldName }, next);
        },
        cb
      );
    }
  };
  Table.prototype.compactCollection = function (cb) {
    this.once("compaction.done", cb);
    this.persistence.compactDatafile();
  };
  Table.prototype.compactCollectionEx = util.promisify(function (cb) {
    this.compactCollection(cb);
  });
  Table.prototype.dropEx = async function () {
    const coll_name = this._name;
    // 删除 META 相当于删库
    if (coll_name == "META") return await this._db.dropEx();
    try {
      await this._db.archive(coll_name);
    } catch (error) {
      logger.debug(`存档失败 ${coll_name}`);
    }
    const coll_meta = this._db.table("META");
    const r = await coll_meta.removeEx(
      { name: coll_name },
      {
        multi: true,
      }
    );
    await coll_meta.compactCollectionEx();
    return r;
  };
  Table.prototype.archive = async function () {
    await this._db.archive(this._name);
  };
  Table.prototype.checkIndex = async function (force) {
    const coll = this;
    if (force) {
      await _check_index_queue.unshiftAsync({ coll, force });
      return;
    }
    const found = lodash.find(
      _check_index_queue.workersList(),
      ({ data }) => data.coll._name === coll._name
    );
    if (found) return;
    await _check_index_queue.pushAsync({ coll, force });
  };
  const funcs = [
    "insert",
    "update",
    "remove",
    "createIndex",
    "ensureIndex",
    "removeIndex",
    "count",
    "findOne",
    "getCandidates",
  ];
  funcs.map((k) => {
    if (["count", "findOne", "getCandidates"].indexOf(k) < 0) {
      // 防止意外修改
      const old_func_name = `${k}Org`;
      Table.prototype[old_func_name] = Table.prototype[k];
      Table.prototype[k] = function (...p) {
        if (this._db.readonly) {
          logger.warn("prevent modify:", p);
          if (["update", "insert", "remove"].indexOf(k) >= 0) return 0;
          return;
        }
        return this[old_func_name](...p);
      };
    }
    if (!Table.prototype[`${k}Ex`])
      Table.prototype[`${k}Ex`] = util.promisify(function (...p) {
        return this[k](...p);
      });
  });
  Cursor.prototype.toArrayEx = util.promisify(function (cb) {
    this.exec(cb);
  });
}

// 检查索引
_check_index = async function ({ coll, force = false }) {
  if (coll._droped) return false;
  const self = coll._db;
  const ident = `${self._name}.${coll._name}`;
  if (!force && coll.indexReady) return false;
  coll.indexReady = true;
  const doc = await coll.findOneEx({});
  if (coll._droped) return false;
  coll._loading = false;
  if (DURIAN_MODE !== "BACKEND") return true;
  if (coll._name === "META") {
    if (!coll.indexes.name)
      coll.createIndex({ name: 1 }, { unique: true }, (e) => null);
    if (!coll.indexes.hash) coll.createIndex({ hash: 1 }, {}, (e) => null);
    return true;
  }
  if (!self._colls.META) await self.table("META").findOneEx({});
  if (coll._droped) return false;
  let need_fix = false;
  let info = await self.table("META").findOneEx({ name: coll._name });
  if (coll._droped) return false;
  if (!info) {
    logger.warn("META中未找到数据表", ident);
    need_fix = true;
    info = {
      name: coll._name,
    };
  }
  if (!info.header || info.header.length === 0) {
    if (!doc) {
      logger.error("数据表为空，无法还原META中缺失的header", ident);
    } else {
      need_fix = true;
      info.header = Object.keys(doc).filter((o) => o !== "_id");
      logger.warn("补全表头", ident, info.header);
    }
  } else {
    if (doc) {
      const header = lodash.union(
        info.header,
        Object.keys(doc).filter((o) => o !== "_id")
      );
      if (header.length > info.header.length) {
        need_fix = true;
        info.header = header;
        logger.warn("补全表头", ident, info.header);
      }
    }
  }
  if (need_fix) {
    delete info._id;
    await self.table("META").updateEx(
      { name: coll._name },
      {
        $set: {
          ...info,
          updatedAt: Date.now(),
        },
      },
      { upsert: true }
    );
  }
  const idx = (info.header || [])
    .map((field) => {
      if (coll.indexes[field]) return;
      return [field, 1];
    })
    .filter((o) => o);
  if (idx.length > 0) {
    logger.warn(
      "补全索引",
      ident,
      idx.map((o) => o[0])
    );
    await coll.createIndexEx(idx);
  }
  await sleep(1);
  return true;
};
const _check_index_queue = Async.queue(async (param) => {
  const ident = `${param.coll._db._name}.${param.coll._name}`;
  logger.debug(
    "开始校验",
    ident,
    _check_index_queue.length() + _check_index_queue.running()
  );
  try {
    var done = await _check_index(param);
  } catch (error) {
    logger.error("校验索引", ident, error);
  } finally {
    setImmediate(
      (ident, done) => {
        logger.debug(
          done ? "完成校验" : "跳过校验",
          ident,
          _check_index_queue.length() + _check_index_queue.running()
        );
      },
      ident,
      done
    );
  }
}, 16);
_check_index_queue.drain(() => {
  setTimeout(() => {
    if (_check_index_queue.length() + _check_index_queue.running() === 0)
      logger.debug("校验队列空闲");
  }, 10);
});

init_table();

module.exports = load;
