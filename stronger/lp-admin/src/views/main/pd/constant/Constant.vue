<template>
	<div class="constant">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="code"
						hideDetails
						label="项目编码"
						:options="auth.proItems"
						:defaultValue="project.code"
						@change="changePro"
					></z-select>
				</v-col>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="dbName"
						hideDetails
						label="常量表"
						:options="constItems"
						@change="changeConst"
					></z-select>
				</v-col>

				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="name"
						clearable
						hideDetails
						label="名称"
						@enter="onSearch"
					></z-text-field>
				</v-col>

				<z-btn
					class="pb-3 pl-3"
					color="primary"
					:disabled="pending"
					:loading="pending"
					@click="onSearch"
				>
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</z-btn>
			</v-row>

			<div class="z-flex pt-6 btns">
				<z-file-input
					formId="files"
					formKey="file"
					accept=".xlsx"
					:action="action"
					chips
					class="pr-3 mt-n3"
					clearable
					:effectData="{ proId: project.code }"
					:headers="fileHeaders"
					hide-details
					label="浏览"
					multiple
					parcel
					prepend-icon="mdi-file-excel-outline"
					width="250"
					@response="handleResponse"
				>
				</z-file-input>

				<z-btn :formId="formId" class="pr-3" color="primary" @click="onExportExcel">
					<v-icon class="text-h6">mdi-export-variant</v-icon>
					导出
				</z-btn>

				<z-btn
					:formId="formId"
					class="pr-3"
					color="error"
					:disabled="!isDeleteMore"
					@click="onDeleteMore"
				>
					<v-icon class="text-h6">mdi-trash-can-outline</v-icon>
					批量删除
				</z-btn>

				<z-btn :formId="formId" class="pr-3" color="error" @click="onDestroy">
					<v-icon class="text-h6">mdi-delete-off-outline</v-icon>
					删除当前常量表
				</z-btn>
			</div>
		</div>

		<div class="table constant-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:loading="pending"
				:max-height="400"
				:size="tableSize"
				@checkbox-all="handleSelectAll"
				@checkbox-change="handleSelectChange"
			>
				<vxe-column type="checkbox" width="60"></vxe-column>

				<template v-for="(item, index) in headers">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="index"
						:width="item.width"
					>
						<template #default="{ row, columnIndex }">
							<span>{{ row.arr[columnIndex - 1] }}</span>
						</template>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:options="pageSizes"
				:pageNum="params.pageIndex"
				:total="pagination.total"
				@page="handlePage"
			></z-pagination>
		</div>

		<!-- 新增/编辑 BEGIN -->
		<edit-dialog
			ref="editDialog"
			:rowInfo="detailInfo"
			:headers="headers"
			@submitted="updateConstant"
		></edit-dialog>
		<!-- 新增/编辑 END -->
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import ButtonMixins from "@/mixins/ButtonMixins";
import SocketsConstMixins from "@/mixins/SocketsConstMixins";
import cells from "./cells";
import { tools as lpTools } from "@/libs/util";

const { baseURLApi } = lpTools.baseURL();
const action = `${baseURLApi}pro-config/project-cost-management/putConstTableByExcel`;

