<template>
	<div class="config-field">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="name"
						hideDetails
						label="字段名称"
					></z-text-field>
				</v-col>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="inputProcess"
						clearable
						hideDetails
						label="录入工序"
						:options="cells.fields[4].options"
					></z-select>
				</v-col>

				<div class="z-flex pb-3 btns">
					<z-btn class="pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn
						:formId="searchFormId"
						btnType="reset"
						class="pl-3"
						color="warning"
						@click="onSearch"
					>
						<v-icon class="text-h6">mdi-reload</v-icon>
						重置
					</z-btn>
					<v-btn btnType="validate" class="ml-3" color="success" @click="onExportField">
						同步字段配置</v-btn
					>
				</div>
			</v-row>

			<z-btn color="primary" @click="onNew">
				<v-icon class="text-h6">mdi-plus</v-icon>
				新增
			</z-btn>
		</div>

		<div class="table config-field-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:height="400"
				:max-height="tableMaxHeight"
				:loading="loading"
				:size="tableSize"
			>
				<vxe-column title="序号" width="60">
					<template #default="item">
						{{ increaseSeq(item) }}
					</template>
				</vxe-column>

				<template v-for="item in cells.headers">
					<!-- 录入工序 BEGIN -->
					<vxe-column
						v-if="item.value === 'inputProcess'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="mt-n2 mb-n1">
								<z-btn-toggle
									formId="process"
									:formKey="`inputProcess${row.ID}`"
									color="primary"
									mandatory
									:options="cells.fields[4].options"
									:defaultValue="row.inputProcess"
									@change="changeInputProcess($event, row)"
								>
								</z-btn-toggle>
							</div>
						</template>
					</vxe-column>
					<!-- 录入工序 END -->

					<!-- 问题件配置 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'opqConfig'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<v-icon @click="opqConfig(row)">mdi-help-circle-outline</v-icon>
						</template>
					</vxe-column>
					<!-- 问题件配置 END -->

					<!-- 导出校验配置 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'exportConfig'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<v-icon @click="exportConfig(row)">mdi-alert-circle-outline</v-icon>
						</template>
					</vxe-column>
					<!-- 导出校验配置 END -->

					<!--字段拦截提示配置BEGIN -->
					<vxe-column
						v-else-if="item.value === 'interceptionConfig'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<v-icon @click="interceptionConfig(row)">mdi-alert</v-icon>
						</template>
					</vxe-column>

					<!--字段拦截提示配置 END -->

					<vxe-column
						v-else-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								<z-btn color="primary" depressed small @click="onEdit(row)">
									编辑
								</z-btn>

								<lp-dropdown
									class="pl-3"
									color="primary"
									depressed
									offset-y
									small
									:options="cells.moreOptions"
									@click="onMore($event, row)"
								>
									更多
									<v-icon>mdi-chevron-down</v-icon>
								</lp-dropdown>
							</div>
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
				:pageNum="page.pageIndex"
				:total="pagination.total"
				:options="pageSizes"
				@page="handlePage"
			></z-pagination>
		</div>

		<!-- 新增/编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:config="formConfig"
			:detail="detailInfo"
			:fieldList="cells.fields"
			:width="900"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->

		<opq-config-dialog ref="opqConfig" :Row="fRow" :issueList="fissueList"></opq-config-dialog>

		<!-- 导出校验配置 BEGIN -->
		<v-dialog v-model="dialog" transition="dialog-top-transition" max-width="1200">
			<template v-slot:default="dialog">
				<v-card>
					<v-toolbar color="primary" dark>导出校验配置</v-toolbar>
					<v-card-text>
						<p style="margin: 12px 0; font-size: 15px">
							字段：{{ exportRow.code + "_" + exportRow.name }}
						</p>
						<v-row cols="12" style="font-weight: 800; font-size: 16px">
							<v-col cols="2"> 录入值执行条件 </v-col>
							<v-col cols="3"> 条件内容 </v-col>
							<v-col cols="6"> 导出校验描述 </v-col>
							<v-col cols="1" style="display: flex; justify-content: flex-end">
								<v-icon
									color="blue darken-2"
									style="cursor: pointer"
									@click="addExport()"
								>
									mdi-view-grid-plus
								</v-icon>
							</v-col>
						</v-row>

						<v-row
							cols="12"
							v-for="(el, index) in cells.exportObject"
							:key="index"
							style="border: 1px solid #dcdddc; height: 65px; margin-bottom: 5px"
						>
							<v-col cols="2">
								<v-select
									v-model="el.checkType"
									:items="el.items"
									outlined
									dense
								></v-select>
							</v-col>
							<v-col cols="3">
								<v-text-field
									v-model="el.value"
									outlined
									clearable
									dense
								></v-text-field>
							</v-col>
							<v-col cols="6">
								<v-text-field
									v-model="el.mark"
									outlined
									clearable
									dense
								></v-text-field>
							</v-col>
							<v-col cols="1" style="display: flex; justify-content: flex-end">
								<v-icon
									color="red darken-2"
									large
									style="cursor: pointer; margin-bottom: 30px; margin-right: 20px"
									@click="deleteExport(index)"
								>
									{{ el.icon }}
								</v-icon>
							</v-col>
						</v-row>
					</v-card-text>
					<v-card-actions class="justify-end">
						<v-btn @click="dialog.value = false">取消</v-btn>
						<v-btn color="primary" @click="ExportConfirm()">确认</v-btn>
					</v-card-actions>
				</v-card>
			</template>
		</v-dialog>
		<!-- 导出校验配置 END -->

		<!-- 字段拦截提示 BEGIN -->
		<v-dialog v-model="interceptionDialog" transition="dialog-top-transition" max-width="1200">
			<template v-slot:default="dialog">
				<v-card>
					<v-toolbar color="primary" dark>字段拦截提示配置</v-toolbar>
					<v-card-text>
						<p style="margin: 12px 0; font-size: 15px">
							字段：{{ interceptionRow.code + "_" + interceptionRow.name }}
						</p>
						<v-row cols="12" style="font-weight: 800; font-size: 16px">
							<v-col cols="2"> 录入值执行条件 </v-col>
							<v-col cols="3"> 条件内容 </v-col>
							<v-col cols="2"> 是否拦截 </v-col>
							<v-col cols="4"> 字段拦截提示 </v-col>
							<v-col cols="1" style="display: flex; justify-content: flex-end">
								<v-icon
									color="blue darken-2"
									style="cursor: pointer"
									@click="addExport()"
								>
									mdi-view-grid-plus
								</v-icon>
							</v-col>
						</v-row>

						<v-row
							cols="12"
							v-for="(el, index) in fields"
							:key="index + el.item"
							style="border: 1px solid #dcdddc; height: 65px; margin-bottom: 5px"
						>
							<v-col cols="2">
								<v-select
									v-model="el.item"
									:items="el.items"
									outlined
									dense
								></v-select>
							</v-col>
							<v-col cols="3">
								<v-text-field
									v-model="el.itemContent"
									outlined
									clearable
									dense
								></v-text-field>
							</v-col>
							<v-col cols="2">
								<div class="mt-n2 mb-n1">
									<v-btn-toggle color="blue" v-model="toggle_exclusive">
										<v-btn v-for="item in el.buttongroups" :key="item.value">{{
											item.label
										}}</v-btn>
									</v-btn-toggle>
								</div>
							</v-col>
							<v-col cols="4">
								<v-text-field
									v-model="el.itemDesc"
									outlined
									clearable
									dense
								></v-text-field>
							</v-col>
							<v-col cols="1" style="display: flex; justify-content: flex-end">
								<v-icon
									color="red darken-2"
									large
									style="cursor: pointer; margin-bottom: 30px; margin-right: 20px"
									@click="deleteExport(index)"
								>
									{{ el.icon }}
								</v-icon>
							</v-col>
						</v-row>
					</v-card-text>
					<v-card-actions class="justify-end">
						<v-btn @click="dialog.value = false">取消</v-btn>
						<v-btn color="primary" @click="ExportConfirm()">确认</v-btn>
					</v-card-actions>
				</v-card>
			</template>
		</v-dialog>
		<!-- 字段拦截提示 END -->
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";

