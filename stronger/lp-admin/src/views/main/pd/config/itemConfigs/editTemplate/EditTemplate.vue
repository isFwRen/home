<template>
	<div class="edit-template">
		<div class="options">
			<div class="template">
				<v-row>
					<v-col cols="2">
						<z-select
							:formId="searchFormId"
							formKey="tempId"
							hideDetails
							label="模板"
							return-object
							:options="tempOptions"
							:defaultValue="tempId"
							@change="changeTemp"
						></z-select>
					</v-col>

					<v-col cols="2">
						<z-text-field
							:formId="searchFormId"
							formKey="blockName"
							hideDetails
							append-icon="mdi-magnify"
							label="分块名称"
							@enter="onSearch"
							@click:append="onSearch"
						></z-text-field>
					</v-col>

					<v-col class="z-flex align-end justify-end" cols="8">
						<z-btn
							v-for="item of cells.configs"
							:key="item.value"
							:class="item.class"
							color="primary"
							:disabled="item.disabled"
							outlined
							small
							@click="onTodo(item)"
						>
							{{ item.label }}</z-btn
						>
					</v-col>
				</v-row>
			</div>

			<div class="z-flex align-center pt-4 chunk">
				<label class="pt-4 mr-3 fw-bold">将分块</label>
				<z-text-field
					:formId="exchangeFormId"
					formKey="startOrder"
					class="mb-n6"
					:validation="[
						{ rule: 'required', message: '序号不能为空.' },
						{ rule: 'numeric', message: '序号只能为正整数.' }
					]"
					label="选中分块对应序号"
				>
				</z-text-field>
				<label class="pt-4 mx-3 fw-bold">插入到</label>
				<z-text-field
					:formId="exchangeFormId"
					formKey="endOrder"
					class="mb-n6"
					:validation="[
						{ rule: 'required', message: '序号不能为空.' },
						{ rule: 'numeric', message: '序号只能为正整数.' }
					]"
					label="指定分块对应序号"
				>
				</z-text-field>
				<z-btn
					:formId="exchangeFormId"
					btnType="validate"
					class="ml-3 mt-4"
					color="primary"
					@click="onExchange"
				>
					确定插入</z-btn
				>
				<v-btn
					btnType="validate"
					class="ml-3 mt-4"
					color="success"
					@click="onExportTemplate"
				>
					同步模板配置</v-btn
				>
				<v-spacer></v-spacer>

				<z-btn class="mt-4" color="primary" depressed fab rounded smaller @click="onNew">
					<v-icon>mdi-plus</v-icon>
				</z-btn>
			</div>

			<div class="mt-6 table edit-template-table">
				<vxe-table
					:data="desserts"
					:border="tableBorder"
					:max-height="tableMaxHeight"
					:size="tableSize"
					:stripe="tableStripe"
				>
					<vxe-column type="seq" title="序号" width="60"></vxe-column>

					<template v-for="item in cells.headers">
						<!-- 分块编码 BEGIN -->
						<vxe-column
							v-if="item.value === 'code'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn color="primary" outlined small @click="goToExchange(row)">
									{{ row.code }}
								</z-btn>
							</template>
						</vxe-column>
						<!-- 分块编码 END -->

						<!-- 分块名称 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'name'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-text-field
									formId="name"
									:formKey="row.code"
									class="mt-n4"
									hideDetails
									:defaultValue="row.name"
									@blur="onEditCell($event.customValue, item, row)"
								></z-text-field>
							</template>
						</vxe-column>
						<!-- 分块名称 END -->

						<!-- F8 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'fEight'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn-toggle
									formId="f8"
									:formKey="row.code"
									class="mt-n4"
									color="primary"
									dense
									mandatory
									:options="cells.choice"
									:defaultValue="row.fEight"
									@click="onEditCell($event, item, row)"
								></z-btn-toggle>
							</template>
						</vxe-column>
						<!-- F8 END -->

						<!-- OCR BEGIN -->
						<vxe-column
							v-else-if="item.value === 'ocr'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn-toggle
									formId="ocr"
									:formKey="row.code"
									class="mt-n4"
									color="primary"
									dense
									mandatory
									:options="cells.codes"
									:defaultValue="row.ocr"
									@click="onEditCell($event, item, row)"
								></z-btn-toggle>
							</template>
						</vxe-column>
						<!-- OCR END -->

						<!-- 释放时间(秒) BEGIN -->
						<vxe-column
							v-else-if="item.value === 'freeTime'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-text-field
									formId="time"
									:formKey="row.code"
									class="mt-n4"
									hideDetails
									:disabled="tempItem?.label === cells.UNDEFINED"
									:defaultValue="row.freeTime"
									@blur="onEditCell($event.customValue, item, row)"
								></z-text-field>
							</template>
						</vxe-column>
						<!-- 释放时间(秒) END -->

						<!-- 关联 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'relation'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-text-field
									formId="relation"
									:formKey="row.code"
									class="mt-n4"
									hideDetails
									:defaultValue="row.relation"
									@blur="onEditCell($event.customValue, item, row)"
								></z-text-field>
							</template>
						</vxe-column>
						<!-- 关联 END -->

						<!-- 截图配置 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'screenshot'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn
									color="primary"
									depressed
									fab
									rounded
									smaller
									@click="openScreenshotDialog(item, row)"
								>
									<v-icon>mdi-monitor-screenshot</v-icon>
								</z-btn>
							</template>
						</vxe-column>
						<!-- 截图配置 END -->

						<!-- 流程 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'isCompetitive'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn-toggle
									formId="competitive"
									:formKey="row.code"
									class="mt-n4"
									color="primary"
									dense
									mandatory
									:options="cells.competitive"
									:defaultValue="row.isCompetitive"
									@click="onEditCell($event, item, row)"
								></z-btn-toggle>
							</template>
						</vxe-column>
						<!-- 流程 END -->

						<!-- 循环分块 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'isLoop'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn-toggle
									formId="loop"
									:formKey="row.code"
									class="mt-n4"
									color="primary"
									dense
									mandatory
									:options="cells.choice"
									:defaultValue="row.isLoop"
									@click="onEditCell($event, item, row)"
								></z-btn-toggle>
							</template>
						</vxe-column>
						<!-- 循环分块 END -->

						<!-- 手机录入 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'isMobile'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:min-width="item.minWidth"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn-toggle
									formId="mobileEntry"
									:formKey="row.code"
									class="mt-n4"
									color="primary"
									dense
									mandatory
									:options="cells.choice"
									:defaultValue="row.isMobile"
									@click="onEditCell($event, item, row)"
								></z-btn-toggle>
							</template>
						</vxe-column>
						<!-- 手机录入 END -->

						<!-- 手机截图 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'mobileScreenshot'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn
									color="primary"
									depressed
									fab
									rounded
									smaller
									@click="openScreenshotDialog(item, row)"
								>
									<v-icon>mdi-cellphone-screenshot</v-icon>
								</z-btn>
							</template>
						</vxe-column>
						<!-- 手机截图 END -->

						<!-- 删除 BEGIN -->
						<vxe-column
							v-else-if="item.value === 'delete'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-btn
									color="error"
									depressed
									fab
									rounded
									smaller
									@click="onDelete(row)"
								>
									<v-icon>mdi-trash-can-outline</v-icon>
								</z-btn>
							</template>
						</vxe-column>
						<!-- 删除 END -->

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
			</div>
		</div>

		<!-- 复制模板 BEGIN -->
		<copy-temp-dialog
			ref="copyTemp"
			tempType="copy"
			:tempItem="tempItem"
			@submitted="updateTempPage"
		></copy-temp-dialog>
		<!-- 复制模板 END -->

		<!-- 修改模板图片 BEGIN -->
		<change-temp-img-dialog
			ref="changeTemplateImg"
			:imageList="imageList"
			@uploaded="handleUploaded"
		></change-temp-img-dialog>
		<!-- 修改模板图片 END -->

		<!-- 截图配置/手机截图 BEGIN -->
		<desensitization-dialog
			ref="desensitization"
			:rowInfo="detailInfo"
			:imageOptions="tempImageOptions"
			:imageList="imageList"
			@dialog="handleDialog"
			@updated="getList"
		></desensitization-dialog>
		<!-- 截图配置/手机截图 END -->

		<!-- 添加新分块 BEGIN -->
		<new-chunk-dialog
			ref="newChunk"
			:tempItem="tempItem"
			:rowInfo="detailInfo"
			@submitted="getList"
		></new-chunk-dialog>
		<!-- 添加新分块 END -->

		<!-- 字段配置 BEGIN -->
		<exchange-field-chunk ref="exchange"></exchange-field-chunk>
		<!-- 字段配置 END -->
	</div>
