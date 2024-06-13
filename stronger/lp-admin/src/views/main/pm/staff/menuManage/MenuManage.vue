<template>
	<div class="menu-manage">
		<div class="z-flex justify-end pb-2">
			<z-btn
				class="mr-2"
				color="success"
				@click="
					onNew({
						ID: '0',
						name: '不存在'
					})
				"
			>
				<v-icon>mdi-plus</v-icon>
				新增
			</z-btn>
		</div>

		<vxe-table
			ref="xTable"
			resizable
			border="inner"
			:max-height="tableMaxHeight"
			:loading="loading"
			:size="tableSize"
			:stripe="tableStripe"
			:tree-config="{
				transform: true,
				rowField: 'ID',
				parentField: 'parentId'
			}"
			:data="desserts"
			@toggle-tree-expand="toggleExpandChangeEvent"
		>
			<template v-for="item in cells.headers">
				<vxe-column
					v-if="item.value === 'options'"
					:field="item.value"
					:title="item.text"
					:tree-node="item.treeNode"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<div class="z-flex">
							<v-icon class="mr-2" color="success" size="26" @click="onNew(row)">
								mdi-plus-circle
							</v-icon>

							<v-icon class="mr-2" color="primary" size="26" @click="onEdit(row)"
								>mdi-pencil-circle</v-icon
							>

							<v-icon color="error" size="26" @click="onDelete(row)"
								>mdi-delete-circle</v-icon
							>
						</div>
					</template>
				</vxe-column>

				<vxe-column
					v-else-if="item.value === 'menuType'"
					:field="item.value"
					:title="item.text"
					:tree-node="item.treeNode"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<span>{{ row[item.value] | chineseMenuType }}</span>
					</template>
				</vxe-column>

				<vxe-column
					v-else-if="item.value === 'isEnable'"
					:field="item.value"
					:title="item.text"
					:tree-node="item.treeNode"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<span
							:class="{
								'success--text': row[item.value],
								'error--text': !row[item.value]
							}"
							>{{ row[item.value] | chineseEnable }}</span
						>
					</template>
				</vxe-column>

				<vxe-column
					v-else-if="item.value === 'isFrame'"
					:field="item.value"
					:title="item.text"
					:tree-node="item.treeNode"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<span
							:class="{
								'success--text': row[item.value],
								'error--text': !row[item.value]
							}"
							>{{ row[item.value] | chineseFrame }}</span
						>
					</template>
				</vxe-column>

				<vxe-column
					v-else-if="item.value === 'CreatedAt'"
					:field="item.value"
					:title="item.text"
					:tree-node="item.treeNode"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<span>{{ row[item.value] | dateFormat("YYYY-MM-DD HH:mm:ss") }}</span>
					</template>
				</vxe-column>

				<vxe-column
					v-else
					:field="item.value"
					:title="item.text"
					:tree-node="item.treeNode"
					:key="item.value"
					:width="item.width"
				>
				</vxe-column>
			</template>
		</vxe-table>

		<update-menu-dialog
			ref="updateMenu"
			:rowInfo="detailInfo"
			@submitted="getList"
		></update-menu-dialog>
	</div>
</template>

<script>
import { R } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "MenuManage",
	mixins: [TableMixins],

	data() {
		return {
			dispatchList: "GET_STAFF_MENU_MANAGE_TREE",
			dispatchDelete: "DELETE_STAFF_MENU_MANAGE_TREE_LEAF",
			cells
		};
	},

	created() {
		this.getAllAPI();
	},

	methods: {
		// 新增
		onNew(row) {
			const detail = {
				parentId: row.ID,
				parentName: row.name
			};

			this.getDetail(detail);
			this.$refs.updateMenu.onOpen({ status: -1, title: "新增" });
		},

		// 编辑
		onEdit(row) {
			const result = R.find(this.sabayon.data, { ID: row.parentId });

			const detail = {
				...row,
				parentName: result ? result.name : "不存在"
			};

			this.getDetail(detail);
			this.$refs.updateMenu.onOpen({ status: 1, title: "编辑" });
		},

		onDelete(row) {
			this.getDetail(row);
			this.deleteItem();
		},

		toggleExpandChangeEvent({ row, expanded }) {
			const $table = this.$refs.xTable;
			console.log("节点展开事件", expanded, "获取父节点：", $table.getParentRow(row));
		},

		async getAllAPI() {
			const result = await this.$store.dispatch("GET_STAFF_MENU_MANAGE_API");
			const list = [];

			if (result.code === 200) {
				if (R.isYummy(result.data)) {
					for (let item of result.data) {
						list.push({
							label: item.title,
							value: item.title,
							id: item.ID,
							action: item.action,
							path: item.path
						});
					}
				}
			}

			this.$store.commit("UPDATE_STAFF", { titleOptions: list });
		}
	},

	filters: {
		chineseMenuType: value => {
			switch (value) {
				case "1":
					return "菜单";
				case "2":
					return "按钮";
				default:
					return "-";
			}
		},

		chineseEnable: value => {
			switch (value) {
				case true:
					return "启用";
				case false:
					return "停用";
				default:
					return "-";
			}
		},

		chineseFrame: value => {
			switch (value) {
				case true:
					return "是";
				case false:
					return "否";
				default:
					return "-";
			}
		}
	},

	components: {
		"update-menu-dialog": () => import("./updateMenuDialog")
	}
};
</script>
