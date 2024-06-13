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
					<z-select
						:formId="searchFormId"
						formKey="queryKey"
						hideDetails
						label="模糊查询字段"
						:options="queryList"
						:defaultValue="defaultHeaderKey"
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
				<z-btn class="pb-3 pl-3" color="primary" @click="openLog"> 更新日志 </z-btn>
			</v-row>

			<div class="z-flex pt-6 btns">
				<z-file-input
					formId="files"
					formKey="file"
					accept=".xlsx"
					chips
					class="pr-3 mt-n3"
					clearable
					:effectData="{ proCode: project.code }"
					:headers="fileHeaders"
					hide-details
					label="浏览"
					multiple
					parcel
					prepend-icon="mdi-file-excel-outline"
					width="250"
					:autoUpload="false"
					@response="handleResponse"
					@change="handleChange"
					@click="handleClick"
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
				<z-btn :formId="formId" class="pr-3" color="success" @click="onNew">
					<v-icon class="text-h6">mdi-plus-box</v-icon>
					批量新增
				</z-btn>
				<z-btn :formId="formId" class="pr-3" dark color="teal" @click="onPublish">
					<v-icon class="text-h6"> mdi-cloud-upload</v-icon>
					发布
				</z-btn>
				<z-btn :formId="formId" class="pr-3" dark color="teal" @click="onBatchPublish">
					<v-icon class="text-h6"> mdi-cloud-upload</v-icon>
					批量发布
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
						:key="item._id"
						:width="item.width"
					>
					</vxe-column>
				</template>
				<vxe-column width="100" field="operation" title="操作" v-if="headers.length">
					<template #default="{ row }">
						<vxe-button
							type="text"
							icon="vxe-icon-edit"
							@click="openEditDialog(row)"
						></vxe-button>
					</template>
				</vxe-column>
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
		<new-dialog ref="newDialog" :headers="headers" @add="handleAdd"></new-dialog>
		<!-- 新增/编辑 END -->

		<div class="loading" v-if="load">正在上传</div>

		<constant-log ref="constantLog" />
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

