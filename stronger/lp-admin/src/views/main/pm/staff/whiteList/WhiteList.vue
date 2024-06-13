<template>
	<div class="white-list">
		<div class="z-flex align-center">
			<z-autocomplete
				:formId="searchFormId"
				formKey="proCode"
				class="mr-4"
				label="项目"
				placeholder="请选择项目"
				width="140"
				:options="proOptions"
				@change="selectPro"
			></z-autocomplete>

			<z-autocomplete
				:formId="searchFormId"
				formKey="tempCode"
				:disabled="!proId"
				label="模板"
				placeholder="请选择模板"
				width="140"
				:options="tempOptions"
				@change="selectTemp"
			></z-autocomplete>
		</div>

		<v-divider></v-divider>

		<div class="z-flex align-end mt-4 mb-4 filters">
			<z-text-field
				:formId="searchFormId"
				formKey="code"
				class="mr-4"
				hint="输入结束后按回车键."
				label="工号"
				:width="120"
				@change="getWhitelist"
			>
			</z-text-field>

			<z-text-field
				:formId="searchFormId"
				formKey="copyCode"
				class="mr-4"
				hint="请输入复制工号，多个工号时用空格分开."
				label="复制工号"
				:width="230"
			>
			</z-text-field>

			<div class="pb-5">
				<z-btn class="mr-3" color="primary" @click="onCopy">
					<v-icon class="text-h6">mdi-content-copy</v-icon>
					复制
				</z-btn>

				<z-btn class="mr-3" color="primary" @click="onSave">
					<v-icon class="text-h6">mdi-content-save</v-icon>
					保存
				</z-btn>

				<z-btn color="primary" @click="onExport">
					<v-icon class="text-h6">mdi-export</v-icon>
					导出
				</z-btn>
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
				<template v-for="(item, index) in headers">
					<!-- checkbox 列 BEGIN -->
					<vxe-column
						v-if="item.value === 'checkbox'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #header>
							<span class="pl-2">#</span>
						</template>

						<template #default="{ rowIndex, row }">
							<z-checkbox
								v-if="row.checkbox"
								:formId="tableSideFormId"
								:formKey="`side_${rowIndex}`"
								class="pa-0 ma-0"
								color="#1976d2"
								:defaultValue="row['checkbox'].selected"
								dense
								hide-details
								:indeterminate="row['checkbox'].indeterminate"
								@change="changeSide($event, { rowIndex })"
							></z-checkbox>
						</template>
					</vxe-column>
					<!-- checkbox 列 END -->

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #header="{ columnIndex }">
							<div class="z-flex align-center">
								<z-checkbox
									:formId="tableHeaderFormId"
									:formKey="`header_${columnIndex}`"
									class="pa-0 ma-0"
									color="#1976d2"
									:defaultValue="item.selected"
									dense
									hide-details
									:indeterminate="item.indeterminate"
									@change="
										changeHeader($event, { columnIndex, headerCell: item })
									"
								></z-checkbox>
								<span>{{ headers[index].text }}</span>
								<span class="ml-1 error--text">
									{{ headers[index].peopleSum }}</span
								>
							</div>
						</template>

						<template #default="{ rowIndex, columnIndex, row }">
							<div v-if="row[item.value]" class="z-flex align-center">
								<z-checkbox
									:formId="tableDessertsFormId"
									:formKey="`cell_${rowIndex}_${columnIndex}`"
									class="pa-0 ma-0"
									color="#1976d2"
									dense
									hide-details
									:defaultValue="row[item.value].selected"
									@change="changeCell($event, { cell: row[item.value] })"
								></z-checkbox>
								<span>{{ row[item.value].code }}{{ row[item.value].name }}</span>
							</div>
						</template>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				@page="handlePage"
			></z-pagination>
		</div>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { R } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import cells from "./cells";

const EXCLUDE = ["checkbox", "_X_ID"];

