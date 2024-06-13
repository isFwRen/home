<template>
	<div class="config">
		<div class="mb-8 lp-filters">
			<div class="z-flex pt-6 btns">
				<z-btn
					v-for="(item, index) in cells.options"
					:key="`options_${index}`"
					:class="item.class"
					color="primary"
					:disabled="item.disabled"
					outlined
					@click="onOptions(item)"
				>
					<v-icon class="text-h6">{{ item.icon }}</v-icon>
					{{ item.text }}
				</z-btn>
			</div>
		</div>

		<div class="table config-table">
			<vxe-table
				:data="dessert"
				:border="tableBorder"
				:max-height="tableMaxHeight"
				:size="tableSize"
				:stripe="tableStripe"
			>
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'type'"
						:field="item.value"
						:key="item.value"
						:width="item.width"
					>
						<template #header>
							<z-select
								formId="type"
								formKey="itemName"
								class="mt-n3"
								hideDetails
								label="项目名称"
								:options="projectOptions"
								:defaultValue="proId"
								@change="selectProject"
							></z-select>
						</template>

						<template #default="{ row }">
							<z-btn color="primary" outlined small @click="onEditItem(row)">
								{{ row.name }}
							</z-btn>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'template'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div class="z-flex pt-1">
								<div class="">
									<z-btn
										v-for="record in row.sysProTemplate"
										:key="record.ID"
										color="primary"
										class="mb-2"
										outlined
										small
										block
										@click="onEditTemplate(row, record)"
									>
										{{ record.name }}
									</z-btn>
								</div>

								<z-btn
									class="align-self-center ml-4"
									color="primary"
									depressed
									fab
									small
									@click="onNewTemplate(row)"
								>
									<v-icon>mdi-plus</v-icon>
								</z-btn>
							</div>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'autoReturn'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<span>{{ row.autoReturn ? "自动回传" : "手动回传" }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default>
							<div class="z-flex flex-wrap justify-start options">
								<z-btn
									v-for="(record, index) in cells.configs"
									:key="`configs_${index}`"
									:class="record.class"
									color="primary"
									depressed
									rounded
									smaller
									@click="onConfig(record)"
								>
									{{ record.text }}
								</z-btn>
							</div>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'settings'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div class="z-flex flex-wrap justify-start">
								<z-btn
									v-for="(record, index) in cells.settings"
									:key="`configs_${index}`"
									:class="record.class"
									color="primary"
									depressed
									rounded
									smaller
									@click="onSettings(row, record)"
								>
									{{ record.text }}
								</z-btn>
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
		</div>

		<!-- 新增/编辑项目 BEGIN -->
		<update-item-dialog
			ref="updateItem"
			:rowInfo="detailInfo"
			@submitted="getList"
		></update-item-dialog>
		<!-- 新增/编辑项目 END -->

		<!-- 新增模板 BEGIN -->
		<new-template-dialog
			ref="newTemplate"
			:rowInfo="detailInfo"
			@submitted="getList"
		></new-template-dialog>
		<!-- 新增模板 END -->

		<!-- 配置 BEGIN -->
		<item-configs ref="itemConfigs" @close="getList"></item-configs>
		<!-- 配置 END -->
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "./ConfigMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

export default {
	name: "Config",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "config",
			dispatchList: "GET_CONFIG_ITEM_LIST",
			cells,
			projectOptions: [],
			dessert: [],
			templateDesserts: []
		};
	},

	beforeDestroy() {
		this.$store.commit("SET_PROJECT_INFO");
	},

	methods: {
		// 配置
		onOptions(item) {
			if (item.key === 1) {
				// 新建项目
				this.getDetail({});
				this.$refs.updateItem.onOpen(-1);
			} else {
				this.$modal({
					visible: true,
					title: item.text,
					content: `请确认是否要${item.text}？`,
					confirm: () => {}
				});
			}
		},

		// 设置
		async onSettings(row, item) {
			const body = {
				isIntranet: lpTools.isIntranet(),
				innerIp: row.innerIp,
				outIp: row.outIp,
				inAppPort: row.inAppPort,
				outAppPort: row.outAppPort
			};

			if (/^(entry|all)$/.test(item.key)) {
				const result = await this.$store.dispatch("REFRESH_CONFIG_PRO_CONFIG", body);

				this.toasted.dynamic(result.msg, result.code);
			}
			if (/^(manage|all)$/.test(item.key)) {
				console.log("111111", row);
				const result = await this.$store.dispatch("REFRESH_CONFIG_PRO_MANAGE_CONFIG", {
					proCode: row.code
				});
				this.toasted.dynamic(result.msg, result.code);
			}
		},

		// 选择项目
		selectProject(value) {
			for (let item of this.desserts) {
				if (item.ID === value) {
					this.updateProject(value);

					this.dessert = [item];
					return;
				}
			}
		},

		// 编辑项目
		onEditItem(row) {
			this.getDetail(row);
			this.$refs.updateItem.onOpen(1);
		},

		// 编辑模板
		onEditTemplate(row, record) {
			this.getDetail(row);

			this.rememberIds({ tempId: record.ID });

			// this.$router.push({ path: 'config/template' })
			this.$router.push({ path: "/main/PD/config/template" });

			this.$refs.itemConfigs.onOpen();
		},

		// 添加新模板
		onNewTemplate(row) {
			this.getDetail(row);
			this.$refs.newTemplate.onOpen("添加新模板");
		},

		// 配置
		onConfig(record) {
			const { path } = record;
			this.$router.push(path);
			this.$refs.itemConfigs.onOpen();
		},

		// 选中项目
		updateProject(proId) {
			const item = tools.find(this.projectOptions, proId);

			const project = {
				id: proId,
				code: item?.code,
				selectItem: item.label || ""
			};

			this.$store.commit("SET_PROJECT_INFO", project);
			this.rememberIds({ proId });
		}
	},

	computed: {
		...mapGetters(["project"])
	},

	watch: {
		desserts: {
			handler(desserts) {
				if (tools.isYummy(desserts)) {
					this.projectOptions = [];

					for (let item of this.auth.perm) {
						if (item.hasPm && item.proId) {
							this.projectOptions.push({
								label: item.proName,
								value: item.proId,
								code: item.proCode
							});
						}
					}

					if (this.proId) {
						this.dessert = [tools.find(this.desserts, this.proId)];
					} else {
						this.proId = this.projectOptions[0].value;
						this.dessert = [desserts[0]];
					}

					this.updateProject(this.proId);

					this.$store.commit("UPDATE_CONFIG", { pro: this.dessert[0] });
				}
			},
			immediate: true,
			deep: true
		},

		proId: {
			async handler(proId, oldProId) {
				if (proId) {
					if (proId !== oldProId) {
						this.$store.commit("UPDATE_CONFIG", { proId });
					}
				}
			},
			immediate: true
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	components: {
		"update-item-dialog": () => import("./updateItemDialog"),
		"new-template-dialog": () => import("./newTemplateDialog"),
		"item-configs": () => import("./itemConfigs")
	}
};
</script>

<style scoped lang="scss">
.options {
	min-width: 150px;
	max-width: 225px;
}
</style>
