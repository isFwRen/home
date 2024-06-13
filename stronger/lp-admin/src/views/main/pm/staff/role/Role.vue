<template>
	<div class="role">
		<div class="z-flex align-end mb-8 filters">
			<z-select
				:formId="searchFormId"
				formKey="status"
				class="pr-4"
				hide-details
				label="状态"
				width="160"
				:options="cells.statusOptions"
				:defaultValue="''"
				@change="onSearch"
			>
			</z-select>

			<z-btn color="primary" @click="onNew">
				<v-icon>mdi-plus</v-icon>
				新增
			</z-btn>
		</div>

		<div class="table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:max-height="tableMaxHeight"
				:loading="loading"
				:size="tableSize"
				:stripe="tableStripe"
			>
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								<z-btn
									class="mr-2"
									color="primary"
									depressed
									rounded
									smaller
									@click="onEdit(row)"
									>编辑</z-btn
								>

								<z-btn
									class="mr-2"
									color="primary"
									depressed
									rounded
									smaller
									@click="onDelete(row)"
									>删除</z-btn
								>

								<z-btn
									color="primary"
									depressed
									rounded
									smaller
									@click="onPermission(row)"
									>权限</z-btn
								>
							</div>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'CreatedAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row[item.value] | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'status'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<span
								:class="{
									'success--text': row[item.value] === 1,
									'error--text': row[item.value] === 2
								}"
								>{{ row[item.value] | chineseStatus }}</span
							>
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination :total="pagination.total" :options="pageSizes"></z-pagination>
		</div>

		<!-- 新增/编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:detail="detailInfo"
			:fieldList="cells.fields"
			:formId="formId"
			:width="600"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->

		<role-permission-dialog ref="permission" :rowInfo="detailInfo"></role-permission-dialog>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import { rocket } from "vue-rocket";

export default {
	name: "Role",
	mixins: [TableMixins],

	data() {
		return {
			formId: "Role",
			dispatchList: "GET_STAFF_ROLE_LIST",
			dispatchDelete: "DELETE_STAFF_ROLE_ITEM",
			cells
		};
	},

	methods: {
		onNew() {
			rocket.emit("ZHT_CLEAR_FORM", this.formId);
			this.$refs.dynamic.open({ title: "新增", status: -1 });
		},

		onEdit(row) {
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		async handleConfirm(effect, form) {
			const data = {
				status: effect.status,
				id: this.detailInfo.ID,
				beforeName: this.detailInfo.name,
				name: form.name,
				roleStatus: form.status,
				remark: form.remark,
				updatedBy: this.storage.get("user").code
			};

			const result = await this.$store.dispatch("UPDATE_STAFF_ROLE_ITEM", data);
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
				this.$refs.dynamic.close();
				rocket.emit("ZHT_RESET_FORM", this.formId);
			}
		},

		onDelete(row) {
			this.getDetail(row);
			this.deleteItem();
		},

		onPermission(row) {
			this.getDetail(row);
			this.$refs.permission.onOpen();
		}
	},

	filters: {
		chineseStatus: value => {
			switch (value) {
				case 1:
					return "正常";
				case 2:
					return "停用";
				default:
					return "-";
			}
		}
	},

	components: {
		"role-permission-dialog": () => import("./rolePermissionDialog")
	}
};
</script>
