import { request } from "@/api/service";
import { MENU_TYPE, BUTTON_TYPE } from "@/views/main/pm/staff/menuManage/updateMenuDialog/cells";

const state = {
	staff: {
		titleOptions: []
	}
};

const getters = {
	staff: () => state.staff
};

const mutations = {
	UPDATE_STAFF(state, data) {
		const { staff } = state;
		state.staff = Object.assign(staff, data);
	}
};

const actions = {
	// 菜单或按钮树
	async GET_STAFF_MENU_MANAGE_TREE() {
		const result = await request({
			url: "sys-menu/tree"
		});

		return result;
	},

	// 所有api
	async GET_STAFF_MENU_MANAGE_API() {
		const result = await request({
			url: "sys-menu/api/get"
		});

		return result;
	},

	// 新增按钮或菜单
	async UPDATE_STAFF_MENU_MANAGE_TREE_LEAF({}, body) {
		const { status } = body;

		let data = {
			ID: body.ID,
			isFrame: false
		};

		switch (body.menuType) {
			case MENU_TYPE:
				data = {
					...data,
					menuType: body.menuType,
					parentId: body.parentId,
					title: body.title,
					component: body.component,
					icon: body.icon,
					name: body.name,
					path: body.path,
					sort: +body.sort,
					isEnable: body.isEnable
				};
				break;

			case BUTTON_TYPE:
				data = {
					...data,
					menuType: body.menuType,
					parentId: body.parentId,
					title: body.title,
					action: body.action,
					api: body.api,
					apiId: body.apiId,
					icon: body.icon,
					name: body.name,
					sort: +body.sort,
					isEnable: body.isEnable
				};
				break;
		}

		const result = await request({
			method: "POST",
			url: `sys-menu/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// 删除
	async DELETE_STAFF_MENU_MANAGE_TREE_LEAF({}, ids) {
		const result = await request({
			method: "DELETE",
			url: "sys-menu/delete",
			data: {
				ids
			}
		});
		return result;
	}
};

export default {
	state,
	getters,
	mutations,
	actions
};
