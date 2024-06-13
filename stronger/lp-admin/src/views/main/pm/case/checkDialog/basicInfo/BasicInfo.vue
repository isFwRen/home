<template>
	<div>
		<v-expansion-panels multiple v-model="flag">
			<v-expansion-panel v-for="(item, i) in title" :key="i">
				<v-expansion-panel-header color="#e9f6ff" class="header">
					{{ item }}
				</v-expansion-panel-header>
				<v-expansion-panel-content v-if="i == 0">
					<vxe-form
						:forms="cells.form1"
						:formRules="cells.formRule"
						:belongModule="'basicInfo'"
						:belongModuleForm="'info1'"
						:cloudData="basicInfo.info1"
					></vxe-form>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else-if="i == 1"
					><vxe-form
						:forms="cells.form2"
						:formRules="cells.formRule"
						:belongModule="'basicInfo'"
						:belongModuleForm="'info2'"
						:cloudData="basicInfo.info2"
					></vxe-form>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else-if="i == 2"
					><vxe-form
						:forms="cells.form3"
						:formRules="cells.formRule"
						:belongModule="'basicInfo'"
						:belongModuleForm="'info3'"
						:cloudData="basicInfo.info3"
					></vxe-form>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else>
					<vxe-form
						:forms="cells.form4"
						:formRules="cells.formRule"
						:belongModule="'basicInfo'"
						:belongModuleForm="'info4'"
						:cloudData="basicInfo.info4"
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
		BasicInfo: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			flag: [0, 1, 2, 3],
			title: ["立案信息", "申请人信息", "投保人信息", "委托人信息"],
			cells,
			basicInfo: {}
		};
	},
	methods: {},

	watch: {
		BasicInfo: {
			handler(newValue) {
				this.basicInfo = JSON.parse(JSON.stringify(newValue));
				// console.log("this.basicInfo", this.basicInfo);
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