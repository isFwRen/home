<template>
	<div class="info-out-of-danger">
		<v-expansion-panels v-model="panel" multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>出险信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
export default {
	name: "InfoOutOfDanger",

	data() {
		return {
			panel: [0]
		};
	}
};
</script>
