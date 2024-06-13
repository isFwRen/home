<template>
	<div class="case">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row>
				<v-col :cols="12">
					<v-row class="z-flex align-end">
						<v-col
							v-for="(item, index) in cells.fields"
							:key="`case_filters_${index}`"
							:cols="item.cols"
						>
							<template v-if="item.inputType === 'input'">
								<z-text-field
									:formId="searchFormId"
									:formKey="item.formKey"
									:hideDetails="item.hideDetails"
									:hint="item.hint"
									:label="item.label"
									:suffix="item.suffix"
									:defaultValue="detail[item.formKey]"
									@change="onChange($event, item.formKey)"
								>
								</z-text-field>
							</template>

							<template v-else-if="item.inputType === 'date'">
								<z-date-picker
									:formId="searchFormId"
									:formKey="item.formKey"
									:clearable="item.clearable"
									:hideDetails="item.hideDetails"
									:hint="item.hint"
									:label="item.label"
									:options="item.options"
									range
									:suffix="item.suffix"
									z-index="10"
									:defaultValue="detail[item.formKey] || item.defaultValue"
									@input="onChange($event, item.formKey)"
								></z-date-picker>
							</template>

							<template v-else>
								<z-select
									:formId="searchFormId"
									:formKey="item.formKey"
									:hideDetails="item.hideDetails"
									:hint="item.hint"
									:label="item.label"
									:clearable="item.clearable"
									:options="cells[`${item.formKey}Options`] || []"
									:suffix="item.suffix"
									:defaultValue="detail[item.formKey]"
									@change="onChange($event, item.formKey)"
								></z-select>
							</template>
						</v-col>

						<div class="btns z-flex">
							<z-btn
								class="pr-3 pb-3"
								color="primary"
								:lockedTime="1000"
								@click="onSearch"
							>
								<v-icon class="text-h6">mdi-magnify</v-icon>
								查询
							</z-btn>

							<z-btn
								:formId="searchFormId"
								btnType="reset"
								class="pr-3 pb-3"
								color="error"
							>
								<v-icon class="text-h6">mdi-reload</v-icon>
								重置
							</z-btn>

							<z-btn class="pb-3" color="primary" :lockedTime="0" @click="onExpand">
								<template v-if="!expanded">
									<v-icon class="text-h6">mdi-chevron-down</v-icon>
									更多
								</template>

								<template v-else>
									<v-icon class="text-h6">mdi-chevron-up</v-icon>
									收起
								</template>
							</z-btn>
						</div>
					</v-row>
				</v-col>

				<v-col v-if="expanded" :cols="12">
					<v-row class="z-flex align-end">
						<v-col
							v-for="(item, index) in cells.moreFields"
							:key="`entry_filters_${index}`"
							:cols="item.cols"
						>
							<template v-if="item.inputType === 'input'">
								<z-text-field
									:formId="searchFormId"
									:formKey="item.formKey"
									:class="item.class"
									:hideDetails="item.hideDetails"
									:hint="item.hint"
									:label="item.label"
									:suffix="item.suffix"
									:defaultValue="item.defaultValue"
								>
									<span v-if="item.appendOuter" slot="prepend-outer">{{
										item.appendOuter
									}}</span>
								</z-text-field>
							</template>

							<template v-else>
								<z-select
									:formId="searchFormId"
									:formKey="item.formKey"
									:hideDetails="item.hideDetails"
									:hint="item.hint"
									:label="item.label"
									:clearable="item.clearable"
									:options="cells[`${item.formKey}Options`]"
									:suffix="item.suffix"
									:defaultValue="item.defaultValue"
								></z-select>
							</template>
						</v-col>
					</v-row>
				</v-col>
			</v-row>
		</div>

		<vxe-table
			:data="desserts"
			resizable
			:border="tableBorder"
			:cell-class-name="returnCellClassName"
			:expand-config="expandConfig"
			:max-height="tableMaxHeight"
			:size="tableSize"
			:sort-config="{
				multiple: true,
				defaultSort: { field: 'scanAt', order: 'asc' },
				trigger: 'cell'
			}"
			@sort-change="handleSort"
		>
			<template v-for="(item, idx) in cells.headers">
				<!-- 时间 BEGIN -->
				<vxe-column
					v-if="item.value === 'CreatedAt'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
					:fixed="item.fixed"
				>
					<template #default="{ row }">
						{{ row.CreatedAt | dateFormat("YYYY-MM-DD") }}
					</template>
				</vxe-column>
				<!-- 时间 END -->

				<!-- 扫描时间 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'scanAt'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ row.CreatedAt | dateFormat("HH:mm:ss") }}
						<!-- <div
							:title="row.scanAt | dateFormat('YYYY-MM-DD HH:mm:ss')"
							v-if="!isScanTime"
						>
							{{ row.scanAt | dateFormat("HH:mm:ss") }}
						</div>
						<span v-else> {{ row.scanAt | dateFormat("HH:mm:ss") }}</span> -->
					</template>
				</vxe-column>
				<!--扫描时间 END -->

				<!-- 案件号 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'billNum'"
					:field="item.value"
					:title="item.text"
					:fixed="item.fixed"
					type="expand"
					:key="idx * Math.random()"
					:width="item.width"
					sortable
				>
					<template #default="{ row }">
						<!-- <div v-if="row.wrongNote == '' || row.delRemarks == ''"  class="writeBlock"></div> -->
						<span
							:class="{
								'wrong-note': row.wrongNote,
								'pack-code': row.packCode
							}"
							@dblclick="onCase(row)"
						>
							{{ row.billNum }}</span
						>
					</template>

					<template #content="{ row }">
						<div
							class="pa-3 expand-wrapper"
							style="background: #f9f4b7"
							v-if="
								row.qualityUserCode !== '' ||
								row.stage !== 3 ||
								['B0103', 'B0106', 'B0110'].includes(row.proCode)
							"
						>
							<!-- <p v-if="row.status == 4" class="ma-0"> -->
							<p
								class="ma-0"
								:style="`width:${fixedWidth - 20}px;word-break: break-all;`"
							>
								{{ resonData(row) }}
							</p>
							<div class="ma-0" v-html="computedWrongNote(row)"></div>
						</div>
					</template>
				</vxe-column>
				<!-- 案件号 END -->

				<!--批次号 BEGIN-->
				<vxe-column
					v-else-if="item.value === 'batchNum'"
					:field="item.value"
					:fixed="item.fixed"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
				</vxe-column>
				<!-- 批次号 END -->

				<!-- 导出时间/回传时间 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'exportAt' || item.value === 'lastUploadAt'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						<div :title="row[item.value] | dateFormat('YYYY-MM-DD HH:mm:ss')">
							{{ row[item.value] | dateFormat("HH:mm:ss") }}
						</div>
					</template>
				</vxe-column>
				<!--导出时间/回传时间 END -->

				<!-- 案件状态 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'status'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						<div v-if="row.status == 4">
							{{ returnChineseName(row.status, cells.statusOptions) }}
						</div>
						<div v-else>
							{{ returnChineseName(row.status, cells.statusOptions) }}
						</div>
					</template>
				</vxe-column>
				<!-- 案件状态 END -->

				<!-- 录入状态 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'stage'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ returnChineseName(row.stage, cells.stageOptions) }}
					</template>
				</vxe-column>
				<!-- 录入状态 END -->

				<!-- 理赔类型 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'claimType'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ returnChineseName(row.claimType, cells.claimTypeOptions) }}
					</template>
				</vxe-column>
				<!-- 理赔类型 END -->

				<!-- 问题件 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'questionNum'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ row.questionNum }}
					</template>
				</vxe-column>
				<!-- 问题件 END -->

				<!-- 加急件 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'stickLevel'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ returnChineseName(row.stickLevel, cells.stickLevelOptions) }}
					</template>
				</vxe-column>
				<!-- 加急件 END -->

				<!-- 质检人 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'qualityUserCode'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:fixed="item.fixed"
					:width="item.width"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ row.qualityUserCode }}{{ row.qualityUserName }}
					</template>
				</vxe-column>
				<!-- 质检人 END -->

				<!-- 质检状态 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'exportStage'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:fixed="item.fixed"
					:width="item.width"
				>
					<template #default="{ row }">
						{{ billStage[row.exportStage] }}
					</template>
				</vxe-column>
				<!-- 质检状态 END -->

				<!-- 备注 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'remark'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:fixed="item.fixed"
					:width="item.width"
					show-overflow="ellipsis"
				>
					<template #default="{ row }">
						<v-tooltip
							top
							v-if="!row.remark"
							:disabled="row.stage === 7 ? true : false"
						>
							<template v-slot:activator="{ on, attrs }">
								<v-icon
									style="cursor: pointer"
									v-on="on"
									v-bind="attrs"
									@click="addRemark(row)"
									:color="row.stage === 5 ? 'grey' : 'primary'"
									>mdi-pen
								</v-icon>
							</template>
							<span>更新备注</span>
						</v-tooltip>

						<span v-else @click="addRemark(row)"> {{ row.remark }}</span>
					</template>
				</vxe-column>
				<!-- 备注 END -->

				<!-- 处理 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'dealWith'"
					:field="item.value"
					:fixed="item.fixed"
					:title="item.text"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<lp-dropdown
							color="primary"
							depressed
							offset-y
							smaller
							z-index="10"
							:options="rmOption(row, cells.moreOptions)"
							@click="onMore($event, row)"
						>
							更多
						</lp-dropdown>
					</template>
				</vxe-column>
				<!-- 处理 END -->

				<!-- 手动回传 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'isAutoUpload'"
					:field="item.value"
					:fixed="item.fixed"
					:title="item.text"
					:key="item.value"
					:width="item.width + 30"
				>
					<template #default="{ row }">
						<z-btn-toggle
							formId="upload"
							:formKey="row.ID"
							class="mt-n1"
							color="primary"
							dense
							:lockedTime="1000"
							mandatory
							:options="cells.autoUploadOptions"
							:defaultValue="row.isAutoUpload"
							@click="onEditAutoUpload($event, row)"
						></z-btn-toggle>
					</template>
				</vxe-column>
				<!-- 手动回传 END -->

				<!-- 操作 BEGIN -->
				<vxe-column
					v-else-if="item.value === 'options'"
					:field="item.value"
					:fixed="item.fixed"
					:title="item.text"
					:key="Math.random() + idx"
					:width="item.width + 60"
				>
					<template #default="{ row }">
						<div style="display: flex">
							<div>
								<v-tooltip top>
									<template v-slot:activator="{ on, attrs }">
										<v-icon
											class="mr-1"
											color="primary"
											v-bind="attrs"
											v-on="on"
											@click="onCheck(row)"
											>mdi-filter-check
										</v-icon>
									</template>
									<span> 质检 </span>
								</v-tooltip>
							</div>
							<div v-if="row.otherInfo != ''">
								<v-tooltip top>
									<template v-slot:activator="{ on, attrs }">
										<v-icon
											class="mr-1"
											color="primary"
											v-bind="attrs"
											v-on="on"
											@click="examine('/main/PM/case/update-result-xml', row)"
											>mdi-download
										</v-icon>
									</template>
									<span> 查看下载报文 </span>
								</v-tooltip>
							</div>
							<div
								v-if="row.stage !== 1 && row.stage !== 2 && row.status != 4"
								:data-row="row.stage"
							>
								<v-tooltip
									v-for="(record, index) in permisOptionsOptions"
									:key="Math.random() + index"
									top
								>
									<template v-slot:activator="{ on, attrs }">
										<v-icon
											:class="record.class"
											color="primary"
											v-bind="attrs"
											v-on="on"
											v-if="
												!((row.stage === 5 || row.stage === 7) &&
												record.value === 4
													? true
													: false)
											"
											@click="onOptions(record, row)"
											>{{ record.icon }}
										</v-icon>
									</template>
									<span> {{ record.label }}</span>
								</v-tooltip>
							</div>
						</div>
					</template>
				</vxe-column>
				<!-- 操作 END -->

				<vxe-column
					v-else
					:field="item.value"
					:fixed="item.fixed"
					:title="item.text"
					:key="item.value"
					:width="item.width"
					:sortable="item.sortable"
				>
				</vxe-column>
			</template>
		</vxe-table>

		<z-pagination
			class="mt-3"
			:total="pagination.total"
			:options="pageSizes"
			:pageNum="page.pageIndex"
			@page="handlePage"
		></z-pagination>

		<case-dialog ref="caseDialog" :rowInfo="detailInfo"></case-dialog>
		<check-dialog ref="checkDialog" :checkDetail="checkRow"></check-dialog>
		<remark-dialog
			ref="remarkDialog"
			@remarksClose="clearRemarks"
			:value="remarkValue"
			@remarksEmit="updateRemarks"
		>
		</remark-dialog>
		<!-- 删除备注 BEGIN -->
		<z-dynamic-form
			ref="delDynamic"
			:fieldList="cells.deleteFields"
			:width="550"
			@confirm="handleDelConfirm"
		></z-dynamic-form>
		<!-- 删除备注 END -->
		<alert-dialog ref="alertRef"> </alert-dialog>
		<!--悬浮时效简报--->
		<time-brief :proCode="timeProCode"></time-brief>
	</div>
