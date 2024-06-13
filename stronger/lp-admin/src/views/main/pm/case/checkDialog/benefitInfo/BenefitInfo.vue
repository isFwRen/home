<template>
	<div>
		<v-expansion-panels multiple v-model="flag">
			<v-expansion-panel v-for="(item, i) in title" :key="i">
				<v-expansion-panel-header color="#e9f6ff" class="header">
					{{ item }}
				</v-expansion-panel-header>
				<v-expansion-panel-content v-if="i == 0">
					<vxe-table
						ref="xTable"
						:data="desserts"
						border
						stripe
						align="center"
						style="margin-top: 15px"
						min-height="90"
						class="mytable-scrollbar"
						:edit-config="{ trigger: 'dblclick', mode: 'cell' }"
						@edit-closed="editClosedEvent"
					>
						<vxe-column type="seq" width="60" title="序号"></vxe-column>
						<template v-for="item in cells.headers">
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
							<vxe-column
								v-else
								:field="item.value"
								:fixed="item.fixed"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:sortable="item.sortable"
								:edit-render="{ autofocus: '.vxe-input--inner' }"
							>
								<template #edit="{ row }">
									<vxe-input v-model="row[item.value]" type="text"></vxe-input>
								</template>
							</vxe-column>
						</template>
					</vxe-table>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else-if="i == 1"
					><vxe-form
						:forms="cells.form2"
						:formRules="cells.formRule"
						:belongModule="'benefitInfo'"
						:belongModuleForm="'info2'"
						:cloudData="benefitInfo.info2"
					></vxe-form>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else-if="i == 2"
					><vxe-form
						:forms="cells.form3"
						:formRules="cells.formRule"
						:belongModule="'benefitInfo'"
						:belongModuleForm="'info3'"
						:cloudData="benefitInfo.info3"
					></vxe-form>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else>
					<vxe-form
						:forms="cells.form4"
						:formRules="cells.formRule"
						:belongModule="'benefitInfo'"
						:belongModuleForm="'info4'"
						:cloudData="benefitInfo.info4"
					></vxe-form
				></v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
import cells from "./cells";
export default {
	props: {
		BenefitInfo: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			flag: [0, 1, 2, 3],
			title: [
				"理赔受益人信息",
				"理赔受益人详细信息",
				"理赔受益人详细信息(自然人)",
				"理赔领款人详细信息(自然人)"
			],
			cells,
			desserts: [
				{
					beneficiaryName: "",
					beneficiaryType: "",
					getMoneyerName: "",
					beneficiaryMoney: "",
					incomeRation: "",
					payment: ""
				}
			],
			benefitInfo: {}
		};
	},
	methods: {
		editClosedEvent() {
			const $table = this.$refs.xTable;
			let content = sessionStorage.get("checkForm");
			content.benefitInfo.desserts = $table.getData();
			sessionStorage.set("checkForm", content);
		}
	},

	watch: {
		BenefitInfo: {
			handler(newValue) {
				this.benefitInfo = JSON.parse(JSON.stringify(newValue));
				if (this.benefitInfo.desserts && this.benefitInfo.desserts.length != 0) {
					this.desserts = this.benefitInfo.desserts;
				}
				// console.log("benefitInfo", newValue);
			},
			immediate: true,
			deep: true
		}
	},

	components: {
		"vxe-form": () => import("../vxeForm")
	}
};
</script>

<style lang="scss">
.header {
	color: #007aff;
	font-weight: bolder;
	font-size: 15px;
}
.v-expansion-panel--active > .v-expansion-panel-header {
	min-height: 0px;
}
</style>