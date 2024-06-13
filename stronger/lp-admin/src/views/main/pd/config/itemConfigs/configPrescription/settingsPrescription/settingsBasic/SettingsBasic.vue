<template>
	<v-expansion-panel>
		<v-expansion-panel-header disable-icon-rotate>
			基础设置
			<template v-slot:actions>
				<z-btn color="primary" depressed outlined small @click.stop="onNew">
					<v-icon class="text-h6" color="primary">mdi-plus</v-icon>
					新增
				</z-btn>
			</template>
		</v-expansion-panel-header>

		<v-expansion-panel-content>
			<div class="settings-basic-table">
				<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
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
							<template #default="{ row }">
								{{ row[item.value] | dateFormat(item.format) }}
							</template>
						</vxe-column>
					</template>
				</vxe-table>
			</div>

			<!-- 新增/编辑 BEGIN -->
			<z-dynamic-form
				ref="dynamic"
				:formId="formId"
				:detail="detailInfo"
				:fieldList="cells.fields"
				@confirm="handleConfirm"
			></z-dynamic-form>
			<!-- 新增/编辑 END -->
		</v-expansion-panel-content>
	</v-expansion-panel>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import { R } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";

export default {
	name: "SettingsBasic",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "SettingsBasic",
			dispatchList: "GET_CONFIG_PRESCRIPTION_LIST",
			dispatchDelete: "DELETE_CONFIG_PRESCRIPTION_ITEM",
			cells,
			effectParams: {
				configType: "base"
			}
		};
	},

	methods: {
		onNew() {
			this.getDetail({});
			this.$refs.dynamic.open({ title: "新增", status: -1 });
		},

		onEdit(row) {
			const detail = { ...row };

			for (let item of this.cells.fields) {
				if (item.dateFormat && R.isYummy(row[item.formKey])) {
					detail[item.formKey] = moment(
						`${moment().format("YYYY-MM-DD")} ${row[item.formKey]}`
					).format(item.dateFormat);
				}
			}

			this.getDetail(detail);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		handleConfirm(effect, form) {
			form = {
				proId: this.config.proId,
				id: this.detailInfo.ID,
				configType: "base",
				status: effect.status,
				...form
			};
			this.updateListItem(form, "UPDATE_CONFIG_PRESCRIPTION");
		},

		onMore({ customValue }, row) {
			this.getDetail(row);
			customValue === "delete" && this.deleteItem();
		}
	},

	computed: {
		...mapGetters(["config"])
	}
};
</script>