</template>

<script>
import moment from "moment";
import { mapState, mapGetters } from "vuex";
import { tools, sessionStorage, localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import TableMixins from "@/mixins/TableMixins";
import CaseMixins from "./CaseMixins";
import cells from "./cells";

const scanProCode = [
	"B0103",
	"B0106",
	"B0108",
	"B0110",
	"B0114",
	"B0113",
	"B0118",
	"B0121",
	"B0122"
];

const { baseURLApi } = lpTools.baseURL();

const keys = new Map([
	["billStatus", "statusOptions"],
	["billInsuranceType", "insuranceTypeOptions"],
	["billClaimType", "claimTypeOptions"],
	["billStickLevel", "stickLevelOptions"],
	["billStage", "stageOptions"]
]);

const moreKeys = new Map([
	[2, "RELOAD_CASE_ITEM"],
	[3, "RECOVER_CASE_ITEM"],
	[4, "EXPORT_UNUSUAL_CASE_ITEM"],
	[5, "EXPORT_FORCE_CASE_ITEM"]
]);

const today = moment().format("YYYY-MM-DD");
const today1 = moment().format("YYYY-MM-DD");

export default {
	name: "Case",
	mixins: [TableMixins, CaseMixins],

	data() {
		return {
			permisOptionsOptions: [],
			permisButtons: [],
			fixedWidth: 0,
			addRow: {},
			billStage: [],
			formId: "Case",
			dispatchList: "GET_CASE_LIST",
			dispatchCellForm: "AUTO_UPLOAD_CASE_ITEM",
			cells,
			timeProCode: "",
			today,
			detail: {
				proCode: undefined,
				time: [today, today1]
			},
			images: [],
			remarkValue: "",
			editRow: {},
			expandConfig: {
				visibleMethod({ row }) {
					if (row.wrongNote === "" && row.remark === "" && row.delRemarks === "") {
						return false;
					}
					// if (row.wrongNote === "" && row.remark === "" && row.delRemarks !== "") {
					// 	return row.status === 4 ? true : false;
					// }

					return true;
				}
			},

			// 重加载开关
			flag: true,
			// 重加载返回结果
			result: "",
			// 所有案件重加载的状态
			allCaseStatus: {},

			checkRow: {}
		};
	},
	computed: {
		...mapState(["forms", "pm/case/case"]),
		...mapGetters(["auth", "project", "cases"]),
		returnChineseName() {
			return (value, items) => {
				const item = tools.find(items, {
					value: String(value)
				});
				return item?.label || "-";
			};
		},
		isScanTime() {
			this.project.code;
			return scanProCode.find(code => code === this.project.code) ? true : false;
		},
		compuWidth() {
			return {
				width: this.fixedWidth
			};
		}
	},
	watch: {
		"page.pageIndex": {
			handler(pageIndex) {
				this.updateCaseSearch("pageIndex", pageIndex);
			},
			immediate: true
		},
		"forms.CaseSearch": {
			handler(form) {
				this.timeProCode = form.proCode;
			}
		}
	},
	created() {
		this.setProOptions();
		this.setConstOptions();
		this.updateCaseSearch();
		this.$EventBus.$on("updateStage", row => {
			let obj = this.desserts.find(ele => row.ID === ele.ID);
		});

		this.detail = sessionStorage.get(this.searchFormId);

		this.filterButton();
	},
	beforeDestroy() {
		this.$store.commit("SET_PROJECT_INFO");
	},
	destroy() {
		this.$EventBus.$off("updateStage");
	},
	methods: {
		filterButton() {
			this.permisOptionsOptions = this.cells.optionsOptions;
			// const auths = localStorage.get("auth");
			// const menus = auths.menus;
			// const findNames = menus.find(item => item.name === "项目管理");
			// const nameChildren = findNames.children;
			// const findChildNames = nameChildren.find(item => item.name === "案件列表");
			// this.permisButtons = findChildNames.children;

			// this.cells.optionsOptions.map(option => {
			// 	const findOp = this.permisButtons.find(item => item.name === option.label);
			// 	if (findOp) {
			// 		this.permisOptionsOptions.push(option);
			// 		console.log(this.permisOptionsOptions, "eeeeeeeeeeeeeeeeee");
			// 	}
			// });
		},
		async onSearch() {
			this.params = {
				...this.params,
				...this.page
			};

			const total = await this.getList();
			if (typeof total !== "number") {
				return;
			}
			const index = this.page.pageIndex - 1;
			if (index * this.page.pageSize > total) {
				this.page.pageIndex = 1;
			}
		},
		async getList() {
			if (this.dispatchList) {
				const params = {
					...this.effectParams,
					...this.params,
					...this.forms[this.searchFormId]
				};
				const result = await this.$store.dispatch(this.dispatchList, params);
				const { list, total } = result.data;
				if (result.code === 200) {
					if (typeof list === "object") {
						if (list instanceof Array) {
							if (this.isScanTime) {
								// 扫描时间取值
								// list.forEach(item => {
								// 	item.scanAt = moment(item.CreatedAt).format("HH:mm:ss");
								// });
							}
							this.desserts = list;
							for (let key in this.desserts) {
								if (!this.allCaseStatus.hasOwnProperty([this.desserts[key].ID])) {
									this.allCaseStatus[this.desserts[key].ID] = true;
								}
							}
						} else {
							this.desserts = [];
						}
						this.pagination.total = total;
					} else {
						this.desserts = result.data;
						this.pagination.total = this.desserts.length;
					}
				} else {
					this.toasted.error(result.msg);

					this.desserts = [];
					this.pagination.total = 0;
				}

				this.sabayon = result;
				return total;
			}

			this.loading = false;

			return this.sabayon;
		},
		computedWrongNote(row) {
			const tableDom = document.querySelector(".vxe-table--fixed-left-wrapper");
			if (tableDom) {
				if (this.fixedWidth === 0) {
					this.fixedWidth = tableDom.offsetWidth;
				}
			}
			let remarkElement = row.remark
				? `<span style="display:block;width: ${this.fixedWidth - 20}px;">备注信息：${
						row.remark
				  }；</span>`
				: "";

			const listNode = row.wrongNote?.replace(/;|；/g, ";")?.split(";");

			let errorElement = "";
			if (listNode.length > 0) {
				listNode.forEach(element => {
					if (element !== "") {
						errorElement += `<span style="display:block;width: ${
							this.fixedWidth - 20
						}px;">错误内容：${element}；</span>`;
					}
				});
			}

			return remarkElement + errorElement;
		},

		returnCellClassName({ column, row }) {
			if (column.property === "billNum") {
				if (row.status === 4) {
					return "show-arrow";
				}

				if (!row.wrongNote && !row.packCode) {
					return "hide-arrow";
				}
			}
		},
		// Optionsbreak
		onOptions(record, row) {
			this.getDetail(row);

			switch (record.value) {
				case 1:
				case 3:
					this.$router.push({ path: record.path });
					this.$refs.caseDialog.onOpen();
					break;
				case 2:
					this.$store.commit("UPDATE_CASE", {
						rowInfo: this.detailInfo
					});
					this.judgeIsQuality(record.path);
					break;

				case 4:
					const noUploads = [
						"不可回传",
						"不允许回传",
						"不需要回传",
						"无需回传",
						"不用回传",
						"不需回传",
						"不能回传",
						"不要回传",
						"不回传",
						"不许回传"
					];
					const noUploadItem = noUploads.find(
						s => row.remark.includes(s) || row.wrongNote.includes(s)
					);
					if (noUploadItem) {
						this.$refs.alertRef.open();
					} else {
						this.upload();
					}
					break;
			}
		},
		// 查看下载报文
		examine(path, row) {
			this.$router.push({
				path,
				query: { row }
			});
			this.$refs.caseDialog.onOpen();
		},

		// 质检
		onCheck(row) {
			console.log(row);
			this.checkRow = row;
			let prefixUrl = `${baseURLApi}files/${row.downloadPath}`;
			let reg = new RegExp("/files/files/", "g");
			prefixUrl = prefixUrl.replace(reg, "/files/");
			const thumbs = [];

			if (!tools.isYummy(row.pictures)) {
				this.toasted.warning("没有图片!");
				return;
			}

			row.pictures.map(image => {
				if (row.proCode == "B0113") {
					thumbs.push({
						thumbPath: `${prefixUrl}${image.replace("/", "/A")}`,
						path: `${prefixUrl}${image}`
					});
				} else {
					thumbs.push({
						thumbPath: `${prefixUrl}A${image}`,
						path: `${prefixUrl}${image}`
					});
				}
			});

			sessionStorage.set("thumbs", thumbs);
			this.$refs.checkDialog.onOpen();
		},

		async updateRemarks(newRemark) {
			if (this.addRow.isFirstRemark) {
				//调用手动回传方法
				this.getDetail(this.addRow);
				this.modifyCell({
					proCode: this.detailInfo.proCode,
					id: this.detailInfo.ID,
					isAutoUpload: false
				});
			}

			// 提交到后台
			this.editRow.remark = newRemark;
			const { id, proCode, remark, editVersion } = this.editRow;

			const data = { id, proCode, remark, editVersion };

			const result = await this.$store.dispatch("UPDATE_REMARKS", data);
			if (result.code === 200) {
				this.toasted.dynamic(result.msg, result.code);
			}
		},
		clearRemarks(remark) {
			if (remark === "confirm") {
				let row = this.desserts.find(item => item.billNum === this.editRow.billNum);
				this.$set(row, "remark", this.editRow.remark);
				this.editRow = {};
				this.remarkValue = "";
			} else if (remark === "close") {
				//this.editRow = {}
				this.remarkValue = "";
			}
			this.addRow = {};
		},
		// 添加备注
		addRemark(row) {
			if (row.stage === 7) {
				this.toasted.dynamic("禁止修改已接收备注", 400);
				return;
			}
			// 第一次remark为空，回传状态改成“手动回传”
			if (!row.remark && row.isAutoUpload) {
				row.isFirstRemark = true;

				const rowsArr = Object.keys(row);
				rowsArr.forEach(key => {
					this.$set(this.addRow, key, row[key]);
				});
			}

			let { billNum, ID, proCode, editVersion } = row;
			this.editRow = {
				billNum,
				id: ID,
				proCode,
				editVersion
			};
			this.$refs.remarkDialog.onOpen();
			this.remarkValue = row.remark ? row.remark : "";
		},
		//处理删除原因和导出案件信息
		resonData(row) {
			let str = "";
			if (row.delRemarks) {
				str = `删除信息${row.delRemarks}`;
			}
			return str;
		},

		// 判断是否有人在质检
		async judgeIsQuality(path) {
			const body = {
				id: this.detailInfo.ID,
				proCode: this.detailInfo.proCode
			};

			const result = await this.$store.dispatch("CASE_JUDGE_IS_QUALITY", body);

			if (result.code === 200) {
				if (result.data) {
					this.$router.push({ path });
					this.$refs.caseDialog.onOpen();
				} else {
					this.$modal({
						visible: true,
						title: "注意",
						content: `该案件已有质检人员处理，是否继续进入？`,
						confirm: async () => {
							this.$router.push({ path });
							this.$refs.caseDialog.onOpen();
						}
					});
				}
			} else {
				this.toasted.warning(result.msg);
			}
		},
		// 回传
		upload() {
			this.$modal({
				visible: true,
				title: "回传提示",
				content: `是否回传案件：${this.detailInfo.billNum}`,
				confirm: async () => {
					const data = {
						proCode: this.detailInfo.proCode,
						id: this.detailInfo.ID
					};
					const result = await this.$store.dispatch("UPLOAD_CASE_ITEM", data);
					this.toasted.dynamic(result.msg, result.code);
					if (result.code === 200) {
						this.getList();
					}
				}
			});
		},
		// 更多
		onMore({ customValue: value }, row) {
			this.getDetail(row);

			const result = tools.find(cells.moreOptions, value);

			if (value === 1) {
				this.$refs.delDynamic.open({ title: "删除备注" });
			} else {
				this.$modal({
					visible: true,
					title: result.title,
					content: `${result.content}${row.billNum}`,
					confirm: async () => {
						if (!moreKeys.get(value)) return;

						const data = {
							proCode: this.detailInfo.proCode,
							id: this.detailInfo.ID
						};

						// 重加载提示
						if (value === 2) {
							if (this.allCaseStatus[row.ID]) {
								this.allCaseStatus[row.ID] = false;
								this.result = await this.$store.dispatch(moreKeys.get(value), data);

								if (this.result.code === 200) {
									this.toasted.dynamic(this.result.data, this.result.code);
									this.getList();
									this.allCaseStatus[row.ID] = true;
								}
							}
							if (!this.allCaseStatus[row.ID]) {
								this.toasted.warning("案件正在重加载中，请勿重复操作");
							}
						} else {
							const result = await this.$store.dispatch(moreKeys.get(value), data);
							this.toasted.dynamic(result.msg, result.code);

							if (result.code === 200) {
								this.getList();
							}
						}
					}
				});
			}
		},

		// 删除
		async handleDelConfirm({}, form) {
			const data = {
				proCode: this.detailInfo.proCode,
				id: this.detailInfo.ID,
				...form
			};

			const result = await this.$store.dispatch("DELETE_CASE_ITEM", data);
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
				this.$refs.delDynamic.close();
			}
		},

		// 手动回传
		//{ customValue: value }
		async onEditAutoUpload(value, row) {
			this.getDetail(row);
			this.$modal({
				visible: true,
				title: "手动回传提示",
				content: "请确认是否修改回传状态？",
				confirm: () => {
					this.modifyCell({
						proCode: this.detailInfo.proCode,
						id: this.detailInfo.ID,
						isAutoUpload: value
					});
				},
				cancel: () => {
					this.desserts = [];
					this.getList();
				}
			});
		},

		// 案件号
		onCase(row) {
			this.getDetail(row);
			if (row.imagesType && row.imagesType.length != 0) {
				let imagesType = row.imagesType.map(el => {
					// if (el.includes("zhenduan")) return "诊断书";
					// else if (el.includes("fapiao")) return "发票";
					// else if (el.includes("qingdan")) return "清单";
					// else if (el.includes("shenfenzheng")) return "身份证";
					// else if (el.includes("yinhangka")) return "银行卡";
					// else if (el.includes("jiesuan")) return "结算单";
					// else if (el.includes("shenqing")) return "申请书";
					// else if (el.includes("hukouben")) return "户口本";
					// else if (el.includes("chushengzhengming")) return "出生证明";
					if (el == "") return "其它";
					else return el;
				});
				sessionStorage.set("imagesType", imagesType);
			}

			let prefixUrl = `${baseURLApi}files/${row.downloadPath}`;
			let reg = new RegExp("/files/files/", "g");
			prefixUrl = prefixUrl.replace(reg, "/files/");
			const thumbs = [];

			if (!tools.isYummy(row.pictures)) {
				this.toasted.warning("没有图片!");
				return;
			}

			row.pictures.map(image => {
				if (row.proCode == "B0113") {
					thumbs.push({
						thumbPath: `${prefixUrl}${image.replace("/", "/A")}`,
						path: `${prefixUrl}${image}`
					});
				} else {
					thumbs.push({
						thumbPath: `${prefixUrl}A${image}`,
						path: `${prefixUrl}${image}`
					});
				}
			});

			sessionStorage.set("thumbs", thumbs);

			window.open(
				`${location.origin}/normal/view-images`,
				"_blank",
				"toolbar=yes, scrollbars=yes, resizable=yes"
			);
		},

		onChange(value, formKey) {
			this.updateCaseSearch(formKey, value);

			if (formKey === "proCode") {
				this.$store.commit("SET_PROJECT_INFO", {
					code: value
				});
				this.rememberCaseInfo({ proCode: value });
			}
		},

		// 设置项目下拉options
		async setProOptions() {
			this.cells.proCodeOptions = tools.deepClone(this.auth.proItems);

			const caseInfo = this.storage.get("caseInfo");

			let proCode =
				caseInfo && caseInfo.proCode
					? caseInfo.proCode
					: this.cells.proCodeOptions[0]?.value;

			this.$store.commit("SET_PROJECT_INFO", {
				code: proCode
			});
			this.rememberCaseInfo({ proCode });

			this.detail = {
				...this.detail,
				proCode
			};
			this.effectParams = {
				proCode: this.detail.proCode,
				time: this.detail.time
			};
		},

		// 设置除项目外的下拉options
		async setConstOptions() {
			const result = await this.$store.dispatch("GET_CASE_CONST_LIST");
			if (result.code === 200) {
				this.billStage = result.data.billStage;
				for (let key in result.data) {
					this.cells[keys.get(key)] = [];

					for (let _key in result.data[key]) {
						this.cells[keys.get(key)].push({
							label: result.data[key][_key],
							value: _key
						});
					}
				}
			}
		},

		// 过滤可操作性选项
		rmOption(row, options) {
			const cloneOptions = [];
			if (row.stage == 1) {
				return cloneOptions;
			}
			if (row.status != 4) {
				cloneOptions.push(options[0]);
			} else {
				cloneOptions.push(options[1]);

				cloneOptions.push(options[2]);

				cloneOptions.push(options[3]);
			}
			if (row.stage == 2) {
				cloneOptions.push(options[4]);
			}

			// const filterOptions = [];

			// cloneOptions.map(option => {
			// 	const findOp = this.permisButtons.find(item => item.name === option.label);
			// 	if (findOp) {
			// 		filterOptions.push(option);
			// 	}
			// });

			return cloneOptions;
		},

		updateCaseSearch(formKey, value) {
			if (!sessionStorage.get(this.searchFormId)) {
				sessionStorage.set(this.searchFormId, {});
			}

			const detail = sessionStorage.get(this.searchFormId);

			if (formKey) {
				sessionStorage.set(this.searchFormId, {
					...detail,
					[formKey]: value
				});
			}
		}
	},

	components: {
		"case-dialog": () => import("./caseDialog"),
		"check-dialog": () => import("./checkDialog"),
		"remark-dialog": () => import("./remarkDialog"),
		"alert-dialog": () => import("./alertDialog/"),
		"time-brief": () => import("./timeBrief/timeBrief.vue")
	}
};
</script>

<style scoped lang="scss">
// ::v-deep .vxe-table--main-wrapper {
// 	.expand-wrapper {
// 		background-color: #fff;
// 		border: 1px solid red;
// 	}
// }

.case {
	// ::v-deep .hide-arrow {
	//   .vxe-cell {
	//     .vxe-table--expanded,
	//     .writeBlock {
	//       display: none;
	//     }
	//   }
	// }

	span.wrong-note {
		background-color: #87d2f5;
		cursor: pointer;
	}

	span.pack-code {
		background-color: #7ee37d;
		cursor: pointer;
	}

	/* .writeBlock{
      position: absolute;
      height: 20px;
      width: 1em;
      background-color: #FFF;
      left:1em;
      top:calc(50% - 10px)
    } */
}
</style>