const getCode = new Map([
	["1", "等于"],
	["2", "包含"],
	["3", "不包含"]
]);
const reverseGetCode = new Map([
	["等于", "1"],
	["包含", "2"],
	["不包含", "3"]
]);

export default {
	name: "ConfigField",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "configField",
			fRow: "",
			exportRow: "",
			fissueList: "",
			dispatchList: "GET_CONFIG_FIELD_LIST",
			dispatchDelete: "DELETE_CONFIG_FIELD_ITEM",
			cells,
			formConfig: {
				inputProcess: {
					items: []
				},

				checkDate: {
					items: []
				},

				validations: {
					items: []
				}
			},
			// 导出校验配置
			dialog: false,
			fields: "",
			interceptionDialog: false,
			toggle_exclusive: 0
		};
	},

	computed: {
		...mapGetters(["config"])
	},

	created() {
		this.getConst();
	},

	methods: {
		async onExportField() {
			const body = {
				proCode: this.config.pro.code,
				mtype: "field",
				templateId: ''
			};

			const result = await this.$store.dispatch("EXPORT_OR_IMPORT_EXPORT", body);
			if(result.code === 200){
				this.toasted.dynamic(result.msg, result.code);
			}
		},
		handleInterceptionChange(value) {
			console.log(this.form["inputProcessisInterception"]);
			console.log(value, "value");
		},
		onNew() {
			const row = {
				code: this.sabayon.data.maxCode,
				myOrder: this.sabayon.data.total || 1
			};
			this.getDetail(row);

			this.detailInfo["validations"] = [];
			this.$refs.dynamic.open({ title: "新增", status: -1 });
		},

		async opqConfig(row) {
			this.fRow = row;
			const result = await this.$store.dispatch("GET_ISSUE_CONFIG_LIST", { fId: row.ID });
			this.fissueList = result.data;
			this.$refs.opqConfig.fissueList = result.data;
			this.$refs.opqConfig.onOpen();
		},

		// 导出配置弹框
		async exportConfig(row) {
			this.exportRow = row;
			console.log(row);
			const result = await this.$store.dispatch("GET_EXPORT_CONFIG_LIST", { fId: row.ID });
			console.log(result.data);
			let arr = result.data.list;
			arr.forEach(el => {
				el.checkType = getCode.get(el.checkType);
				el.items = ["等于", "包含", "不包含"];
				el.icon = "mdi-delete-empty";
			});
			this.cells.exportObject = arr;
			this.dialog = true;
		},

		// 字段拦截提示
		interceptionConfig(row) {
			this.interceptionRow = row;
			this.fields = this.cells.interceptionObject;
			this.interceptionDialog = true;
		},

		// 导出配置确认按钮
		ExportConfirm() {
			let arr = this.cells.exportObject;
			arr.forEach(el => {
				el.checkType = reverseGetCode.get(el.checkType);
				el.code = this.exportRow.code;
				el.proId = this.exportRow.proId;
				el.fId = this.exportRow.ID;
				delete el.icon;
				delete el.items;
			});
			this.$store.dispatch("ADD_EXPORT_CONFIG", { fId: this.exportRow.ID, list: arr });
			this.dialog = false;
		},

		// 导出配置添加按钮
		addExport() {
			this.cells.exportObject.push({
				items: ["等于", "包含", "不包含"],
				checkType: "等于",
				value: "",
				mark: "",
				icon: "mdi-delete-empty"
			});
		},

		// 导出配置删除按钮
		deleteExport(index) {
			this.cells.exportObject.splice(index, 1);
		},

		onEdit(row) {
			row = {
				...row,
				validations: row.validations || []
			};
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		handleConfirm(effect, form) {
			form = {
				myOrder: this.detailInfo.myOrder,
				id: this.detailInfo.ID,
				proId: this.config.proId,
				status: effect.status,
				...form
			};

			this.updateListItem(form, "UPDATE_CONFIG_FIELD");
		},

		async changeInputProcess(value, row) {
			this.getDetail(row);

			const body = {
				...row,
				id: row.ID,
				inputProcess: value
			};

			const result = await this.$store.dispatch("UPDATE_CONFIG_FIELD_INPUT_PROCESS", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		onMore({ customValue }, row) {
			if (customValue === "copy") {
				const _row = {
					code: this.sabayon.data.maxCode,
					myOrder: this.sabayon.data.total,
					name: `${row.name}_${this.sabayon.data.total}`
				};

				row = {
					...row,
					..._row
				};

				this.getDetail(row);

				this.$modal({
					visible: true,
					title: "复制提示",
					content: "请确认是否要复制当前字段？",
					confirm: () => {
						this.onCopy();
					}
				});
			} else {
				this.getDetail(row);
				this.deleteItem();
			}
		},

		// 复制
		async onCopy() {
			const body = {
				status: -1,
				...this.detailInfo
			};

			const result = await this.$store.dispatch("UPDATE_CONFIG_FIELD", body);

			if (result.code === 200) {
				this.toasted.success("复制成功!");
				this.getList();
			} else {
				this.toasted.error(result.msg);
			}
		},

		// 常量
		async getConst() {
			const result = await this.$store.dispatch("GET_CONFIG_CONSTANT");

			if (result.code !== 200) return;

			const { inputProcess, dateLimit, validation } = result.data;

			const [inputProcessItems, dateLimitItems, validationItems] = [[], [], []];

			for (let key in inputProcess) {
				inputProcessItems.push({
					label: inputProcess[key],
					value: +key
				});
			}

			for (let key in dateLimit) {
				dateLimitItems.push({
					label: dateLimit[key],
					value: +key
				});
			}

			for (let key in validation) {
				validationItems.push({
					label: validation[key],
					value: +key
				});
			}

			this.formConfig = {
				inputProcess: {
					items: inputProcessItems
				},

				checkDate: {
					items: dateLimitItems
				},

				validations: {
					items: validationItems
				}
			};

			// this.cells.fields[4].options = []
			// this.cells.fields[5].options = []
			// this.cells.fields[17].options = []

			// // 录入工序
			// for(let key in inputProcess) {
			//   this.cells.fields[4].options.push({
			//     label: inputProcess[key],
			//     value: +key
			//   })
			// }

			// // 录入工序
			// for(let key in dateLimit) {
			//   this.cells.fields[5].options.push({
			//     label: dateLimit[key],
			//     value: +key
			//   })
			// }

			// // validations
			// for(let key in validation) {
			//   this.cells.fields[17].options.push({
			//     label: validation[key],
			//     value: +key
			//   })
			// }
		}
	},

	components: {
		"opq-config-dialog": () => import("./opqConfigDialog")
	}
};
</script>
