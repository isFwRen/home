<template>
	<div class="yield">
		<lp-tabs
			class="mb-6"
			:options="cells.tabsOptions"
			:defaultValue="currentTab"
			@change="onTab"
		></lp-tabs>

		<router-view></router-view>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import cells from "./cells";

export default {
	name: "Yield",

	data() {
		return {
			formId: "yield",
			cells,
			currentTab: 0
		};
	},

	created() {
		this.saveProToStore();
	},

	methods: {
		saveProToStore() {
			const proOptions = [];
			const proMap = {};
			for (let index = 0; index < this.auth?.Perm?.length; index++) {
				const element = this.auth?.Perm[index];
				proOptions.push({
					label: element["proCode"],
					value: element["proCode"]
				});
				proMap[element["proCode"]] = element;
			}
			console.log(proOptions);
			console.log(proMap);
			// this.storage.set("pro", { proOptions, proMap });
			this.$store.commit("UPDATE_PRO", { proOptions, proMap });
		},

		onTab({ path }) {
			this.$router.push({ path });
		}
	},

	computed: {
		...mapGetters(["auth"])
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
