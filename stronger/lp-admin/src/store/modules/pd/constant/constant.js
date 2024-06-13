import { request } from "@/api/service";
import { tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

import PocketDB from "@/api/pocketdb";

const { baseURL: constBaseURL } = lpTools.constBaseURL();

const actions = {
	// 项目编码
	async GET_PROJECT_OPTIONS() {
		const result = await request({
			url: "pro-config/sys-project/list"
		});

		return result;
	},

	// 常量表信息
	async GET_CONSTANT_OPTIONS({}, proCode) {
		const result = await request({
			url: "pro-config/project-cost-management/list-const-name",
			params: {
				proCode
			}
		});

		return result;
	},

	// 获取表头
	async GET_CONSTANT_HEADERS({}, dbName) {
		const data = {
			dbName: dbName
		};
		const result = await request({
			url: "pro-config/project-cost-management/list-table-top",
			params: data
		});
		return result;
	},

	// 获取常量表
	async GET_CONSTANT_EXCEL({}, { pageIndex, pageSize, dbName, name }) {
		// 常量表名为空
		if (tools.isLousy(dbName)) {
			this.toasted.warning("常量表名不能为空！");
			return {
				data: [],
				total: 0
			};
		}

		const remoteDB = new PocketDB(`${constBaseURL}${dbName}`, { skip_setup: true });
		const remoteInfo = await remoteDB.info();

		// 远程不存在该数据库
		if (!remoteInfo.total && !remoteInfo.updateSeq) {
			this.toasted.warning("当前项目无常量表！");

			return {
				data: [],
				total: 0
			};
		}

		function setRegex() {
			if (!name) return "";
			const arr = name.split(" ");
			let str = "";
			for (let val of arr) {
				str += `(?=.*${val})`;
			}

			return `^${str}.+$`;
		}

		// 按条件查询
		let result = null;
		if (setRegex() !== "") {
			result = await remoteDB.find({
				skip: (pageIndex - 1) * pageSize,
				selector: {
					arr: {
						$elemMatch: { $regex: setRegex() }
					}
				}
			});
		} else {
			// console.log({
			//   limit: pageSize,
			//   skip: (pageIndex - 1) * pageSize,
			//   selector: {
			//     arr: {
			//       $elemMatch: { $regex: setRegex() }
			//     }
			//   }
			// }, 'couchDB')
			result = await remoteDB.find({
				limit: pageSize,
				skip: (pageIndex - 1) * pageSize,
				selector: {
					_id: {
						$gt: null
					}
				}
			});
		}

		const desserts = tools.deepClone(result.data);

		// 删除不需要的数据
		for (let i = 0; i < desserts.length; i++) {
			if (desserts[i].content || desserts[i].tabletop) {
				desserts.splice(i, 1);
				--i;
			}
		}

		return {
			ok: result.ok,
			msg: result.msg,
			data: desserts,
			total: remoteInfo.total - 3
		};
	},

	// 更新当前常量表
	async UPDATE_CONSTANT_EXCEL_ITEM({}, { dbName, doc }) {
		const remoteDB = new PocketDB(`${constBaseURL}${dbName}`);

		let remoteDoc;

		if (tools.isObject(doc)) {
			const result = await remoteDB.get(doc._id);

			remoteDoc = {
				...doc,
				_rev: result._rev
			};
		}

		let result = await remoteDB.put(remoteDoc);

		if (result.ok) {
			const localDB = new PocketDB(dbName);

			result = await localDB.put(doc);
		}

		return result;
	},

	// 导出常量表
	async EXPORT_CONSTANT_EXCEL({}, { code, proCode, dbName }) {
		const data = {
			exportCode: code,
			proCode,
			constName: dbName
		};

		const result = await request({
			url: "pro-config/project-cost-management/export-const",
			params: data
		});

		return result;
	},

	// 删除当前常量表
	async DESTROY_CONSTANT_EXCEL({}, dbName) {
		const data = {
			dbName: dbName
		};

		const result = await request({
			method: "POST",
			url: "pro-config/project-cost-management/delete-table",
			params: data
		});

		if (result.code === 200) {
			// const remoteDB = new PocketDB(`${ constBaseURL }${ dbName }`)

			// let result = await remoteDB.destroy()

			// if(result.ok) {
			//   const localDB = new PocketDB(dbName)

			//   result = await localDB.destroy()
			// }

			// return result

			const localDB = new PocketDB(dbName);
			await localDB.destroy();
		}

		return result;
	},

	// 批量删除(删除当前常量表下的数据)
	async DELETE_CONSTANT_EXCEL_ITEM({}, { dbName, del }) {
		const data = {
			dbName,
			table: del
		};

		const result = await request({
			method: "POST",
			url: "pro-config/project-cost-management/delete-line",
			data
		});

		if (result.code === 200) {
			const localDB = new PocketDB(dbName);
			await localDB.remove(del);
		}

		return result;

		// const remoteDB = new PocketDB(`${ constBaseURL }${ dbName }`)

		// let remoteDel

		// if(tools.isObject(del)) {
		//   const result = await remoteDB.get(del._id)
		//   remoteDel = {
		//     _id: del._id,
		//     _rev: result._rev
		//   }
		// }
		// else {
		//   remoteDel = []
		//   for(let item of del) {
		//     const result = await remoteDB.get(item._id)
		//     remoteDel = [...remoteDel, { _id: item._id, _rev: result._rev }]
		//   }
		// }

		// let result = await remoteDB.remove(remoteDel)

		// if(result.ok) {
		//   const localDB = new PocketDB(dbName)

		//   result = await localDB.remove(del)
		// }

		// return result
	}
};

export default {
	actions
};
