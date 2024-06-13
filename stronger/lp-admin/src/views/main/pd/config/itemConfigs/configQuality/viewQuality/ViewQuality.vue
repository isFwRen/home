<template>
	<div class="view-quality">
		<lp-dialog ref="dialog" fullscreen title="案件信息">
			<div class="pt-3" slot="main">
				<v-subheader class="px-0">案件号：1111111111</v-subheader>

				<v-divider></v-divider>

				<lp-tabs class="mb-6" :options="panelInfos" @change="onTab"></lp-tabs>

				<div v-for="(ele, i) in panelInfos" :key="i">
					<div v-if="currentTab === ele.value">
						<div class="info-basic">
							<v-expansion-panels v-model="ele.panel" multiple>
								<v-expansion-panel v-for="(item, k) in ele.items" :key="k">
									<v-expansion-panel-header>{{
										item.label
									}}</v-expansion-panel-header>
									<v-expansion-panel-content>
										<v-row>
											<v-col
												v-for="(cItem, cIndex) in info[item.key]"
												:key="`BasicInfo-Child${cIndex}`"
												:cols="4"
											>
												<div>
													<label class="mr-3"
														>{{ cItem.fieldCode
														}}{{ cItem.fieldName }}</label
													>

													<z-text-field
														v-if="cItem.inputType === 1"
														:formId="cItem.ID"
														:formKey="cItem.belongType + ''"
														:readonly="true"
														width="250"
														:defaultValue="cItem.xmlNodeName"
													>
													</z-text-field>

													<z-select
														v-if="cItem.inputType === 2"
														:formId="cItem.ID"
														:formKey="cItem.belongType + ''"
														width="250"
														:options="[
															{
																label: cItem.xmlNodeName,
																value: cItem.xmlNodeName
															}
														]"
														:defaultValue="cItem.xmlNodeName"
													></z-select>

													<div
														v-if="cItem.inputType === 3"
														:class="['z-flex', item.class]"
													>
														<z-select
															:formId="cItem.ID"
															:formKey="cItem.belongType + ''"
															width="100"
															:options="[
																{
																	label: cItem.xmlNodeName,
																	value: cItem.xmlNodeName
																}
															]"
															:defaultValue="cItem.xmlNodeName"
														></z-select>

														<z-text-field
															:formId="cItem.ID"
															:formKey="cItem.belongType + ''"
															:readonly="true"
															width="150"
															:defaultValue="cItem.xmlNodeName"
														>
														</z-text-field>
													</div>
												</div>
											</v-col>
										</v-row>
									</v-expansion-panel-content>
								</v-expansion-panel>
							</v-expansion-panels>
						</div>
					</div>
				</div>
				<!-- <info-basic v-if="currentTab === 1"></info-basic>

        <info-beneficiary
          v-else-if="currentTab === 2"
        ></info-beneficiary>

        <info-bill v-else-if="currentTab === 3"></info-bill>

        <info-out-of-danger v-else></info-out-of-danger> -->
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
import { mapGetters } from "vuex";
export default {
	name: "ViewQuality",
	mixins: [DialogMixins],

	data() {
		return {
			panelInfos: [
				{
					items: [
						{
							key: "1",
							label: "申请人信息"
						},
						{
							key: "2",
							label: "被保人信息"
						},
						{
							key: "3",
							label: "受托人信息"
						},
						{
							key: "4",
							label: "其他信息"
						}
					],
					panel: [0, 1, 2, 3],
					value: 1,
					label: "基础信息"
				},
				{
					items: [
						{
							key: "5",
							label: "受益人信息"
						},
						{
							key: "6",
							label: "领款人详细信息"
						}
					],
					panel: [0, 1],
					value: 2,
					label: "受益人信息"
				},
				{
					items: [
						{
							key: "7",
							label: "账单信息"
						}
					],
					panel: [0],
					value: 3,
					label: "账单信息"
				},
				{
					items: [
						{
							key: "8",
							label: "出险信息"
						}
					],
					panel: [0],
					value: 4,
					label: "出险信息"
				}
			],
			currentTab: 1,
			info: {}
		};
	},
	methods: {
		// 获取数据
		async getConfigQualityFormat() {
			const result = await this.$store.dispatch("GET_CONFIG_QUALITY_FORMAT", {
				proId: this.storage.get("config").proId
			});
			console.log(result);
			this.toasted.dynamic(result.msg, result.code);
			if (result.code === 200) {
				this.info = result.data;
			}
		},
		onTab(item) {
			this.currentTab = item.value;
		}
	},
	computed: {
		...mapGetters(["config"])
	},
	components: {
		"info-basic": () => import("./infoBasic"),
		"info-beneficiary": () => import("./infoBeneficiary"),
		"info-bill": () => import("./infoBill"),
		"info-out-of-danger": () => import("./infoOutOfDanger")
	}
};
</script>