export default {
	name: "Constant",
	mixins: [TableMixins, ButtonMixins, SocketsConstMixins],

	data() {
		return {
			formId: "constant",
			cells,
			action,
			firstSearch: true,
			user: {},
			obsreveForm: "",
			fileHeaders: {},

			dbName: "",
			constItems: [],

			headers: [],

			socketPath: "constReply"
		};
	},
	watch: {
		obsreveForm: function (val) {
			this.firstSearch = val ? false : true;
		}
	},
	created() {
		this.getUserInfo();
	},

	beforeDestroy() {
		this.$store.commit("SET_PROJECT_INFO");
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

			this.user = user;

			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"x-code": String(code)
			};
		},

		// 上传常量表后返回的结果
		uploadedExcel(result) {
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				const { files } = result;
				const filenameList = [];

				for (let file of files) {
					filenameList.push(file.name);
				}

				const data = JSON.stringify({
					proCode: this.project.code,
					relationship: result.data.dbname
				});

				this.$socket.emit("save", data);
			}
		},

		// 检索
		onSearch() {
			const form = this.forms[this.searchFormId];
			this.obsreveForm = form?.name;
			this.params = {
				...this.params,
				pageSize: 10,
				pageIndex: 1,
				...form
			};

			this.desserts = [];
			this.getConstantExcel();
		},

		// 分页
		handlePage({ pageSize, pageNum: pageIndex }) {
			this.params = {
				...this.params,
				pageSize,
				pageIndex
			};

			this.desserts = [];

			this.getConstantExcel();
		},

		// 获取当前常量表 headers
		async getConstantExcelHeaders() {
			this.headers = [];

			const result = await this.$store.dispatch("GET_CONSTANT_HEADERS", this.dbName);

			const list = result.list;

			for (let value of list) {
				this.headers.push({
					text: value,
					value
				});
			}
		},

		// 获取当前常量表
		async getConstantExcel() {
			this.isPending(true);

			const result = await this.$store.dispatch("GET_CONSTANT_EXCEL", this.params);

			this.desserts = result.data;
			this.pagination.total = this.firstSearch ? result.total : result.data.length;

			this.isPending(false);
		},

		// 导出常量表
		async onExportExcel() {
			const body = {
				code: this.user.code,
				proCode: this.project.code,
				dbName: this.dbName
			};

			const result = await this.$store.dispatch("EXPORT_CONSTANT_EXCEL", body);

			this.toasted.dynamic(result.msg, result.code);
		},

		// 删除当前常量表
		async onDestroy() {
			this.$modal({
				visible: true,
				title: "删除提示",
				content: "请确认是否要删除当前常量表？",
				confirm: async () => {
					const result = await this.$store.dispatch(
						"DESTROY_CONSTANT_EXCEL",
						this.dbName
					);

					this.toasted.dynamic(result.msg, result.code);

					if (result.code === 200) {
						this.getConstantOptions();

						this.desserts = [];
						this.pagination.total = 0;
					}
				}
			});
		},

		// 打开编辑弹框
		openEditDialog(row) {
			this.getDetail(row);

			this.$refs.editDialog.onOpen(1);
		},

		// 更新当前项
		async updateConstant(doc) {
			const body = {
				dbName: this.dbName,
				doc
			};

			const result = await this.$store.dispatch("UPDATE_CONSTANT_EXCEL_ITEM", body);

			this.toasted.dynamic(result.msg, result.ok);

			if (result.ok) {
				this.getConstantExcel();
				this.$refs.editDialog.onClose();
			}
		},

		// 更多
		onMore({ customValue }, row) {
			this.getDetail(row);

			if (customValue === "delete") {
				this.$modal({
					visible: true,
					title: "删除提示",
					content: "请确认是否要删除？",
					confirm: async () => {
						const { ok, msg } = await this.$store.dispatch(
							"DELETE_CONSTANT_EXCEL_ITEM",
							{
								dbName: this.dbName,
								del: {
									_id: row._id,
									_rev: row._rev
								}
							}
						);

						this.toasted.dynamic(msg, ok);

						if (ok) {
							this.getConstantExcel();
						}
					}
				});
			}
		},

		// 批量删除
		onDeleteMore() {
			this.$modal({
				visible: true,
				title: "批量删除提示",
				content: "请确认是否要批量删除？",
				confirm: async () => {
					let dels = [];

					for (let { _id, _rev } of this.selected) {
						dels = [...dels, { _id, _rev }];
					}

					const result = await this.$store.dispatch("DELETE_CONSTANT_EXCEL_ITEM", {
						dbName: this.dbName,
						del: dels
					});

					this.toasted.dynamic(result.msg, result.code);

					if (result.code === 200) {
						this.selected = [];
						this.getConstantExcel();
					}
				}
			});
		},

		// 切换项目编码
		changePro(value) {
			this.params = {
				...this.params,
				pageIndex: 1
			};

			this.$store.commit("SET_PROJECT_INFO", { code: value });
			const secret = this.storage.get("secret");
			let code = "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			this.fileHeaders = {
				...this.fileHeaders,
				"pro-code": this.project.code,
				"x-code": String(code)
			};

			this.getConstantOptions();
		},

		// 常量表信息
		async getConstantOptions() {
			this.constItems = [];

			const result = await this.$store.dispatch("GET_CONSTANT_OPTIONS", this.project.code);

			if (result.code === 200) {
				if (tools.isYummy(result.data.list)) {
					for (let item of result.data.list) {
						this.constItems.push({
							label: item.chineseName,
							value: item.dbName
						});
					}
				} else {
					this.toasted.warning("常量表为空！");
				}
			} else {
				this.toasted.error(result.msg);
			}

			this.dbName = tools.isYummy(this.constItems) ? this.constItems[0].value : "";

			this.params = {
				...this.params,
				dbName: this.dbName
			};

			// this.getConstantExcel()
		},

		// 切换常量表信息
		changeConst(value) {
			this.dbName = value;

			this.params = {
				...this.params,
				pageIndex: 1,
				dbName: this.dbName
			};

			this.getConstantExcelHeaders();

			this.getConstantExcel();
		},

		handleResponse({ result }) {
			if (result.code === 200) {
				this.toasted.success("上传常量表成功，等待服务器保存.");
			} else {
				this.toasted.error("上传常量表失败!");
			}
		}
	},

	computed: {
		...mapGetters(["auth", "project"])
	},

	components: {
		"edit-dialog": () => import("./editDialog")
	}
};
</script>
