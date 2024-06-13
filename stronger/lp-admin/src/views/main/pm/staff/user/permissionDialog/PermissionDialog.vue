<template>
	<lp-dialog
		ref="dialog"
		:title="`${detail.nickName}的项目权限`"
		:width="1080"
		transition="dialog-bottom-transition"
		@dialog="handleDialog"
	>
		<div class="pt-6" slot="main">
			<div class="table">
				<vxe-table
					:data="desserts"
					:border="tableBorder"
					:max-height="tableMaxHeight"
					:loading="loading"
					:size="tableSize"
				>
					<template v-for="item in cells.headers">
						<!-- 权限设置 BEGIN -->
						<vxe-colgroup
							v-if="item.value === 'permission'"
							:key="item.value"
							:title="item.text"
						>
							<vxe-column :field="item.value" :width="item.width">
								<template #header="{ _columnIndex }">
									<v-row class="pa-0 ma-0" v-if="dialog">
										<v-col
											v-for="item in permissionHeadersCheckboxs"
											class="pa-0 ma-0"
											:cols="2"
											:key="item.formKey"
										>
											<z-checkbox
												:formId="tableHeaderFormId"
												:formKey="item.formKey"
												:class="item.class"
												:dense="item.dense"
												:hide-details="item.hideDetails"
												:indeterminate="item.indeterminate"
												:defaultValue="item.selected"
												@change="
													changeHeader($event, { ...item, _columnIndex })
												"
											></z-checkbox>
										</v-col>
									</v-row>
								</template>

								<template #default="{ _columnIndex, rowIndex, row }">
									<v-row class="pa-0 ma-0">
										<v-col
											v-for="(
												item, index
											) in cells.permissionDessertsCheckboxs"
											class="pa-0 ma-0"
											:cols="2"
											:key="`${item.initFormKey}_${rowIndex}_${_columnIndex}_${index}`"
										>
											<z-checkbox
												:formId="tableDessertsFormId"
												:formKey="`${item.initFormKey}_${rowIndex}_${_columnIndex}_${index}`"
												:label="item.label"
												:class="item.class"
												:dense="item.dense"
												:hide-details="item.hideDetails"
												:defaultValue="row[item.value]"
												@change="
													changeCell($event, {
														...item,
														_columnIndex,
														rowIndex
													})
												"
											></z-checkbox>
										</v-col>
									</v-row>
								</template>
							</vxe-column>
						</vxe-colgroup>
						<!-- 权限设置 END -->

						<!-- 内外网 BEGIN -->
						<vxe-colgroup
							v-else-if="item.value === 'network'"
							:key="item.value"
							:title="item.text"
						>
							<vxe-column :field="item.value" :width="item.width">
								<template #header="{ _columnIndex }">
									<v-row class="pa-0 ma-0" v-if="dialog">
										<v-col
											v-for="item in netHeadersCheckboxs"
											class="pa-0 ma-0"
											:cols="6"
											:key="item.formKey"
										>
											<z-checkbox
												:formId="tableHeaderFormId"
												:formKey="item.formKey"
												:class="item.class"
												:dense="item.dense"
												:hide-details="item.hideDetails"
												:indeterminate="item.indeterminate"
												:defaultValue="item.selected"
												@change="
													changeHeader($event, { ...item, _columnIndex })
												"
											></z-checkbox>
										</v-col>
									</v-row>
								</template>

								<template #default="{ rowIndex, _columnIndex, row }">
									<v-row class="pa-0 ma-0">
										<v-col
											v-for="(item, index) in cells.netDessertsCheckboxs"
											class="pa-0 ma-0"
											:cols="6"
											:key="`${item.initFormKey}_${rowIndex}_${_columnIndex}_${index}`"
										>
											<z-checkbox
												:formId="tableDessertsFormId"
												:formKey="`${item.initFormKey}_${rowIndex}_${_columnIndex}_${index}`"
												:label="item.label"
												:class="item.class"
												:dense="item.dense"
												:hide-details="item.hideDetails"
												:defaultValue="row[item.value]"
												@change="
													changeCell($event, {
														...item,
														_columnIndex,
														rowIndex
													})
												"
											></z-checkbox>
										</v-col>
									</v-row>
								</template>
							</vxe-column>
						</vxe-colgroup>
						<!-- 内外网 END -->

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
		</div>

		<div class="z-flex" slot="actions">
			<z-btn class="mr-3" color="normal" @click="onClose">取消</z-btn>

			<z-btn class="mr-2" color="primary" @click="onConfirm">确认</z-btn>
		</div>
	</lp-dialog>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells, { permissionHeadersCheckboxs, netHeadersCheckboxs } from "./cells";

