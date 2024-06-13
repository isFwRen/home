<template>
	<div class="user">
		<div class="z-flex flex-column mb-8 filters">
			<div class="z-flex">
				<div
					v-for="(item, index) in cells.filterFields"
					:key="`prescriptionFilters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'text'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:class="item.class"
							:hideDetails="item.hideDetails"
							:label="item.label"
							:width="item.width"
						>
						</z-text-field>
					</template>

					<template v-else-if="item.inputType === 'select'">
						<z-select
							:formId="searchFormId"
							:formKey="item.formKey"
							:class="item.class"
							:hideDetails="item.hideDetails"
							:label="item.label"
							:width="item.width"
							:options="item.options"
							:defaultValue="item.defaultValue"
						></z-select>
					</template>

					<template v-else>
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:class="item.class"
							:hideDetails="item.hideDetails"
							:label="item.label"
							:range="item.range"
							:width="item.width"
						></z-date-picker>
					</template>
				</div>

				<z-btn class="mt-3" color="primary" @click="onSearch">
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</z-btn>
				<z-btn class="mt-3 ml-3" color="primary" @click="onSync">
					<v-icon class="text-h6">mdi-cloud-sync-outline</v-icon>
					同步
				</z-btn>
			</div>

			<div class="z-flex mt-6">
				<!--				<z-file-input-->
				<!--					formId="files"-->
				<!--					formKey="file"-->
				<!--					accept=".xlsx"-->
				<!--					:action="action"-->
				<!--					chips-->
				<!--					class="pr-3 mt-n3"-->
				<!--					clearable-->
				<!--					:headers="fileHeaders"-->
				<!--					hide-details-->
				<!--					label="导入"-->
				<!--					multiple-->
				<!--					parcel-->
				<!--					prepend-icon="mdi-file-excel-outline"-->
				<!--					width="250"-->
				<!--					@response="handleImportResponse"-->
				<!--					@click="getUserInfo"-->
				<!--				>-->
				<!--				</z-file-input>-->

				<!--				<z-btn class="mr-3" color="primary" @click="onExport">-->
				<!--					<v-icon class="text-h6">mdi-export</v-icon>-->
				<!--					导出-->
				<!--				</z-btn>-->
			</div>
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
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								<!--								<z-btn-->
								<!--									class="mr-2"-->
								<!--									color="primary"-->
								<!--									depressed-->
								<!--									rounded-->
								<!--									smaller-->
								<!--									@click="onEdit(row)"-->
								<!--									>编辑</z-btn-->
								<!--								>-->

								<z-btn
									color="primary"
									class="mr-2"
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
						v-else-if="item.value === 'status'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<span
								:class="{
									'success--text': row[item.value],
									'error--text': !row[item.value]
								}"
								>{{ row[item.value] | chineseStatus }}</span
							>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.date"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<span
								:class="{
									'primary--text': item.value === 'entryDate',
									'success--text': item.value === 'mountGuardDate',
									'error--text': item.value === 'leaveDate'
								}"
								>{{ row[item.value] | dateFormat("YYYY-MM-DD") }}</span
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

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				@page="handlePage"
			></z-pagination>
		</div>

		<!-- 导出 BEGIN -->
		<z-dynamic-form
			ref="exportDynamic"
			:config="{
				proCode: {
					items: [{ label: '全部', value: 'all' }, ...auth.proItems]
				}
			}"
			:fieldList="cells.exportDialogFields"
			:width="600"
			@confirm="handleExportConfirm"
		></z-dynamic-form>
		<!-- 导出 END -->

		<!-- 编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:detail="detailInfo"
			:fieldList="cells.dialogFields"
			:width="600"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 编辑 END -->

		<permission-dialog
			ref="permission"
			:rowInfo="detailInfo"
			@submitted="getList"
		></permission-dialog>

		<lp-spinners :overlay="overlay">
			<h3 slot="tips">同步中，请耐心等待...</h3>
		</lp-spinners>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();
const action = `${baseURLApi}sys-base/user-management/sys-pro-permission/import`;

export default {
	name: "User",
	mixins: [TableMixins],

	data() {
		return {
			formId: "User",
			dispatchList: "STAFF_GET_USER_LIST",
			cells,
			action,
			fileHeaders: {},
			overlay: false
		};
	},

	computed: {
		...mapGetters(["auth"])
	},

	created() {
		this.getUserInfo();
		this.getRole();
	},

	methods: {
		// 获取用户信息
		getUserInfo() {
			const user = this.storage.get("user");
			const token = this.storage.get("token");
			const secret = this.storage.get("secret");
			let code = "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			console.log(code, "code");
			this.user = user;

			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"x-code": String(code)
			};
		},

		async getRole() {
			const result = await this.$store.dispatch("GET_STAFF_ROLE_LIST", {
				pageIndex: 1,
				pageSize: 999,
				status: 1
			});
			if (result.code != 200) {
				return;
			}
			this.cells.dialogFields[0].options = [];
			result.data.list.forEach(element => {
				this.cells.dialogFields[0].options.push({
					value: element.ID,
					label: element.name
				});
			});
		},

		async onSync() {
			this.overlay = true;

			const result = await this.$store.dispatch("STAFF_SYNC_USER_LIST");

			this.toasted.dynamic(result.msg, result.code);

			this.overlay = false;
		},

		// 导入结果
		handleImportResponse({ result }) {
			this.toasted.dynamic(result.msg, result.code);
		},

		// 打开导出弹框
		onExport() {
			this.$refs.exportDynamic.open({ title: "导出" });
		},

		// 导出
		async handleExportConfirm({}, form) {
			const body = {
				proCode: form.proCode
			};

			const result = await this.$store.dispatch("STAFF_EXPORT_PERMISSION", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				const urls = baseURLApi + result.data.replace(/^\.\//, "");
				const res = await await this.$store.dispatch("STAFF_EXPORT_EXCEL", urls);
				this.downloadFile(form.proCode, res);
			}
		},
		downloadFile(proCode, file) {
			var anchor = document.createElement("a");
			anchor.download = `${proCode}用户管理表.xlsx`;
			anchor.style.display = "none";

			anchor.href = URL.createObjectURL(file);
			document.body.appendChild(anchor);
			anchor.click();
			document.body.removeChild(anchor);
		},
		onEdit(row) {
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		handleConfirm({}, form) {
			form = {
				roleId: form.roleId,
				id: this.detailInfo.ID,
				updatedAt: this.detailInfo.UpdatedAt
			};
			this.updateListItem(form, "STAFF_UPDATE_USER");
		},

		// onReset() {
		//   this.$modal({
		//       visible: true,
		//       title: '重置提示',
		//       content: '请确认是否要重置密码？',
		//       confirm: () => {}
		//     })
		// },

		onPermission(row) {
			this.getDetail(row);
			this.$refs.permission.onOpen();
		}
	},

	filters: {
		chineseStatus: value => {
			switch (value) {
				case true:
					return "在职";
				case false:
					return "离职";
				default:
					return "-";
			}
		}
	},

	components: {
		"permission-dialog": () => import("./permissionDialog")
	}
};
</script>
