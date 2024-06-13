<template>
	<div class="config-prescription">
		<lp-tabs class="mb-6" :options="tabsOptions" @change="onTab"></lp-tabs>

		<settings-prescription v-if="currentTab === 1"></settings-prescription>

		<settings-contrach v-else-if="currentTab === 3"></settings-contrach>

		<settings-holiday v-else></settings-holiday>
	</div>
</template>

<script>
export default {
	name: "ConfigPrescription",

	data() {
		return {
			tabsOptions: [
				{
					value: 1,
					label: "时效设置"
				},

				{
					value: 2,
					label: "节假日设置"
				},

				{
					value: 3,
					label: "合同设置"
				}
			],
			currentTab: 1
		};
	},

	methods: {
		onTab(item) {
			this.currentTab = item.value;
		}
	},

	components: {
		"settings-prescription": () => import("./settingsPrescription"),
		"settings-holiday": () => import("./settingsHoliday"),
		"settings-contrach": () => import("./settingsContract")
	}
};
</script>