export default {
	name: "PermissionDialog",
	mixins: [TableMixins, DialogMixins],

	data() {
		return {
			formId: "PermissionDialog",
			tableHeaderFormId: "StaffDialogTableHeaders",
			tableDessertsFormId: "StaffDialogTableDesserts",
			dispatchForm: "STAFF_UPDATE_USER_PERMISSION",
			cells,
			permissionHeadersCheckboxs: tools.deepClone(permissionHeadersCheckboxs),
			netHeadersCheckboxs: tools.deepClone(netHeadersCheckboxs)
		};
	},

	methods: {
		// 获取用户项目权限列表
		async getPermissionList() {
			const params = {
				userId: this.detail.id
			};

			const result = await this.$store.dispatch("STAFF_GET_USER_PERMISSION_LIST", params);

			if (result.code === 200) {
				this.desserts = result.data?.list;
			}
		},

		// 当前 header 操作的 checkbox
		changeHeader(value, item) {
			// desserts
			this.desserts.map(record => {
				record[item.value] = value;
			});
		},

		// 当前 cell 操作的 checkbox
		changeCell(value, item) {
			const currentCol = [];

			// desserts
			this.desserts.map((record, index) => {
				// 当前 checkbox 状态
				if (item.rowIndex === index) {
					record[item.value] = value;
				}

				// 当前 checkbox 列
				for (let key in record) {
					if (key === item.value) {
						currentCol.push(record[key]);
					}
				}
			});

			const everyTrue = currentCol.every(value => value);
			const someTrue = currentCol.some(value => value);

			// 当前 cell 操作的 checkbox 对应 headers checkbox 的状态
			switch (item.initFormKey) {
				case "op":
					this.permissionHeadersCheckboxs.map(record => {
						if (record.value === item.value) {
							record.selected = everyTrue;

							record.indeterminate = everyTrue ? false : someTrue;
						}
					});
					break;

				case "net":
					this.netHeadersCheckboxs.map(record => {
						if (record.value === item.value) {
							record.selected = everyTrue;

							record.indeterminate = everyTrue ? false : someTrue;
						}
					});
					break;
			}
		},

		onConfirm() {
			const sysProPermission = [];

			this.desserts.map(item => {
				sysProPermission.push({
					hasOp0: item.hasOp0,
					hasOp1: item.hasOp1,
					hasOp2: item.hasOp2,
					hasOpq: item.hasOpq,
					hasInNet: item.hasInNet,
					hasOutNet: item.hasOutNet,
					proCode: item.proCode,
					proName: item.proName,
					userCode: item.userCode,
					hasPm: item.hasPm,
					objectId: item.objectId,
					userId: item.userId,
					proId: item.proId,
					ID: item.ID
				});
			});

			this.submit(sysProPermission);
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					this.permissionHeadersCheckboxs = tools.deepClone(permissionHeadersCheckboxs);
					this.netHeadersCheckboxs = tools.deepClone(netHeadersCheckboxs);

					this.$nextTick(this.getPermissionList);
				}
			},
			immediate: true
		}
	}
};
</script>