</template>

<script>
import storage from "@/libs/util.storage";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "../../ConfigMixins";
import { tools } from "vue-rocket";
import cells from "./cells";
import { mapGetters } from "vuex";
export default {
	name: "EditTemplateDialog",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "EditTemplateDialog",
			exchangeFormId: "EditTemplateExchange",
			dispatchList: "GET_CONFIG_TEMP_CHUNK_LIST",
			dispatchCellForm: "UPDATE_CONFIG_TEMP_CHUNK",
			dispatchDelete: "DEL_CONFIG_TEMP_CHUNK",
			cells,

			tempId: storage.get("config").tempId,
			tempOptions: [],
			tempItem: {},
			effectParams: {
				tempId: storage.get("config").tempId
			},
			tempImageOptions: [],
			imageList: []
		};
	},

	created() {
		this.setTempOptions();
	},
	computed: {
		...mapGetters(["project"])
	},

	methods: {
		onTodo(item) {
			this[item.fn]();
		},
		async onExportTemplate() {
			const body = {
				proCode: this.project.code,
				mtype: "template",
				templateId: this.tempItem?.label
			};

			const result = await this.$store.dispatch("EXPORT_OR_IMPORT_EXPORT", body);
			if (result.code === 200) {
				this.toasted.dynamic(result.msg, result.code);
			}
		},
		onImportTemplate() {},
		goToExchange(row) {
			this.rememberIds({ chunkId: row.ID });
			this.$store.commit("UPDATE_EDIT_TEMP", { chunkId: row.ID });
			this.$refs.exchange.onOpen("字段配置");
		},

		// 复制模板
		copyTemplate() {
			this.$refs.copyTemp.onOpen("复制模板");
		},

		// 修改模板图片
		modifyTemplateImg() {
			this.$refs.changeTemplateImg.onOpen();
		},

		handleDialog(dialog) {
			console.log(dialog);
		},

		// 脱敏配置
		desensitization() {
			this.$refs.desensitization.onOpen();
		},

		// 添加新分块
		onNew() {
			const row = {
				myOrder: this.sabayon.data.maxCode || 0
			};
			this.getDetail(row);
			this.$refs.newChunk.onOpen(-1);
		},

		// 修改单个单元格
		onEditCell(value, item, row) {
			if (typeof value === "object") {
				value = value.customValue;
			}

			this.modifyCell({
				tempId: this.tempId,
				id: row.ID,
				...row,
				[item.value]: value
			});
		},

		// 截图配置
		openScreenshotDialog(item, row) {
			if (tools.isLousy(this.tempImageOptions)) {
				this.toasted.warning(`当前分块暂无${item.text}！`);
			} else {
				this.getDetail({
					...row,
					coordinateType: item.coordinateType
				});
				this.$refs.desensitization.onOpen(item.text);
			}
		},

		// 删除
		onDelete(row) {
			this.getDetail(row);

			this.$modal({
				visible: true,
				title: "删除提示",
				content: "请确认是否要删除？",
				confirm: () => {
					this.delRows();
				}
			});
		},

		// 交换位置
		async onExchange() {
			const { startOrder, endOrder } = this.forms[this.exchangeFormId];

			const startItem = this.desserts.find(d => d.myOrder === +startOrder) || {};
			const endItem = this.desserts.find(d => d.myOrder === +endOrder) || {};

			const form = {
				startId: startItem.ID,
				startOrder,
				endId: endItem.ID,
				endOrder
			};

			const result = await this.$store.dispatch("EXCHANGE_CONFIG_TEMP_CHUNK", form);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		// 模板下拉选项
		async setTempOptions() {
			const result = await this.$store.dispatch("GET_CONFIG_TEMP_OPTIONS", {
				proId: this.proId
			});
			const list = [];

			if (result.code === 200) {
				result.data.map(item => {
					list.push({
						label: item.name,
						value: item.ID,
						imageList: item.images || []
					});
				});
			}

			this.tempOptions = list;
			this.tempItem = tools.find(this.tempOptions, { value: this.tempId });

			this.setImageOptions();

			return result;
		},

		// 切换模板选项
		changeTemp(temp) {
			this.tempItem = temp;
			this.tempId = this.tempItem.value;
			this.effectParams = { tempId: this.tempItem.value };

			this.setImageOptions();
			this.rememberIds({ tempId: this.tempItem.value });
			this.getList();
		},

		// 设置图片列表
		setImageOptions() {
			this.imageList = this.tempItem?.imageList;

			this.tempImageOptions = [];

			this.imageList?.map((image, index) => {
				const splitImage = image.split("/");
				const lastIndex = splitImage.length - 1;

				this.tempImageOptions.push({
					value: index,
					label: splitImage[lastIndex]
				});
			});
		},

		// 复制模板后更新【修改模板】页面
		async updateTempPage() {
			const result = await this.setTempOptions();

			if (result.code === 200) {
				this.changeTemp(this.tempItem);
			}
		},

		// 修改模板图片后更新模板列表
		handleUploaded() {
			this.setTempOptions();
		}
	},

	watch: {
		"sabayon.data.list": {
			handler(list) {
				if (tools.isYummy(list)) {
					this.$store.commit("UPDATE_EDIT_TEMP", { chunkList: list });
				}
			},
			immediate: true
		}
	},

	components: {
		"copy-temp-dialog": () => import("../../newTemplateDialog"),
		"desensitization-dialog": () => import("./desensitizationDialog"),
		"change-temp-img-dialog": () => import("./changeTempImgDialog"),
		"new-chunk-dialog": () => import("./newChunkDialog"),
		"exchange-field-chunk": () => import("./exchangeFieldChunk")
	}
};
</script>