export default {
	name: "Constant",
	mixins: [TableMixins, ButtonMixins, SocketsConstMixins],

	data() {
		return {
			formId: "constant",
			cells,
			firstSearch: true,
			user: {},
			obsreveForm: "",
			fileHeaders: {},
			queryList: [],
			dbName: "",
			load: false,
			constItems: [],
			headers: [],
			dbNames: [],
			ConstantOptions: [],
			socketPath: "constReply",
			defaultHeaderKey: "",
			firstSelect: true,
			itemBefore: {}
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
		handleClick(event) {
			if (!this.project.code) {
				this.toasted.warning("请先选择上传项目编码");
				event.preventDefault();
			}
		},
		async handleChange(files) {
			if (files.length === 0) {
				return;
			}
			this.load = true;
			const formData = new FormData();

			const body = {
				formData,
				proCode: this.project.code
			};

			for (let i = 0, len = files.length; i < len; i++) {
				formData.append("files", files[i]);
			}

			const result = await this.$store.dispatch("UPLOAD_CONSTANT_FILE_CLIENT", body);
			if (result.status === 200) {
				this.load = false;
				this.toasted.dynamic(result.msg, result.status);
				this.getConstantOptions();
			}
		},
		async handleAdd(arr) {
			const body = {
				proCode: this.project.code,
				name: this.dbName,
				items: arr
			};
			const result = await this.$store.dispatch("BATCH_ADD_EXCEL_CLIENT", body);

			this.toasted.dynamic(result.msg, result.status);

			if (result.status == 200) {
				this.load = false;
				this.getConstantExcelHeaders();
				this.$refs.newDialog.onClose();
			}
		},
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
		matchSymbol(str) {
			const specialSymbols = [
				"(",
				")",
				"-",
				".",
				"+",
				"=",
				"@",
				"!",
				"^",
				"%",
				"#",
				"$",
				"*",
				"[",
				"]",
				"{",
				"}",
				",",
				";",
				"<",
				">"
			];
			specialSymbols.forEach(reg => {
				str = str.replace(reg, `\\${reg}`);
			});

			return str;
		},
		openLog() {
			this.$refs.constantLog.openDialog();
		},
		queryAttr(form) {
			this.obsreveForm = form?.name;
			let queryNames = {};
			if (form.queryKey && form.name) {
				let str = form.name;
				str = str.trim();
				// 特殊符号在正则中有特殊意义，需要过滤
				str = this.matchSymbol(str);
				const reg = /\s/;
				if (reg.test(" ") && str.split(" ").length > 0) {
					const strs = str.split(" ");
					let regex;
					const len = strs.length;
					for (let i = 0; i < len; i++) {
						regex = i > 0 ? `(${regex})(.)*(${strs[i]})` : strs[i];
					}
					queryNames = {
						[form.queryKey]: {
							$regex: `/${regex}/`
						}
					};
				} else {
					queryNames = {
						[form.queryKey]: {
							$regex: `/${str}/`
						}
					};
				}
			}
			return queryNames;
		},
		// 检索
		onSearch() {
			const form = this.forms[this.searchFormId];
			this.obsreveForm = form?.name;
			this.params = {
				...this.params,
				...form
			};

			if (!form.queryKey) {
				this.toasted.warning("请选择模糊查询字段!");
				return;
			}

			this.desserts = [];
			if (form.name) {
				const queryNames = this.queryAttr(form);
				const body = {
					pageSize: 10,
					pageIndex: 1,
					name: this.dbName,
					queryNames
				};
				// 客户端
				this.getConstantExcelHeaders(body);
			} else {
				const body = {
					pageSize: this.params.pageSize,
					pageIndex: this.params.pageIndex,
					name: this.dbName,
					queryNames: {}
				};
				this.getConstantExcelHeaders(body);
			}
		},

		// 分页
		handlePage({ pageSize, pageNum: pageIndex }) {
			const form = this.forms[this.searchFormId];
			this.params = {
				...this.params,
				pageSize,
				pageIndex,
				...form
			};

			this.headers = [];
			this.desserts = [];
			const queryNames = this.queryAttr(form);
			const body = {
				pageSize,
				pageIndex,
				name: this.dbName,
				queryNames
			};
			// 客户端
			this.getConstantExcelHeaders(body);
		},

		// 获取当前常量表 headers
		//*客户端*
		async getConstantExcelHeaders(params = {}) {
			this.headers = [];
			this.queryList = [];

			const body = {
				name: this.dbName,
				proCode: this.project.code,
				queryNames: {},
				respNames: [],
				...params
			};
			const result = await this.$store.dispatch("GET_CONSTANT_HEADERS_CLIENT", body);

			const list = this.ConstantOptions.find(item => item.name === this.dbName);

			for (let value of list.header) {
				this.headers.push({
					text: value,
					value
				});
				this.queryList.push({
					label: value
				});
				if (value == "项目名称" && this.firstSelect === true) {
					this.defaultHeaderKey = value;
				}
			}

			this.firstSelect = false;
			this.desserts = result.list;
			this.pagination.total = result.total;
		},
		//*客户端*

		// 导出常量表
		//*客户端*
		async onExportExcel() {
			const body = {
				code: this.user.code,
				proCode: this.project.code,
				dbName: this.dbName
			};

			const result = await this.$store.dispatch("EXPORT_CONSTANT_EXCEL_CLIENT", body);

			const fileName = `${this.project.code}常量管理`;
			this.downloadFile(fileName, result);

			this.toasted.dynamic(result.msg, result.status);
		},
		//*客户端*
		downloadFile(fileName, file) {
			var anchor = document.createElement("a");
			anchor.download = `${fileName}.xlsx`;
			anchor.style.display = "none";

			anchor.href = URL.createObjectURL(file);
			document.body.appendChild(anchor);
			anchor.click();
			document.body.removeChild(anchor);
		},
		getConstantIds() {
			let ids = [];
			for (let { _id } of this.selected) {
				ids = [...ids, _id];
			}
			return ids;
		},
		async onPublish() {
			const ids = this.getConstantIds();
			const body = {
				name: [this.dbName],
				proCode: this.project.code
			};
			const result = await this.$store.dispatch("CONSTANT_PUBLISH_CLIENT", body);
			this.toasted.dynamic(result.msg, result.status);
		},
		// 批量发布
		async onBatchPublish() {
			const ids = this.getConstantIds();
			const body = {
				name: this.dbNames,
				proCode: this.project.code
			};
			const result = await this.$store.dispatch("CONSTANT_PUBLISH_CLIENT", body);
			this.toasted.dynamic(result.msg, result.status);
		},
		// 批量新增
		onNew() {
			if (this.headers.length === 0) {
				this.toasted.warning("常量表不能为空");
				return;
			}
			this.$refs.newDialog.onOpen("新增");
		},
		// 删除当前常量表
		//*客户端*
		async onDestroy() {
			this.$modal({
				visible: true,
				title: "删除提示",
				content: "请确认是否要删除当前常量表？",
				confirm: async () => {
					const result = await this.$store.dispatch("DESTROY_CONSTANT_EXCEL_CLIENT", {
						proCode: this.project.code,
						name: [this.dbName]
					});

					this.toasted.dynamic(result.msg, result.status);

					if (result.status === 200) {
						this.getConstantOptions();
						this.desserts = [];
						this.pagination.total = 0;
					}
				}
			});
		},
		//*客户端*

		// 打开编辑弹框
		openEditDialog(row) {
			this.getDetail(row);
			this.itemBefore = row;

			this.$refs.editDialog.onOpen(1);
		},

		// 编辑行
		//*客户端*
		async updateConstant(doc) {
			const body = {
				proCode: this.project.code,
				name: this.dbName,
				item: doc.item,
				id: doc.id,
				itemBefore: this.itemBefore
			};
			const result = await this.$store.dispatch("UPDATE_CONSTANT_EXCEL_ITEM_CLIENT", body);
			this.toasted.dynamic(result.msg, result.status);
			if (result.status === 200) {
				this.getConstantExcelHeaders();
				this.$refs.editDialog.onClose();
				this.itemBefore = {};
			}
		},
		//*客户端*

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
							"DELETE_CONSTANT_EXCEL_ITEM_CLIENT",
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
							// 客户端
							this.getConstantExcelHeaders();
						}
					}
				});
			}
		},

		// 批量删除
		//*客户端*
		onDeleteMore() {
			this.$modal({
				visible: true,
				title: "批量删除提示",
				content: "请确认是否要批量删除？",
				confirm: async () => {
					const ids = this.getConstantIds();

					const items = this.desserts
						.map(dessert => {
							if (ids.includes(dessert._id)) return dessert;
						})
						.filter(item => item);
					const body = {
						proCode: this.project.code,
						ids,
						name: this.dbName,
						items
					};
					console.log(body, "body");
					const result = await this.$store.dispatch(
						"DELETE_CONSTANT_EXCEL_ITEM_CLIENT",
						body
					);

					this.toasted.dynamic(result.msg, result.status);

					if (result.status === 200) {
						this.selected = [];
						this.getConstantExcelHeaders();
					}
				}
			});
		},
		//*客户端*

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
				// "pro-code": this.project.code,
				"x-code": String(code)
			};
			this.getConstantOptions();
		},

		// 常量表信息
		//*客户端*
		async getConstantOptions() {
			this.constItems = [];

			const result = await this.$store.dispatch(
				"GET_CONSTANT_OPTIONS_CLIENT",
				this.project.code
			);

			if (result.status === 200) {
				this.dbNames = result.list.map(item => item.name);
				if (tools.isYummy(result.list)) {
					this.ConstantOptions = result.list;
					for (let item of result.list) {
						this.constItems.push({
							label: item.name,
							id: item.id
						});
					}
				} else {
					this.toasted.warning("常量表为空！");
				}
			} else {
				this.toasted.error(result.msg);
			}

			this.dbName = tools.isYummy(this.constItems) ? this.constItems[0].id : "";

			this.params = {
				...this.params,
				name: this.dbName
			};
		},
		//*客户端*

		// 切换常量表信息
		changeConst(value) {
			this.firstSelect = true;
			this.defaultHeaderKey = "";
			this.dbName = value;
			this.headers = [];

			this.forms[this.searchFormId].queryKey = "";
			const body = {
				pageSize: 10,
				pageIndex: 1,
				name: this.dbName
			};

			this.getConstantExcelHeaders(body);
		},

		handleResponse({ result }) {
			if (result.code === 200) {
				this.toasted.success("上传常量表成功，等待服务器保存.");
			} else {
				//this.toasted.error("上传常量表失败!");
			}
		}
	},

	computed: {
		...mapGetters(["auth", "project"])
	},

	components: {
		"edit-dialog": () => import("./editDialog"),
		"new-dialog": () => import("./newDialog"),
		"constant-log": () => import("./constantLog/index.vue")
	}
};
</script>

<style lang="scss" scoped>
.loading {
	width: 110px;
	position: absolute;
	top: 20px;
	right: 70px;
	background-color: #0ea5e9;
	border: 1px solid #eee;
	text-align: center;
	padding: 5px 5px;
	padding-right: 17px;
	color: #fff;
	border-radius: 5px;
	animation: dot 3s infinite steps(3, start);

	&::after {
		color: #fff;
		content: "";
		position: absolute;
		top: 0%;
		bottom: 0;
		animation: dot 3s infinite steps(3, start);
		line-height: 40px;
	}
}

@keyframes dot {
	33.33% {
		content: ".";
	}

	66.67% {
		content: "..";
	}

	100% {
		content: "....";
	}
}
</style>
