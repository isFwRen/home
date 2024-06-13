<template>
	<lp-dialog
		ref="dialog"
		title="角色权限"
		transition="dialog-bottom-transition"
		@dialog="handleDialog"
	>
		<div class="pt-6" slot="main">
			<div class="z-flex mb-6">
				<z-btn class="mr-4" color="primary" lockedTime="0" outlined small @click="onToggle">
					<template v-if="!openAll">
						<v-icon size="18">mdi-format-indent-increase </v-icon>
						全部展开
					</template>
					<template v-else>
						<v-icon size="18">mdi-format-indent-decrease</v-icon>
						全部收缩
					</template>
				</z-btn>
			</div>

			<v-row>
				<v-col :cols="10" class="pl-14"> 名称 </v-col>

				<v-col :cols="2" class="text-center"> 禁用/启用 </v-col>
			</v-row>

			<v-divider></v-divider>

			<v-treeview v-if="reload" :items="items" :open-all="openAll">
				<template v-slot:append="{ item }">
					<z-switch
						:formId="formId"
						:formKey="`view_${item.id}`"
						class="mr-3"
						:defaultValue="item.isSelect"
						@change="updateLeaf($event, item)"
					></z-switch>
				</template>
			</v-treeview>
		</div>
	</lp-dialog>
</template>

<script>
import { R } from "vue-rocket";
import DialogMixins from "@/mixins/DialogMixins";

const retainKeys = ["isSelect", "id", "name", "roleMenuRelationId", "children"];

export default {
	name: "RolePermissionDialog",
	mixins: [DialogMixins],

	props: {
		rowInfo: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			formId: "RolePermissionDialog",
			openAll: false,
			reload: true,
			items: [],

			treeItems: []
		};
	},

	methods: {
		// 获取树
		async getTree() {
			const data = {
				roleId: this.rowInfo.ID
			};

			const result = await this.$store.dispatch("GET_STAFF_ROLE_AUTH_TREE", data);

			if (result.code === 200) {
				this.items = this.recursion(R.deepClone(result.data) || []);
			}
			console.log("item", this.items);
			this.$refs.dialog.onOpen();
		},

		// 更新树节点
		async updateLeaf(value, leaf) {
			console.log("value", value);
			console.log("leaf", leaf);
			const data = [];

			var dataitem = {
				id: leaf.roleMenuRelationId,
				isSelect: value,
				menuId: leaf.id,
				roleId: this.rowInfo.ID
			};

			data.push(dataitem);

			if (leaf.children.length != 0) {
				for (var i = 0; i < leaf.children.length; i++) {
					leaf.children[i].isSelect = value;
					var dataitem = {
						id: leaf.children[i].roleMenuRelationId,
						isSelect: value,
						menuId: leaf.children[i].id,
						roleId: this.rowInfo.ID
					};

					data.push(dataitem);
				}
			}

			var editId = leaf.id;
			for (var j = 0; j < this.items.length; j++) {
				if (this.items[j].id == editId) {
					this.items[j].isSelect = value;
				}
				for (var k = 0; k < this.items[j].children.length; k++) {
					if (this.items[j].children[k].id == editId) {
						this.items[j].children[k].isSelect = value;
					}
				}
			}
			console.log("data", data);
			const result = await this.$store.dispatch("UPDATE_STAFF_ROLE_AUTH_TREE_LEAF", data);

			this.toasted.dynamic(result.msg, result.code);
		},

		// 去掉不需要的key
		recursion(treeItems) {
			for (let item of treeItems) {
				item.id = item.ID;
				item.name = item.title;

				for (let key in item) {
					if (!retainKeys.includes(key)) {
						delete item[key];
					}
				}

				if (item.children) {
					this.recursion(item.children);
				}
			}

			return treeItems;
		},

		handleDialog(dialog) {
			if (dialog) {
				this.getTree();
			}
		},

		onToggle() {
			this.reload = false;
			this.$nextTick(() => {
				this.openAll = !this.openAll;
				this.reload = true;
			});
		}
	}
};
</script>
