<template>
	<div class="config-path">
		<v-expansion-panels v-model="panel" multiple>
			<path-table></path-table>

			<program-table></program-table>
		</v-expansion-panels>
	</div>
</template>

<script>
export default {
	name: "ConfigPath",

	data() {
		return {
			panel: [0, 1]
		};
	},

	components: {
		"path-table": () => import("./pathTable"),
		"program-table": () => import("./programTable")
	}
};
</script>