export default {
	name: "WhiteList",
	mixins: [TableMixins, SocketsMixins],

	data() {
		return {
			formId: "WhiteList",
			tableHeaderFormId: "WhiteListTableHeaders",
			tableSideFormId: "WhiteListSideHeaders",
			tableDessertsFormId: "WhiteListTableDesserts",
			cells,
			// dispatchList: 'GET_STAFF_WHITELIST',
			socketPath: "whitelist",
			proOptions: [],
			proId: "",
			tempOptions: [],
			tempId: "",

			headers: [],
			desserts: []
		};
	},

	created() {
		this.createProOptions();
	},

	methods: {
		// 项目
		createProOptions() {
			const perm = this.auth.perm;
			const options = [];
			for (let item of perm) {
				if (item.hasPm) {
					options.push({
						label: item.proCode,
						value: item.proCode,
						proId: item.proId
					});
				}
			}
			this.proOptions = options;
		},

		// 模板
		async createTempOptions() {
			const params = {
				proId: this.proId
			};
			const result = await this.$store.dispatch("GET_STAFF_WHITELIST_TEMP_OPTIONS", params);

			if (result.code === 200) {
				const options = [];
				for (let item of result.data) {
					options.push({
						label: item.name,
						value: item.name,
						tempId: item.ID
					});
				}
				this.tempOptions = options;
			} else {
				this.toasted.error(result.msg);
			}
		},

		// 分块
		async getChunkList() {
			const params = {
				tempId: this.tempId
			};
			const result = await this.$store.dispatch(
				"GET_STAFF_WHITELIST_TEMP_CHUNK_LIST",
				params
			);

			if (result.code === 200) {
				const form = this.forms[this.searchFormId];

				const body = {
					proCode: form.proCode,
					tempCode: form.tempCode
				};

				const result2 = await this.$store.dispatch("GET_STAFF_WHITELIST_Top_Sum", body);

				console.log("rere", result2);

				if (result2.code === 200) {
					const headers = [cells.headersFirstChild];

					for (let item of result.data.list) {
						const width = item.name.length * 33 < 100 ? 100 : item.name.length * 35;

						headers.push({
							text: item.name,
							value: item.code,
							indeterminate: false,
							peopleSum: 0,
							width
						});
					}

					for (var i = 0; i < headers.length; i++) {
						for (const key of Object.keys(result2.data.list.blockSum)) {
							var arr = key.split("-");

							const peopleSum = result2.data.list.blockSum[key];

							if (headers[i].text == arr[1]) {
								headers[i].peopleSum = peopleSum;
							}
						}
					}

					this.headers = headers;

					this.getWhitelist();
				}
			} else {
				this.toasted.error(result.msg);
			}
		},

		// 分页
		handlePage({ pageNum, pageSize }) {
			this.params = {
				...this.params,
				pageIndex: pageNum
			};

			this.getWhitelist();
		},

		// 获取白名单列表
		async getWhitelist() {
			const form = this.forms[this.searchFormId];

			const body = {
				pageIndex: this.params.pageIndex,
				pageSize: this.params.pageSize,
				proCode: form.proCode,
				tempCode: form.tempCode,
				code: form.code
			};

			const result = await this.$store.dispatch("GET_STAFF_WHITELIST", body);

			const desserts = [];

			if (result.code === 200) {
				result.data.list.map((item, index) => {
					item.blockPermissions = item.blockPermissions || [];
					desserts[index] = {};
					this.headers.map((header, headerIndex) => {
						desserts[index][header.value] = {
							code: item.userCode,
							name: item.userName,
							blockCode: header.value,
							rowIndex: index,
							columnIndex: headerIndex,
							selected: item.blockPermissions.includes(header.value) ? true : false,
							indeterminate: false
						};
					});
				});
			}

			this.desserts = desserts;

			this.pagination.total = result.data.total;
		},

		// 选择项目
		selectPro(value) {
			this.tempOptions = [];
			this.desserts = [];

			const item = R.find(this.proOptions, value);
			this.proId = item.proId;

			this.createTempOptions();
		},

		// 选择模板
		selectTemp(value) {
			const item = R.find(this.tempOptions, value);
			this.tempId = item.tempId;
			this.getChunkList();
		},

		// header checkbox
		changeHeader(value, { columnIndex, headerCell }) {
			const dessertsForm = this.forms[this.tableDessertsFormId];
			const length = this.desserts.length;

			for (let i = 0; i < length; i++) {
				dessertsForm[`cell_${i}_${columnIndex}`] = value;
			}

			this.desserts.map(item => {
				item[headerCell.value].selected = value;
				// 不确定状态保持为：false
				item[headerCell.value].indeterminate = false;

				// 找到当前行
				const currentRow = [];

				for (let key in item) {
					if (!EXCLUDE.includes(key)) {
						currentRow.push(item[key]);
					}
				}

				// 设置 checkbox 状态
				const everyTrue = currentRow.every(e => e.selected);
				const someTrue = currentRow.some(s => s.selected);

				item["checkbox"].selected = everyTrue;

				item["checkbox"].indeterminate = everyTrue ? false : someTrue;
			});
		},

		// side checkbox
		changeSide(value, { rowIndex }) {
			const dessertsForm = this.forms[this.tableDessertsFormId];
			const length = this.headers.length;

			for (let i = 1; i < length; i++) {
				for (let i = 0; i < length; i++) {
					dessertsForm[`cell_${rowIndex}_${i}`] = value;
				}
			}

			for (let key in this.desserts[rowIndex]) {
				if (key !== "_X_ID") {
					this.desserts[rowIndex][key].selected = value;
					// 不确定状态保持为：false
					if (key === "checkbox") {
						this.desserts[rowIndex]["checkbox"].indeterminate = false;
					}
				}
			}

			// 找到当前列
			this.headers.map((item, index) => {
				if (!EXCLUDE.includes(item.value)) {
					const currentCol = [];
					this.desserts.map((dItem, dIndex) => {
						currentCol.push(dItem[item.value]);
					});

					// 设置 checkbox 状态
					const everyTrue = currentCol.every(e => e.selected);
					const someTrue = currentCol.some(s => s.selected);

					item.selected = everyTrue;

					item.indeterminate = everyTrue ? false : someTrue;
				}
			});
		},

		// desserts checkbox
		changeCell(value, { cell }) {
			this.desserts.map((item, index) => {
				if (index === cell.rowIndex) {
					item[cell.blockCode].selected = value;
				}
			});

			let [currentRow, currentCol] = [[], []];

			// 找到当前 cell 当前所在的行跟列
			this.desserts.map((item, index) => {
				// 当前行
				if (cell.rowIndex === index) {
					for (let key in item) {
						console.log(key);
						if (!EXCLUDE.includes(key)) {
							currentRow.push(item[key]);
						}
					}
				}

				// 当前列
				for (let key in item) {
					if (cell.blockCode === key) {
						if (!EXCLUDE.includes(key)) {
							currentCol.push(item[key]);
						}
					}
				}
			});

			// 找到当前 cell 对应的 header
			this.headers.map(item => {
				if (cell.blockCode === item.value) {
					const everyTrue = currentCol.every(e => e.selected);
					const someTrue = currentCol.some(s => s.selected);

					item.selected = everyTrue;

					item.indeterminate = everyTrue ? false : someTrue;
				}
			});

			// 找到当前 cell 对应的 checkbox
			this.desserts.map((item, index) => {
				if (cell.rowIndex === index) {
					const everyTrue = currentRow.every(e => e.selected);
					const someTrue = currentRow.some(s => s.selected);

					item["checkbox"].selected = everyTrue;

					item["checkbox"].indeterminate = everyTrue ? false : someTrue;
				}
			});
		},

		// 复制
		async onCopy() {
			const form = this.forms[this.searchFormId];
			console.log("form", form);
			const data = {
				proCode: form.proCode,
				tempName: form.tempCode,
				code: form.code,
				copyCode: form.copyCode.split(" ")
			};
			console.log("da", data);

			const result = await this.$store.dispatch("COPY_STAFF_WHITELIST_ITEM", data);

			this.toasted.dynamic(result.msg, result.code);
		},

		// 导出
		async onExport() {
			const form = this.forms[this.searchFormId];
			console.log("form", form);

			const result = await this.$store.dispatch("EXPORT_STAFF_WHITELIST", form);

			this.toasted.dynamic(`${result.msg}，正在导出...`, result.code);
		},

		// 保存
		async onSave() {
			const { proCode, tempCode } = this.forms[this.searchFormId];

			const staffs = [];

			this.desserts.map((item, index) => {
				staffs[index] = {
					userCode: item.checkbox.code,
					userName: item.checkbox.name,
					blockPermissions: []
				};

				for (let key in item) {
					if (!EXCLUDE.includes(key)) {
						if (item[key]["selected"]) {
							staffs[index].blockPermissions.push(item[key]["blockCode"]);
						}
					}
				}
			});

			const data = {
				proCode,
				tempName: tempCode,
				staffs
			};

			const result = await this.$store.dispatch("UPDATE_STAFF_WHITELIST_ITEM", data);

			this.toasted.dynamic(result.msg, result.code);
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
