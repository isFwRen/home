<template>
	<div class="info-beneficiary">
		<v-expansion-panels v-model="panel" multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>受益人信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>

			<v-expansion-panel>
				<v-expansion-panel-header>受益人详细信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>

			<v-expansion-panel>
				<v-expansion-panel-header>领款人详细信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
export default {
	name: "InfoBeneficiary",

	data() {
		return {
			panel: [0, 1, 2]
		};
	}
};
</script>
