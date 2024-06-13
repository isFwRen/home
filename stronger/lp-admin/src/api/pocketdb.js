import { tools } from "vue-rocket";
import PouchDB from "pouchdb";
import PouchDBFind from "pouchdb-find";

PouchDB.plugin(PouchDBFind);

class PocketDB {
	constructor(dbName, options) {
		this.db = new PouchDB(dbName, options);
	}

	/**
	 * @description 数据库信息
	 * @see https://pouchdb.com/api.html#database_information
	 */
	info() {
		return new Promise((resolve, reject) => {
			this.db.info((err, info) => {
				try {
					if (!!err) {
						resolve({ ...err, ok: false });
					} else {
						const { db_name, doc_count, update_seq, host } = info;

						resolve({
							ok: true,
							dbName: db_name,
							total: doc_count,
							updateSeq: update_seq,
							host
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	// 删除数据库
	destroy() {
		return new Promise((resolve, reject) => {
			this.db.destroy((err, res) => {
				try {
					if (err) {
						resolve({
							ok: false,
							...err,
							msg: "删除数据库失败！"
						});
					} else {
						resolve({
							ok: true,
							msg: "删除数据库成功！"
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	/**
	 * @description 所有数据
	 * @see https://pouchdb.com/api.html#query_index
	 * @param {object} options 条件
	 */
	allDocs(options) {
		const assignOptions = Object.assign(
			{
				include_docs: true
			},
			options
		);

		return new Promise((resolve, reject) => {
			this.db.allDocs(assignOptions, (err, docs) => {
				try {
					if (err) {
						resolve({
							ok: false,
							...err,
							msg: "文档获取失败！"
						});
					} else {
						const { rows: data, total_rows: total } = docs;
						resolve({
							ok: true,
							total,
							data,
							msg: "文档获取成功！"
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	get(_id) {
		return new Promise((resolve, reject) => {
			this.db.get(_id, (err, res) => {
				try {
					if (err) {
						resolve({
							ok: false,
							...err
						});
					} else {
						resolve({
							ok: true,
							...res
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	put(doc) {
		return new Promise((resolve, reject) => {
			this.db.put(doc, (err, res) => {
				try {
					if (err) {
						resolve({
							ok: false,
							...err,
							msg: "文档更新失败！"
						});
					} else {
						resolve({
							ok: true,
							msg: "文档更新成功！"
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	/**
	 * @description 批量保存数据
	 * @see https://pouchdb.com/api.html#batch_create
	 * @param {object} docs 文档
	 * @param {object} options 条件
	 */
	bulkDocs(docs, options) {
		const assignOptions = Object.assign(
			{
				new_edits: false
			},
			options
		);

		return new Promise((resolve, reject) => {
			this.db.bulkDocs(docs, assignOptions, (err, res) => {
				try {
					if (err) {
						resolve({
							ok: false,
							...err,
							msg: "文档保存失败！"
						});
					} else {
						resolve({
							ok: true,
							msg: "文档保存成功！"
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	/**
	 * @description 查询
	 * @see https://pouchdb.com/api.html#query_index
	 * @param {object} options 查询条件
	 */
	async find(options) {
		const assignOptions = Object.assign(
			{
				selector: {},
				// limit: 10,
				skip: 0
			},
			options
		);
		return new Promise((resolve, reject) => {
			this.db.find(assignOptions, (err, res) => {
				try {
					if (err) {
						resolve({
							ok: false,
							...err,
							msg: "文档查询失败！"
						});
					} else {
						resolve({
							ok: true,
							data: res.docs,
							msg: "文档查询成功！"
						});
					}
				} catch (error) {
					reject(error);
				}
			});
		});
	}

	/**
	 * @description 删除
	 * @see https://pouchdb.com/api.html#delete_document
	 * @param {object | array} info 删除条件
	 */
	remove(info) {
		return new Promise((resolve, reject) => {
			if (tools.isObject(info)) {
				const { _id, _rev } = info;

				this.db.remove(_id, _rev, (err, res) => {
					try {
						if (err) {
							resolve({
								ok: false,
								...err,
								msg: "删除失败！"
							});
						} else {
							resolve({
								ok: true,
								msg: "删除成功！"
							});
						}
					} catch (error) {
						reject(error);
					}
				});
			} else {
				let docs = [];
				for (let item of info) {
					docs = [...docs, { ...item, _deleted: true }];
				}

				this.db.bulkDocs(docs, (err, res) => {
					try {
						if (err) {
							resolve({
								ok: false,
								...err,
								msg: "批量删除失败！"
							});
						} else {
							resolve({
								ok: true,
								msg: "批量删除成功！"
							});
						}
					} catch (error) {
						reject(error);
					}
				});
			}
		});
	}
}

export default PocketDB;
