<template>
	<div class="info-bill">
		<v-expansion-panels v-model="panel" multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>账单信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
export default {
	name: "InfoBill",

	data() {
		return {
			panel: [0]
		};
	}
};
</script>
