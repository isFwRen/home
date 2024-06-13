import { request } from "@/api/service";

const actions = {
	// 获取角色列表
	async GET_STAFF_ROLE_LIST({}, params) {
		console.log(params);

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			status: params.status
		};

		const result = await request({
			url: "sys-base/role-management/list",
			params: data
		});

		return result;
	},

	// 添加或修改角色
	async UPDATE_STAFF_ROLE_ITEM({}, form) {
		const { status } = form;
		console.log("UPDATE_STAFF_ROLE_ITEM", status);
		var data = {};
		if (status === -1) {
			data = {
				id: form.id,
				beforeName: form.beforeName,
				roleName: form.name,
				roleStatus: form.roleStatus,
				roleRemark: form.remark,
				createdBy: form.updatedBy
			};
		} else {
			data = {
				id: form.id,
				beforeName: form.beforeName,
				newName: form.name,
				roleStatus: form.roleStatus,
				roleRemark: form.remark,
				updatedBy: form.updatedBy
			};
		}
		console.log("data", data);

		const result = await request({
			method: "POST",
			url: `sys-base/role-management/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// 删除角色
	async DELETE_STAFF_ROLE_ITEM({}, ids) {
		const result = await request({
			method: "DELETE",
			url: "sys-base/role-management/delete",
			data: {
				ids
			}
		});

		return result;
	},

	// 获取角色权限树
	async GET_STAFF_ROLE_AUTH_TREE({}, body) {
		const data = {
			roleId: body.roleId
		};

		const result = await request({
			url: "sys-base/role-management/sys-menu/tree",
			params: data
		});

		return result;
	},

	// 更新角色权限树节点
	async UPDATE_STAFF_ROLE_AUTH_TREE_LEAF({}, body) {
		const data = {
			list: body
		};
		//   const data = [{
		//     ID: body.id,
		//     isSelect: body.isSelect,
		//     menuId: body.menuId,
		//     roleId: body.roleId
		//   },

		//   {
		//     ID: body.id,
		//     isSelect: body.isSelect,
		//     menuId: body.menuId,
		//     roleId: body.roleId
		//   }
		// ]

		const result = await request({
			method: "POST",
			url: "sys-base/role-management/sys-menu/relation-set",
			data
		});

		return result;
	}
};

export default {
	actions
};
