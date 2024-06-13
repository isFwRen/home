<template>
	<div class="quality">
		<lp-tabs
			class="mb-4"
			:options="cells.tabsOptions"
			:defaultValue="currentTab"
			@change="onTab"
		></lp-tabs>

		<router-view></router-view>
	</div>
</template>

<script>
import cells from "./cells";

export default {
	name: "Quality",

	data() {
		return {
			cells,
			currentTab: 0
		};
	},

	methods: {
		onTab({ path }) {
			this.$router.push({ path });
		}
	},

	watch: {
		$route: {
			handler({ meta }) {
				this.cells.tabsOptions.map((tab, tabIndex) => {
					if (tab.key === meta.key) {
						this.currentTab = tabIndex;
					}
				});
			},
			immediate: true
		}
	}
};
</script>
